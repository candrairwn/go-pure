package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Valuer interface {
	Value(key interface{}) interface{}
}

func Context(key string, value Valuer) zap.Field {
	return zap.Object(key, valuer{value})
}

type (
	valuer struct {
		Valuer
	}

	contextKey struct{}
)

func (v valuer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	fs, _ := v.Value(contextKey{}).([]zap.Field)
	for i := range fs {
		fs[i].AddTo(enc)
	}

	return nil
}

func NewLogger() *zap.SugaredLogger {
	config := zap.Config{
		Encoding:         "json", // Set encoding to JSON
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	// Build the logger
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync() // Flushes buffer, if any

	// Create a SugaredLogger from the built logger
	sugar := logger.Sugar()

	return sugar
}
