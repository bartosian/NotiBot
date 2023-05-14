package twiliogw

import (
	"fmt"

	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"dstwilio/internal/core/domain"
)

// CreateCall creates a new call using the Twilio gateway.
func (t *Gateway) CreateCall(call domain.Call) error {
	callParams := t.buildCallParams(call)

	result, err := t.client.Api.CreateCall(callParams)
	if err != nil {
		t.logger.Error("failed to execute Call", err, callParams)
	}

	t.logger.Info("call successfully executed", result)

	return err
}

// buildCallParams builds the Twilio call parameters based on the domain.Call object.
func (t *Gateway) buildCallParams(call domain.Call) *openapi.CreateCallParams {
	return &openapi.CreateCallParams{
		From:  &call.From,
		To:    &call.To,
		Twiml: buildTwiml(call.TwimlMessage),
	}
}

// buildTwiml builds the TwiML (Twilio Markup Language) for the given message.
func buildTwiml(message string) *string {
	twiml := fmt.Sprintf("<Response><Say>%s</Say></Response>", message)

	return &twiml
}
