package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MateusLeviDev/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Worker struct {
	Client    *redis.Client
	Health    *HealthCheckService
	Repo      *PaymentRepository
	WorkerNum int
}

func (w *Worker) Start(ctx context.Context) {
	processingQueue := fmt.Sprintf("payments:processing:%d", w.WorkerNum)
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
		},
	}
	for {
		result, err := w.Client.RPopLPush(ctx, "payments:queue", processingQueue).Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(1 * time.Second)
				continue
			}
			log.Printf("[Worker %d] Redis error: %v", w.WorkerNum, err)
			time.Sleep(1 * time.Second)
			continue
		}
		var payment domain.Payment
		if err := json.Unmarshal([]byte(result), &payment); err != nil {
			log.Printf("[Worker %d] Failed to unmarshal payment: %v", w.WorkerNum, err)
			w.Client.LRem(ctx, processingQueue, 1, result)
			continue
		}
		processor := w.Health.GetCurrent()
		if !w.processPayment(ctx, payment, processor, client) {
			w.Client.LPush(ctx, "payments:queue", result)
			w.Client.LRem(ctx, processingQueue, 1, result)
			continue
		}
		w.Client.LRem(ctx, processingQueue, 1, result)
	}
}

func (w *Worker) processPayment(ctx context.Context, payment domain.Payment, processor ProcessorStatus, client *http.Client) bool {
	if w.Repo == nil {
		log.Printf("[Worker %d] CRITICAL: Repo is nil!", w.WorkerNum)
		return false
	}
	if processor.URL == "" {
		log.Printf("[Worker %d] CRITICAL: Processor URL is empty!", w.WorkerNum)
		return false
	}
	body := map[string]interface{}{
		"correlationId": payment.CorrelationId.String(),
		"amount":        payment.Amount,
		"requestedAt":   payment.RequestedAt.Format(time.RFC3339Nano),
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", processor.URL, bytes.NewReader(b))
	if err != nil {
		log.Printf("[Worker %d] Failed to create request: %v", w.WorkerNum, err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Worker %d] Processor %s HTTP error: %v", w.WorkerNum, processor.Service, err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("[Worker %d] Processor %s returned status: %v", w.WorkerNum, processor.Service, resp.StatusCode)
		return false
	}
	payment.Processor = processor.Service
	if err := w.Repo.StorePayment(ctx, payment); err != nil {
		log.Printf("[Worker %d] CRITICAL: Payment accepted by processor but failed to save in Redis: %v", w.WorkerNum, err)
	}
	return true
}
