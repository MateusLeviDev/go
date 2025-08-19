package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/MateusLeviDev/internal/pubsub"
)

type Handler struct {
	logger    *slog.Logger
	publisher *pubsub.Publisher
}

func NewHandler(logger *slog.Logger, publisher *pubsub.Publisher) *Handler {
	return &Handler{logger, publisher}
}

func (h *Handler) SendEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	event := pubsub.EventData{
		Email: "testing@pubsub.com",
	}

	id, err := h.publisher.PublishEvent(ctx, event)
	if err != nil {
		h.logger.Error("error sending pubsub event", slog.String("error", err.Error()))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(fmt.Appendf(nil, `{"status":"created","serverID":"%s"}`, id))
}
