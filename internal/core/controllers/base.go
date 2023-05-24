package controllers

import (
	"github.com/bartosian/notibot/internal/core/ports"
	"github.com/bartosian/notibot/pkg/l0g"
)

// NotifierController handles the logic for sending notifications.
type NotifierController struct {
	notifierGateway ports.NotifierGateway // Gateway for sending notifications
	discordGateway  ports.DiscordGateway  // Gateway for communicating with Discord API
	logger          l0g.Logger            // Logger for logging purposes
}

// NewNotifierController creates a new instance of the NotifierController.
func NewNotifierController(
	notifierGateway ports.NotifierGateway,
	discordGateway ports.DiscordGateway,
	logger l0g.Logger,
) ports.NotifierController {
	return &NotifierController{
		notifierGateway: notifierGateway, // Assign the provided notifier gateway
		discordGateway:  discordGateway,  // Assign the provided discord gateway
		logger:          logger,          // Assign the provided logger
	}
}
