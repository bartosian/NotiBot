package ports

import (
	"github.com/bwmarrin/discordgo"
)

type NotifierGateway interface {
	CreateCall(message string) error
	SendMessage(message string) error
}

type (
	MessageHandler = func(s *discordgo.Session, m *discordgo.MessageCreate)

	DiscordGateway interface {
		AddHandler(handler MessageHandler) error
	}
)
