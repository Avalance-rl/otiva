package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func New() *Logger {
	return &Logger{setupLogger()}
}

func setupLogger() *zap.Logger {
	encoderConfig := ecszap.EncoderConfig{
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   ecszap.FullCallerEncoder,
	}

	stdoutCore := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	stderrCore := ecszap.NewCore(encoderConfig, os.Stderr, zap.ErrorLevel)

	core := zapcore.NewTee(stdoutCore, stderrCore)

	logger := zap.New(core, zap.AddCaller())

	return logger
}
