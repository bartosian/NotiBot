package twilioclient

import (
	"net/http"
	"os"

	"github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
)

type NoConnReuseClient struct {
	twilioClient.Client
}

func NewTwilioClientWithoutKeepAlives() *twilio.RestClient {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	transport := &http.Transport{
		DisableKeepAlives: true,
	}

	customClient := &NoConnReuseClient{
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

	customClient.SetAccountSid(accountSID)

	return twilio.NewRestClientWithParams(twilio.ClientParams{
		Client: customClient,
	})
}
