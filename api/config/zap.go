package config

import "go.uber.org/zap"

func NewLogger() *zap.SugaredLogger {
	log := zap.NewExample().Sugar()
	defer log.Sync()

	return log
}
