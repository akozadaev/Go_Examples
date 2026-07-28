package main

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (Logger, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: zapLogger}, nil
}

func (zl *ZapLogger) Info(msg string) {
	zl.logger.Info(msg)
}

func (zl *ZapLogger) Error(msg string) {
	zl.logger.Error(msg)
}

func main() {
	logger, err := NewZapLogger()
	if err != nil {
		panic(err)
	}
	defer logger.(*ZapLogger).logger.Sync() // flush buffers

	logger.Info("This is an informational message.")
	logger.Error("This is an error message.")
}
