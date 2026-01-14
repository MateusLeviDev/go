package http

import (
	"encoding/json"
	"fn-contract-settled/internal/models"
	"fn-contract-settled/internal/ports"
	"net/http"

	"github.com/google/uuid"
)

type ContractServer struct {
	producer ports.ContractProducer
}

func NewContractServer(producer ports.ContractProducer) *ContractServer {
	return &ContractServer{producer: producer}
}

func (s *ContractServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contracts/settled" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var payload models.ContractSettled
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	correlationID := r.Header.Get("X-Correlation-Id")
	if correlationID == "" {
		correlationID = uuid.NewString()
	}

	if err := s.producer.ProduceContractSettled(payload, correlationID); err != nil {
		http.Error(w, "failed to publish event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
