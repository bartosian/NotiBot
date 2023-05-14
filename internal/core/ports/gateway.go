package ports

import "dstwilio/internal/core/domain"

type NotifierGateway interface {
	CreateCall(call domain.Call) error
	SendMessage(message domain.Message) error
}
