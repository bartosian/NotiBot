package twiliogw

import (
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendMessage sends a message using the Twilio gateway.
func (t *Gateway) SendMessage(message string) error {
	messageParams := t.buildMessageParams(message)

	result, err := t.client.Api.CreateMessage(messageParams)
	if err != nil {
		t.logger.Error("failed to send Message", err, messageParams)
	}

	t.logger.Info("message successfully sent", result)

	return err
}

// buildMessageParams builds the Twilio message parameters based on the domain.Message object.
func (t *Gateway) buildMessageParams(message string) *openapi.CreateMessageParams {
	return &openapi.CreateMessageParams{
		From: &t.fromPhoneNumber,
		To:   &t.toPhoneNumber,
		Body: &message,
	}
}
