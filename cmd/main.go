package main

import (
	"errors"
	"flag"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"dstwilio/internal/core/controllers"
	"dstwilio/internal/core/ports"
	"dstwilio/internal/gateways/discordgw"
	"dstwilio/internal/gateways/twiliogw"
	"dstwilio/pkg/l0g"
)

const (
	envFilePath = "../.envrc"
)

var (
	requiredTwilioVars  = []string{"TWILIO_ACCOUNT_SID", "TWILIO_AUTH_TOKEN", "TWILIO_PHONE_NUMBER", "CLIENT_PHONE_NUMBER"}
	requiredDiscordVars = []string{"DISCORD_BOT_TOKEN", "DISCORD_CHANNEL"}
	requiredAlertVars   = []string{"ALERT_MANAGER_URL"}
	errEnvVarNotFound   = errors.New("environment variable not found")
)

func main() {
	logger := l0g.NewLogger()
	flagSet := parseFlags()

	checkMonitors(flagSet, logger)

	if flagSet.readEnv {
		loadEnvFile(logger)
	}

	checkEnvVars(requiredTwilioVars, logger)

	// Instantiate gateways
	notifierGateway := twiliogw.NewTwilioGateway(logger)
	discordGateway, err := discordgw.NewDiscordGateway(logger)
	if err != nil {
		logger.Error("error creating discord gateway", err, nil)

		os.Exit(1)
	}

	// Instantiate controllers
	notifierController := controllers.NewNotifierController(notifierGateway, discordGateway, logger)

	monitor(flagSet, notifierController, logger)
}

type flags struct {
	readEnv        bool
	monitorAlerts  bool
	monitorDiscord bool
}

// parseFlags parses and returns command line flags.
func parseFlags() *flags {
	readEnv := flag.Bool("envrc", false, "read environment variables from .envrc file")
	monitorAlerts := flag.Bool("alerts", false, "receive alerts from grafana alert manager")
	monitorDiscord := flag.Bool("discord", false, "receive notifications from discord channels")

	flag.Parse()

	return &flags{
		readEnv:        *readEnv,
		monitorAlerts:  *monitorAlerts,
		monitorDiscord: *monitorDiscord,
	}
}

// checkMonitors checks if both alert and discord monitors are disabled. If they are, it logs an info message and exits the program.
func checkMonitors(flagSet *flags, logger l0g.Logger) {
	if !flagSet.monitorAlerts && !flagSet.monitorDiscord {
		logger.Info("all monitors disabled - dismissing", nil)

		os.Exit(0)
	}
}

// loadEnvFile loads environment variables from a .envrc file. If an error occurs, it logs the error and exits the program.
func loadEnvFile(logger l0g.Logger) {
	if err := godotenv.Load(envFilePath); err != nil {
		logger.Error("error loading .envrc file:", err, nil)

		os.Exit(1)
	}
}

// checkEnvVars checks for the presence of each environment variable in the provided slice. If any are missing, it logs an error and exits the program.
func checkEnvVars(envVars []string, logger l0g.Logger) {
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			logger.Error("error lookup for environment variable", errEnvVarNotFound, envVar)

			os.Exit(1)
		}
	}
}

// monitor checks which monitors are enabled via the provided flags, and starts those monitors. If an error occurs during monitoring, it logs the error and exits the program.
func monitor(flagSet *flags, notifierController ports.NotifierController, logger l0g.Logger) {
	var errGroup errgroup.Group

	if flagSet.monitorDiscord {
		checkEnvVars(requiredDiscordVars, logger)

		errGroup.Go(func() error {
			return notifierController.MonitorDiscord()
		})
	}

	if flagSet.monitorAlerts {
		checkEnvVars(requiredAlertVars, logger)

		errGroup.Go(func() error {
			return notifierController.MonitorAlerts()
		})
	}

	if err := errGroup.Wait(); err != nil {
		logger.Error("unexpected error occurred:", err, nil)

		os.Exit(1)
	}
}
