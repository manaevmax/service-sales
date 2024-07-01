package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(level string) *Logger {
	c := zap.NewProductionConfig()

	zapLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}

	c.DisableCaller = true
	c.DisableStacktrace = true
	c.Encoding = "json"
	c.Level = zapLevel
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	l, err := c.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{l.Sugar()}
}

func NoOpLogger() *Logger {
	return &Logger{zap.NewNop().Sugar()}
}
