package application

import (
	"context"
	"time"

	"github.com/MateusLeviDev/internal/domain"
)

type SummaryService struct {
	Repo domain.PaymentRepository
}

func (s *SummaryService) GetSummary(from, to *time.Time) (domain.Summary, error) {
	ctx := context.Background()
	payments, err := s.Repo.GetAllPayments(ctx)
	if err != nil {
		return domain.Summary{}, err
	}
	var defaultCount, fallbackCount int
	var defaultAmount, fallbackAmount float64
	for _, p := range payments {
		if from != nil && !from.IsZero() && p.RequestedAt.Before(*from) {
			continue
		}
		if to != nil && !to.IsZero() && p.RequestedAt.After(*to) {
			continue
		}
		switch p.Processor {
		case "default":
			defaultCount++
			defaultAmount += p.Amount
		case "fallback":
			fallbackCount++
			fallbackAmount += p.Amount
		}
	}
	return domain.Summary{
		Default: domain.SummaryItem{
			TotalRequests: defaultCount,
			TotalAmount:   defaultAmount,
		},
		Fallback: domain.SummaryItem{
			TotalRequests: fallbackCount,
			TotalAmount:   fallbackAmount,
		},
	}, nil
}
