package config

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/candrairwn/go-pure/api/utils"
)

type Bootstrap struct {
	Ctx context.Context
}

func Run(configBoot *Bootstrap) error {
	// Set up context
	ctx, cancel := signal.NotifyContext(configBoot.Ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Set up logging
	log, err := NewLogger()
	if err != nil {
		return err
	}

	// set up viper
	viperCustom, err := NewViper(log)
	if err != nil {
		return err
	}

	// Set up database
	db, err := NewDatabasePostgres(viperCustom, log)
	if err != nil {
		return err
	}

	err = utils.NewJWTUtil(viperCustom, log).LoadFileKeys()
	if err != nil {
		return err
	}

	// Set Up Port default 80
	var port uint
	if viperCustom.GetUint("WEB_PORT") == 0 {
		port = 80
	} else {
		port = viperCustom.GetUint("WEB_PORT")
	}

	// Set up http mux
	Mux := http.NewServeMux()

	handlerBootstrap := BootsrapWire(&BootstrapWireConfig{
		App:   Mux,
		Log:   log,
		Viper: viperCustom,
		DB:    db,
	})

	// Create a new http server instance
	serverConfig := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handlerBootstrap,
	}

	// Run the server
	server, err := AppRunServe(ctx, serverConfig, viperCustom, log)
	if err != nil {
		return err
	}

	// Wait for the server to shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	return server.Shutdown(ctx)
}
