package controllers

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/bartosian/notibot/internal/core/ports"
)

const (
	callVoiceTemplate   = "RECEIVED MESSAGE FROM %s IN %s CHANNEL"
	messageTextTemplate = "📢 [ RECEIVED MESSAGE FROM %s IN %s CHANNEL ]"
)

var targetChannel = os.Getenv("DISCORD_CHANNEL")

// MonitorDiscord starts monitoring the Discord channel for new messages by adding the message handler.
// It returns an error if there is an issue adding the message handler.
func (c *NotifierController) MonitorDiscord() error {
	messageHandler := c.newMessageHandler()

	if err := c.discordGateway.AddHandler(messageHandler); err != nil {
		c.logger.Error("error adding message handler", err, nil)

		return err
	}

	return nil
}

func (c *NotifierController) newMessageHandler() ports.MessageHandler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		targetChannel = os.Getenv("DISCORD_CHANNEL")

		messageChannel, err := s.Channel(m.ChannelID)
		if err != nil {
			c.logger.Error("error getting channel by id", err, nil)

			return
		}

		if messageChannel.Name == targetChannel {
			err := c.notifierGateway.CreateCall(buildCallVoice(m.Author.Username, targetChannel))
			if err != nil {
				c.logger.Error("error creating call", err, nil)

				os.Exit(1)
			}

			err = c.notifierGateway.SendMessage(buildMessageText(m.Author.Username, targetChannel))
			if err != nil {
				c.logger.Error("error sending message", err, nil)

				os.Exit(1)
			}
		}
	}
}

// buildDiscordCallVoice creates the call voice content using the provided channel names.
func buildCallVoice(fromChannel, toChannel string) string {
	return fmt.Sprintf(callVoiceTemplate, fromChannel, toChannel)
}

// buildDiscordMessageText creates the message text using the provided channel names.
func buildMessageText(fromChannel, toChannel string) string {
	return fmt.Sprintf(messageTextTemplate, fromChannel, toChannel)
}
