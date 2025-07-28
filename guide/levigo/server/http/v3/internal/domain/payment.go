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
	Default  SummaryItem
	Fallback SummaryItem
}

type SummaryItem struct {
	TotalRequests int
	TotalAmount   float64
}
