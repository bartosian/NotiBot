package domain

type Call struct {
	From         string // The phone number or client identifier initiating the call
	To           string // The phone number or client identifier receiving the call
	TwimlMessage string // The TwiML message to be played during the call
}
