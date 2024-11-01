package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewViper(log *zap.SugaredLogger) (*viper.Viper, error) {
	viper := viper.New()

	viper.AutomaticEnv()
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to read config")
		return nil, err
	}

	log.Info("success read config")

	return viper, nil
}
