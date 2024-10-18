package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func gracefulShutdown() (ctx context.Context) {
	// Create a context that will be canceled on shutdown
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())

	// Set up a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// Block until a signal is received
		sig := <-sigChan
		slog.Info("Received signal to shutdown", slog.String("signal", sig.String()))
		// Trigger graceful shutdown by canceling the context
		cancel()
	}()

	return
}
