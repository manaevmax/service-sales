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

func (l *Logger) Printf(format string, args ...any) {
	l.With()
	l.Infof(format, args...)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.SugaredLogger.With(args...)}
}

func NoOpLogger() *Logger {
	return &Logger{zap.NewNop().Sugar()}
}
