package domain

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	CorrelationId uuid.UUID
	Amount        float64
	RequestedAt   time.Time
	Processor     string
}

type Summary struct {
	Default  SummaryItem `json:"default"`
	Fallback SummaryItem `json:"fallback"`
}

type SummaryItem struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}
