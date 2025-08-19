package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub/v2"
)

const (
	SendEventTopic = "send-event-topic"
)

type application struct {
	logger       *slog.Logger
	pubSubClient *pubsub.Client
}

type EventData struct {
	Email string `json:"email"`
}

type PubSubEvent struct {
	Subscription string `json:"subscription"`
	Message      struct {
		Data string `json:"data"`
	} `json:"message"`
}

func (app application) sendEvent(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	topic := app.pubSubClient.Publisher(SendEventTopic)

	eventData := EventData{
		Email: "testing@pubsub.com",
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		app.logger.Error("error encoding event data", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}

	result := topic.Publish(c, &pubsub.Message{
		Data: data,
	})

	ID, err := result.Get(c)
	if err != nil {
		app.logger.Error("error sending pubsub event", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(fmt.Appendf(nil, "{\"status\": \"created\", serverID: %s}", ID))
}

func (app application) routes() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("POST /save", app.sendEvent)
	return r
}

func main() {
	var gcpProjectID string
	flag.StringVar(&gcpProjectID, "gcp-project-id", "demo-test", "gcp project id")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pubSubClient, err := pubsub.NewClient(context.Background(), gcpProjectID)
	if err != nil {
		logger.Error("error creating pubsub client")
		os.Exit(1)
	}
	defer pubSubClient.Close()

	app := application{
		logger,
		pubSubClient,
	}

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: app.routes(),
	}

	logger.Info("Starting server at port 8080")
	if err := server.ListenAndServe(); err != nil {
		logger.Error("error starting server", slog.String("error", err.Error()))
	}
}
