package domain

import "context"

type PaymentRepository interface {
	AddToStream(payment Payment) error
	StorePayment(ctx context.Context, payment Payment) error
	GetAllPayments(ctx context.Context) ([]Payment, error)
}
