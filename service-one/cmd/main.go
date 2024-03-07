package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shandler/go-expert-observabilidade/service-one/internal/core/search"
	"github.com/shandler/go-expert-observabilidade/service-one/internal/infra"
	"github.com/shandler/go-expert-observabilidade/shared"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := shared.InitProvider(
		viper.GetString("OTEL_SERVICE_NAME"),
		viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	config := shared.NewConfig()

	clientZipCode := search.New(http.DefaultClient, config)
	server := infra.NewServer(config, clientZipCode)
	router := server.CreateServer()

	go func() {
		log.Println("Starting server on port", viper.GetString("HTTP_PORT"))
		if err := http.ListenAndServe(viper.GetString("HTTP_PORT"), router); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	// Create a timeout context for the graceful shutdown
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
