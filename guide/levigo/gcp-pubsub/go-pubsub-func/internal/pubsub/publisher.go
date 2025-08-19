package pubsub

import (
	"context"
	"encoding/json"
	"log/slog"

	"cloud.google.com/go/pubsub/v2"
)

const SendEventTopic = "send-event-topic"

type Publisher struct {
	client *pubsub.Client
	logger *slog.Logger
}

func NewPublisher(client *pubsub.Client, logger *slog.Logger) *Publisher {
	return &Publisher{client, logger}
}

type EventData struct {
	Email string `json:"email"`
}

func (p *Publisher) PublishEvent(ctx context.Context, event EventData) (string, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	topic := p.client.Publisher(SendEventTopic)
	result := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	return result.Get(ctx)
}