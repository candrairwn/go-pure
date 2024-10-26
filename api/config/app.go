package config

import (
	"context"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func AppRunServe(ctx context.Context, server *http.Server, viperCustom *viper.Viper, log *zap.SugaredLogger) (*http.Server, error) {
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
