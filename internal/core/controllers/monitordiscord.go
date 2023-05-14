package controllers

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

const (
	callVoiceTemplate   = "RECEIVED MESSAGE FROM %s IN %s CHANNEL"
	messageTextTemplate = "📢 [ RECEIVED MESSAGE FROM %s IN %s CHANNEL ]"
)

var targetChannel = os.Getenv("DISCORD_CHANNEL")

// MonitorDiscord starts monitoring the Discord channel for new messages by adding the message handler.
// It returns an error if there is an issue adding the message handler.
func (c *NotifierController) MonitorDiscord() error {
	if err := c.discordGateway.AddHandler(c.newMessageHandler); err != nil {
		c.logger.Error("error adding message handler", err, nil)

		return err
	}

	return nil
}

// newMessageHandler is the handler function for new messages in the Discord channel.
// It checks if the message author is the channel owner and skips further processing if true.
// Otherwise, it retrieves the Discord channel by ID and performs actions based on the target channel.
func (c *NotifierController) newMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// if message originator is the channel owner - dismissing
	if message.Author.ID == session.State.User.ID {
		return
	}

	discordChannel, err := session.Channel(message.ChannelID)
	if err != nil {
		c.logger.Error("error getting channel by id", err, nil)

		return
	}

	if discordChannel.Name == targetChannel {
		err := c.notifierGateway.CreateCall(buildCallVoice(message.Author.Username, targetChannel))
		if err != nil {
			c.logger.Error("error creating call", err, nil)

			os.Exit(1)
		}

		err = c.notifierGateway.SendMessage(buildMessageText(message.Author.Username, targetChannel))
		if err != nil {
			c.logger.Error("error sending message", err, nil)

			os.Exit(1)
		}
	}
}

// buildCallVoice creates the call voice content using the provided channel names.
func buildCallVoice(fromChannel, toChannel string) string {
	return fmt.Sprintf(callVoiceTemplate, fromChannel, toChannel)
}

// buildMessageText creates the message text using the provided channel names.
func buildMessageText(fromChannel, toChannel string) string {
	return fmt.Sprintf(messageTextTemplate, fromChannel, toChannel)
}
