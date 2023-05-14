package twiliogw

import (
	"github.com/twilio/twilio-go"

	"dstwilio/internal/core/ports"
	"dstwilio/pkg/l0g"
	"dstwilio/pkg/twilioclient"
)

type Gateway struct {
	client *twilio.RestClient // Twilio REST client used for communication
	logger l0g.Logger         // Logger for logging purposes
}

// NewTwilioGateway creates a new instance of the Twilio gateway.
func NewTwilioGateway(logger l0g.Logger) ports.NotifierGateway {
	return &Gateway{
		client: twilioclient.NewTwilioClientWithoutKeepAlives(), // Create a Twilio REST client without keep-alives
		logger: logger,                                          // Assign the provided logger to the gateway
	}
}
