package event

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	functions.CloudEvent("ProcessEvent", ProcessEvent)
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data []byte `json:"data"`
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

func ProcessEvent(ctx context.Context, e event.Event) error {
	logger.Info("receiving data", slog.String("id", e.ID()))

	var pse PubSubEvent
	if err := e.DataAs(&pse); err != nil {
		logger.Error("cannot read DataAs")
		return fmt.Errorf("event.DataAs: %v", err)
	}

	payload, err := base64.StdEncoding.DecodeString(pse.Message.Data)
	if err != nil {
		logger.Error("error decoding payload", slog.String("err", err.Error()))
		return err
	}

	var evtData EventData
	if err := json.Unmarshal(payload, &evtData); err != nil {
		logger.Error("cannot parse the email json unmarshalled", slog.String("error", err.Error()))
		return err
	}

	logger.Info("Hello", slog.String("email", evtData.Email))
	return nil
}
