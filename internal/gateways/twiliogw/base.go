package twiliogw

import (
	"github.com/twilio/twilio-go"
	"os"

	"dstwilio/internal/core/ports"
	"dstwilio/pkg/l0g"
	"dstwilio/pkg/twilioclient"
)

type Gateway struct {
	client          *twilio.RestClient // Twilio REST client used for communication
	fromPhoneNumber string             // Twilio phone number
	toPhoneNumber   string             // Client phone number to receive alerts
	logger          l0g.Logger         // Logger for logging purposes
}

// NewTwilioGateway creates a new instance of the Twilio gateway.
func NewTwilioGateway(logger l0g.Logger) ports.NotifierGateway {
	return &Gateway{
		client:          twilioclient.NewTwilioClientWithoutKeepAlives(),
		fromPhoneNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
		toPhoneNumber:   os.Getenv("CLIENT_PHONE_NUMBER"),
		logger:          logger,
	}
}
