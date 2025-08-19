package http

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /save", h.SendEvent)
	return mux
}
