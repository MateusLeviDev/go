package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gcpPubsub "cloud.google.com/go/pubsub/v2"
	"github.com/MateusLeviDev/config"
	httpHandlers "github.com/MateusLeviDev/internal/http"
	internalPubsub "github.com/MateusLeviDev/internal/pubsub"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pubSubClient, err := gcpPubsub.NewClient(context.Background(), cfg.GCPProjectID)
	if err != nil {
		logger.Error("error creating pubsub client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pubSubClient.Close()

	publisher := internalPubsub.NewPublisher(pubSubClient, logger)
	handler := httpHandlers.NewHandler(logger, publisher)

	server := http.Server{
		Addr:    cfg.ServerAddr,
		Handler: handler.Routes(),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server", slog.String("addr", cfg.ServerAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.String("error", err.Error()))
		}
	}()

	<-stop
	logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("http server shutdown failed", slog.String("error", err.Error()))
	} else {
		logger.Info("http server stopped gracefully")
	}

	if err := pubSubClient.Close(); err != nil {
		logger.Error("error closing pubsub client", slog.String("error", err.Error()))
	} else {
		logger.Info("pubsub client closed gracefully")
	}
}
