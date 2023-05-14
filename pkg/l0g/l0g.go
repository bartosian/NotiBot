package l0g

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type Logger interface {
	Info(msg string, data interface{})
	Warn(msg string, data interface{})
	Error(msg string, err error, data interface{})
	Sugar() *zap.SugaredLogger
	Sync() error
}

type Log struct {
	internal *zap.Logger
}

type Field = zap.Field

var (
	Any           = zap.Any
	String        = zap.String
	Bool          = zap.Bool
	Object        = zap.Object
	Namespace     = zap.Namespace
	AddCallerSkip = zap.AddCallerSkip
)

func NewLogger(options ...zap.Option) Logger {
	// Automatically add CallerSkip for service loggers
	options = append([]zap.Option{AddCallerSkip(1)}, options...)
	logger, _ := zap.NewProduction(options...)

	return &Log{
		internal: logger,
	}
}

func NewUnaryLogger(options ...zap.Option) *zap.Logger {
	logger, _ := zap.NewProduction(options...)

	return logger
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

func (logger *Log) Sugar() *zap.SugaredLogger {
	return logger.internal.Sugar()
}

func (logger *Log) Sync() error {
	return logger.internal.Sync()
}

// PrintDelimiter prints a delimiter line to the console. This can be used to visually separate sections of console output.
func PrintDelimiter() {
	delimiter := strings.Repeat("-", 60)

	fmt.Println(delimiter)
}
