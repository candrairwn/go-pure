package config

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/candrairwn/go-pure/api/delivery/http/routes"
)

type Bootstrap struct {
	Ctx    context.Context
	Stdout io.Writer
	Args   []string
}

func Run(configBoot *Bootstrap) error {
	ctx, cancel := signal.NotifyContext(configBoot.Ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Parse flags
	var port uint
	fs := flag.NewFlagSet(configBoot.Args[0], flag.ExitOnError)
	fs.SetOutput(configBoot.Stdout)
	fs.UintVar(&port, "port", 80, "port to listen on")
	if err := fs.Parse(configBoot.Args[1:]); err != nil {
		return err
	}

	// Set up logging
	// slog.SetDefault(slog.New(slog.NewJSONHandler(configBoot.Stdout, nil)))
	log := NewLogger()

	// Start the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: routes.Route(log, "1.0"),
	}

	// Start the server
	errChan := make(chan error, 1)
	go func() {
		log.Info("starting server on port " + server.Addr)
		// slog.InfoContext(ctx, "starting", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}

	}()

	// Wait for a signal
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		slog.InfoContext(ctx, "shutting down")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
