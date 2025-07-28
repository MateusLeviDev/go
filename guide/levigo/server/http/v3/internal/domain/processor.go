package domain

import (
	"context"
)

type ProcessorService interface {
	ProcessPayment(ctx context.Context, payment Payment) (string, error)
}
