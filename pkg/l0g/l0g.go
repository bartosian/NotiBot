package l0g

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, data interface{})
	Warn(msg string, data interface{})
	Error(msg string, err error, data interface{})
}

type Log struct {
	internal *zap.Logger
}

type Field = zap.Field

var AddCallerSkip = zap.AddCallerSkip

func NewLogger(options ...zap.Option) Logger {
	options = append([]zap.Option{AddCallerSkip(1)}, options...)
	logger, _ := zap.NewProduction(options...)

	return &Log{
		internal: logger,
	}
}

func (logger *Log) Info(msg string, data interface{}) {
	logger.internal.Info(msg, zap.Any("data", data))
}

func (logger *Log) Warn(msg string, data interface{}) {
	logger.internal.Warn(msg, zap.Any("data", data))
}

func (logger *Log) Error(msg string, err error, data interface{}) {
	logger.internal.Error(msg, zap.Error(err), zap.Any("data", data))
}
