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

func Run(config *Bootstrap) error {
	ctx, cancel := signal.NotifyContext(config.Ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Parse flags
	var port uint
	fs := flag.NewFlagSet(config.Args[0], flag.ExitOnError)
	fs.SetOutput(config.Stdout)
	fs.UintVar(&port, "port", 80, "port to listen on")
	if err := fs.Parse(config.Args[1:]); err != nil {
		return err
	}

	// Set up logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(config.Stdout, nil)))

	// Start the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: routes.Route(slog.Default(), "1.0"),
	}

	// Start the server
	errChan := make(chan error, 1)
	go func() {
		slog.InfoContext(ctx, "starting", slog.String("addr", server.Addr))
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
