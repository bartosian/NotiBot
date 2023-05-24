package twiliogw

import (
	"os"

	"github.com/twilio/twilio-go"

	"github.com/bartosian/notibot/internal/core/ports"
	"github.com/bartosian/notibot/pkg/l0g"
	"github.com/bartosian/notibot/pkg/twilioclient"
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
		client:          twilioclient.NewTwilioClient(),
		fromPhoneNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
		toPhoneNumber:   os.Getenv("CLIENT_PHONE_NUMBER"),
		logger:          logger,
	}
}
