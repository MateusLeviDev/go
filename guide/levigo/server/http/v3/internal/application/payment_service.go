package application

import (
	"context"

	"github.com/MateusLeviDev/internal/domain"
)

type PaymentService struct {
	Repo      domain.PaymentRepository
	Processor domain.ProcessorService
}

func (s *PaymentService) SubmitPayment(ctx context.Context, payment domain.Payment) error {
	return s.Repo.AddToStream(payment)
}
