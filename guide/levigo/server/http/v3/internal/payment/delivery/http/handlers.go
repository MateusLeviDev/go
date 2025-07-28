package payment

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MateusLeviDev/config"
	"github.com/MateusLeviDev/internal/application"
	"github.com/MateusLeviDev/internal/domain"
	"github.com/google/uuid"
)

type Handler struct {
	cfg        *config.Config
	PaymentSvc *application.PaymentService
	SummarySvc *application.SummaryService
}

type paymentRequest struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func NewHandler(paymentSvc *application.PaymentService, summarySvc *application.SummaryService, cfg *config.Config) *Handler {
	return &Handler{
		PaymentSvc: paymentSvc,
		SummarySvc: summarySvc,
		cfg:        cfg,
	}
}

func (h *Handler) HandlePayments(w http.ResponseWriter, r *http.Request) {
	var req paymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[Handler] Invalid body: %v", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(req.CorrelationId)
	if err != nil {
		log.Printf("[Handler] Invalid correlationId: %v", err)
		http.Error(w, "invalid correlationId", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		log.Printf("[Handler] Invalid amount: %v", req.Amount)
		http.Error(w, "amount must be > 0", http.StatusBadRequest)
		return
	}
	payment := domain.Payment{
		CorrelationId: id,
		Amount:        req.Amount,
		RequestedAt:   time.Now().UTC(),
	}

	if err := h.PaymentSvc.SubmitPayment(r.Context(), payment); err != nil {
		log.Printf("[Handler] Failed to enqueue payment: %v", err)
		http.Error(w, "failed to enqueue payment", http.StatusInternalServerError)
		return
	}

	log.Printf("[Handler] Payment enqueued: %s", req.CorrelationId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response := map[string]string{
		"status":        "success",
		"message":       "Payment request accepted",
		"correlationId": req.CorrelationId,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleSummary(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	var from, to *time.Time
	if fromStr != "" {
		f, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &f
		}
	}
	if toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}
	summary, err := h.SummarySvc.GetSummary(from, to)
	if err != nil {
		log.Printf("[Handler] Failed to get summary: %v", err)
		http.Error(w, "failed to get summary", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		log.Printf("[Handler] Failed to encode summary: %v", err)
		http.Error(w, "failed to encode summary", http.StatusInternalServerError)
	}
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
