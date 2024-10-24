package config

import (
	"context"
	"os/signal"
	"syscall"
	"time"
)

type Bootstrap struct {
	Ctx context.Context
}

func Run(configBoot *Bootstrap) error {
	ctx, cancel := signal.NotifyContext(configBoot.Ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Set up logging
	log := NewLogger()

	// set up viper
	viperCustom := NewViper(log)

	// Set up database
	// db := NewDatabasePostgres(viperCustom, log)

	server, err := AppRunServe(ctx, viperCustom, log)
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
