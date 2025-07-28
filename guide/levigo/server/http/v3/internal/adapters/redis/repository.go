package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MateusLeviDev/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PaymentRepository struct {
	Client *redis.Client
}

func NewPaymentRepository(client *redis.Client) *PaymentRepository {
	return &PaymentRepository{Client: client}
}

func (r *PaymentRepository) AddToStream(payment domain.Payment) error {
	data, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}
	return r.Client.LPush(context.Background(), "payments:queue", data).Err()
}

func (r *PaymentRepository) StorePayment(ctx context.Context, payment domain.Payment) error {
	paymentData := map[string]interface{}{
		"correlationId": payment.CorrelationId.String(),
		"amount":        payment.Amount,
		"processor":     payment.Processor,
		"requestedAt":   payment.RequestedAt.Format(time.RFC3339Nano),
	}
	paymentJSON, err := json.Marshal(paymentData)
	if err != nil {
		return fmt.Errorf("failed to marshal payment data: %w", err)
	}
	return r.Client.HSet(ctx, "payments", payment.CorrelationId.String(), paymentJSON).Err()
}

func (r *PaymentRepository) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	paymentsData, err := r.Client.HGetAll(ctx, "payments").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payments: %w", err)
	}
	var payments []domain.Payment
	for _, paymentDataJSON := range paymentsData {
		var paymentData map[string]interface{}
		if err := json.Unmarshal([]byte(paymentDataJSON), &paymentData); err != nil {
			continue
		}
		payment, err := parsePaymentFromData(paymentData)
		if err != nil {
			continue
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func parsePaymentFromData(data map[string]interface{}) (domain.Payment, error) {
	correlationIdStr, ok := data["correlationId"].(string)
	if !ok {
		return domain.Payment{}, fmt.Errorf("invalid correlationId")
	}
	amount, ok := data["amount"].(float64)
	if !ok {
		return domain.Payment{}, fmt.Errorf("invalid amount")
	}
	processor, ok := data["processor"].(string)
	if !ok {
		return domain.Payment{}, fmt.Errorf("invalid processor")
	}
	requestedAtStr, ok := data["requestedAt"].(string)
	if !ok {
		return domain.Payment{}, fmt.Errorf("invalid requestedAt")
	}
	correlationId, err := uuid.Parse(correlationIdStr)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("failed to parse correlationId: %w", err)
	}
	requestedAt, err := time.Parse(time.RFC3339Nano, requestedAtStr)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("failed to parse requestedAt: %w", err)
	}
	return domain.Payment{
		CorrelationId: correlationId,
		Amount:        amount,
		RequestedAt:   requestedAt.UTC(),
		Processor:     processor,
	}, nil
}
