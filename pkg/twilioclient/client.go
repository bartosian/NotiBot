package twilioclient

import (
	"net/http"
	"os"

	"github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
)

type Client struct {
	twilioClient.Client
}

func NewTwilioClient() *twilio.RestClient {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	transport := &http.Transport{
		DisableKeepAlives: true,
	}

	client := &Client{
		Client: twilioClient.Client{
			Credentials: &twilioClient.Credentials{
				Username: accountSID,
				Password: authToken,
			},
			HTTPClient: &http.Client{
				Transport: transport,
			},
		},
	}

	client.SetAccountSid(accountSID)

	return twilio.NewRestClientWithParams(twilio.ClientParams{
		Client: client,
	})
}
