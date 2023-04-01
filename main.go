package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	twiliorest "github.com/twilio/twilio-go/rest/api/v2010"

	"dstwilio/twilioclient"
)

func main() {
	err := godotenv.Load(".envrc")
	if err != nil {
		fmt.Println("Error loading .envrc file:", err)
		return
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
		makePhoneCall()
		sendTextMessage(fmt.Sprintf("📢 [ RECEIVED MESSAGE FROM: %s IN %s CHANNEL]", m.Author.Username, targetChannel))
	}
}

func makePhoneCall() {
	client := twilioclient.NewTwilioClientWithoutKeepAlives()

	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER")
	toPhone := os.Getenv("YOUR_PHONE_NUMBER")
	twimlURL := "http://demo.twilio.com/docs/voice.xml"

	params := &twiliorest.CreateCallParams{
		From: &fromPhone,
		To:   &toPhone,
		Url:  &twimlURL,
	}

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
