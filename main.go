package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"dstwilio/twilioclient"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	twiliorest "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	alertManagerAPI   = "http://localhost:9093/api/v2/alerts"
	callVoiceTemplate = "RECEIVED MESSAGE FROM %s IN %s CHANNEL"
	messageTemplate   = "📢 [ RECEIVED MESSAGE FROM %s IN %s CHANNEL ]"

	alertVoiceTemplate   = "ALERT TRIGGERED: %s (SEVERITY: %s, LABELS: %v)"
	alertMessageTemplate = "🔥 [ ALERT TRIGGERED: %s (SEVERITY: %s, LABELS: %v) ]"
)

func main() {
	readEnvrc := flag.Bool("envrc", false, "read environment variables from .envrc file")
	enableAlerts := flag.Bool("alerts", false, "receive alerts from grafana alert manager")
	flag.Parse()

	if *readEnvrc {
		err := godotenv.Load(".envrc")
		if err != nil {
			fmt.Println("Error loading .envrc file:", err)
			return
		}
	}

	requiredVars := []string{"TWILIO_ACCOUNT_SID", "TWILIO_AUTH_TOKEN", "TWILIO_PHONE_NUMBER", "YOUR_PHONE_NUMBER", "DISCORD_BOT_TOKEN", "DISCORD_CHANNEL"}
	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			fmt.Printf("Error: Environment variable %s is not set\n", envVar)
			return
		}
	}

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	if *enableAlerts {
		go func() {
			for {
				checkAlerts()
				time.Sleep(time.Minute)
			}
		}()
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var targetChannel = os.Getenv("DISCORD_CHANNEL")

	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel:", err)
		return
	}

	if channel.Name == targetChannel {
		getDelimiter()
		getCurrentTime()
		makePhoneCall(fmt.Sprintf(callVoiceTemplate, m.Author.Username, targetChannel))
		sendTextMessage(fmt.Sprintf(messageTemplate, m.Author.Username, targetChannel))
		getDelimiter()
	}
}

func makePhoneCall(messageContent string) {
	client := twilioclient.NewTwilioClientWithoutKeepAlives()

	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER")
	toPhone := os.Getenv("YOUR_PHONE_NUMBER")

	params := &twiliorest.CreateCallParams{
		From: &fromPhone,
		To:   &toPhone,
	}

	params.SetTwiml(fmt.Sprintf("<Response><Say>%s</Say></Response>", messageContent))

	_, err := client.Api.CreateCall(params)
	if err != nil {
		fmt.Println("Error making phone call:", err)
		return
	}

	fmt.Println("Phone call initiated.")
}

func sendTextMessage(messageContent string) {
	client := twilioclient.NewTwilioClientWithoutKeepAlives()

	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER")
	toPhone := os.Getenv("YOUR_PHONE_NUMBER")

	params := &twiliorest.CreateMessageParams{
		From: &fromPhone,
		To:   &toPhone,
		Body: &messageContent,
	}

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending text message:", err)
		return
	}

	fmt.Println("Text message sent.")
}

type AlertList struct {
	Data []struct {
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
	} `json:"data"`
}

func checkAlerts() {
	resp, err := http.Get(alertManagerAPI)
	if err != nil {
		fmt.Println("Error retrieving alerts:", err)
		return
	}

	var alerts AlertList
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	if err != nil {
		fmt.Println("Error parsing alert list:", err)
		return
	}

	if len(alerts.Data) > 0 {
		for _, alert := range alerts.Data {
			alertName := alert.Labels["alertname"]
			alertSeverity := alert.Labels["severity"]
			alertLabels := alert.Labels

			getDelimiter()
			getCurrentTime()
			makePhoneCall(fmt.Sprintf(alertVoiceTemplate, alertName, alertSeverity, alertLabels))
			sendTextMessage(fmt.Sprintf(alertMessageTemplate, alertName, alertSeverity, alertLabels))
			getDelimiter()
		}
	}
}

func getCurrentTime() {
	currentTime := time.Now()
	formattedTime := currentTime.Format("Monday, January 2, 2006 at 3:04pm")

	fmt.Println("The current time is:", formattedTime)
}

func getDelimiter() {
	delimiter := strings.Repeat("-", 60)

	fmt.Println(delimiter)
}
