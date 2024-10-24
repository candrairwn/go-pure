package config

import (
	"context"
	"fmt"
	"net/http"

	"github.com/candrairwn/go-pure/api/delivery/http/routes"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func AppRunServe(ctx context.Context, viperCustom *viper.Viper, log *zap.SugaredLogger) (*http.Server, error) {
	// port configuration
	var port uint
	if viperCustom.GetUint("WEB_PORT") == 0 {
		port = 80
	} else {
		port = viperCustom.GetUint("WEB_PORT")
	}

	// Create a new http server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: routes.Route(log, viperCustom.GetString("APP_VERSION")),
	}

	// Run the server
	errChan := make(chan error, 1)
	go func() {
		log.Info("starting server on port " + server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for a signal
	select {
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		log.Info("shutting down", Context("context", ctx))
	}

	return server, nil
}
