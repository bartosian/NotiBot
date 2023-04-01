package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
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
	callVoiceTemplate = "RECEIVED MESSAGE FROM %s IN %s CHANNEL"
	messageTemplate   = "📢 [ RECEIVED MESSAGE FROM %s IN %s CHANNEL ]"

	alertVoiceTemplate   = "NEW ALERT TRIGGERED!\n\n*ALERT NAME:* %s\n*SUMMARY:* %s\n*START TIME:* %s\n*END TIME:* %s\n*STATUS:* %s\n\n"
	alertMessageTemplate = "🚨 NEW ALERT TRIGGERED!\n\n🔴 ALERT NAME: %s\n\n🔴 SUMMARY: %s\n\n🔴 START TIME: %s\n\n🔴 END TIME: %s\n\n🔴 STATUS: %s\n\n"
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
		if os.Getenv("ALERT_MANAGER_URL") == "" {
			fmt.Printf("Error: Environment variable %s is not set\n", "ALERT_MANAGER_URL")

			return
		}

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

	fmt.Println(messageContent)
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

type AlertReceiver struct {
	Name string `json:"name"`
}

type AlertStatus struct {
	InhibitedBy []string `json:"inhibitedBy"`
	SilencedBy  []string `json:"silencedBy"`
	State       string   `json:"state"`
}

type AlertLabels struct {
	Alertname string `json:"alertname"`
}

type AlertAnnotations struct {
	Summary string `json:"summary"`
}

type AlertData struct {
	Annotations  AlertAnnotations `json:"annotations"`
	EndsAt       string           `json:"endsAt"`
	Fingerprint  string           `json:"fingerprint"`
	Receivers    []AlertReceiver  `json:"receivers"`
	StartsAt     string           `json:"startsAt"`
	Status       AlertStatus      `json:"status"`
	UpdatedAt    string           `json:"updatedAt"`
	GeneratorURL string           `json:"generatorURL"`
	Labels       AlertLabels      `json:"labels"`
}

func checkAlerts() {
	resp, err := http.Get(os.Getenv("ALERT_MANAGER_URL"))
	if err != nil {
		log.Fatal("Error fetching alerts: ", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	var alertList []AlertData

	err = json.Unmarshal(data, &alertList)
	if err != nil {
		log.Fatal("error parsing alert list", err)
	}

	if len(alertList) == 0 {
		log.Println("no active alerts found")

		return
	}

	var (
		voiceBody   string
		messageBody string
	)

	for _, alert := range alertList {
		voiceBody += fmt.Sprintf(alertVoiceTemplate,
			alert.Labels.Alertname,
			alert.Annotations.Summary,
			alert.StartsAt,
			alert.EndsAt,
			alert.Status.State,
		)

		messageBody += fmt.Sprintf(alertMessageTemplate,
			alert.Labels.Alertname,
			alert.Annotations.Summary,
			alert.StartsAt,
			alert.EndsAt,
			alert.Status.State,
		)
	}

	getDelimiter()
	getCurrentTime()
	makePhoneCall(voiceBody)
	sendTextMessage(messageBody)
	getDelimiter()
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
