package payment

import "net/http"

func (h *Handler) MapRoutes(mux *http.ServeMux) {
	mux.HandleFunc(h.cfg.Http.PaymentsPath, h.HandlePayments)
	mux.HandleFunc(h.cfg.Http.SummaryPath, h.HandleSummary)
	mux.HandleFunc(h.cfg.Http.HealthPath, h.HandleHealth)
}
