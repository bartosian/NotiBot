package domain

type Message struct {
	From string // The phone number or identifier from which the message is being sent
	To   string // The phone number or identifier to which the message is being sent
	Body string // The content of the message
}
