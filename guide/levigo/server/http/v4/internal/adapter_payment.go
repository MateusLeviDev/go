package internal

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRetriesAreOver       = errors.New("retries are over")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrUnavailableProcessor = errors.New("unavailable processor")
)

type PaymentProcessorAdapter struct {
	client               *http.Client
	db                   *redis.Client
	repo                 *PaymentRepository
	healthStatusDefault  atomic.Value
	healthStatusFallback atomic.Value
	defaultUrl           string
	fallbackUrl          string
	retryQueue           chan PaymentRequestProcessor
	workers              int
}

func NewPaymentProcessorAdapter(
	client *http.Client,
	db *redis.Client,
	repo *PaymentRepository,
	defaultUrl string,
	fallbackUrl string,
	retryQueue chan PaymentRequestProcessor,
	workers int,
) *PaymentProcessorAdapter {
	a := &PaymentProcessorAdapter{
		client:      client,
		db:          db,
		repo:        repo,
		defaultUrl:  defaultUrl,
		fallbackUrl: fallbackUrl,
		retryQueue:  retryQueue,
		workers:     workers,
	}

	a.healthStatusDefault.Store(HealthCheckResponse{
		Failing:         false,
		MinResponseTime: 0,
	})
	a.healthStatusFallback.Store(HealthCheckResponse{
		Failing:         false,
		MinResponseTime: 0,
	})

	return a
}

func (a *PaymentProcessorAdapter) Process(payment PaymentRequestProcessor) {
	err := a.innerProcess(payment)
	if err != nil {
		a.retryQueue <- payment
	}
}

func (a *PaymentProcessorAdapter) innerProcess(payment PaymentRequestProcessor) error {
	healthStatusDefault := a.healthStatusDefault.Load().(HealthCheckResponse)
	healthStatusFallback := a.healthStatusFallback.Load().(HealthCheckResponse)

	var err error
	if !healthStatusDefault.Failing && healthStatusDefault.MinResponseTime < 80 {
		err = a.sendPayment(
			payment,
			a.defaultUrl+"/payments",
			time.Second*10,
			PaymentEndpointDefault,
		)
	} else if !healthStatusFallback.Failing && healthStatusFallback.MinResponseTime < 80 {
		err = a.sendPayment(
			payment,
			a.fallbackUrl+"/payments",
			time.Second*10,
			PaymentEndpointFallback,
		)
	} else {
		return ErrUnavailableProcessor
	}

	if errors.Is(err, ErrInvalidRequest) {
		return nil
	}

	return err
}

func (a *PaymentProcessorAdapter) sendPayment(
	payment PaymentRequestProcessor,
	url string,
	timeout time.Duration,
	endpoint PaymentEndpoint,
) error {
	slog.Debug("sending the request", "body", payment, "url", url)
	start1 := time.Now()

	payment.UpdateRequestTime()
	reqBody, err := sonic.ConfigFastest.Marshal(payment)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")

	res, err := a.client.Do(req)
	slog.Debug("response from api", "url", url, "res", res, "payment", payment)

	if res != nil && res.StatusCode != 200 {
		return ErrUnavailableProcessor
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return ErrUnavailableProcessor
	}
	if err != nil || res == nil {
		slog.Debug("failed to process the request", "err", err, "res", res)
		return ErrUnavailableProcessor
	}

	start2 := time.Now()
	err = a.repo.Add(payment, endpoint)
	if time.Since(start1).Milliseconds() > 25 {
		slog.Debug("time of the complete request and db",
			"dbTimeMs", time.Since(start2).Milliseconds(),
			"requestTimeMs", time.Since(start1).Milliseconds(),
			"healthStatusDefault", a.healthStatusDefault.Load().(HealthCheckResponse),
			"healthStatusFallback", a.healthStatusFallback.Load().(HealthCheckResponse),
			"endpoint", endpoint,
			"err", err,
			"requestAt", *payment.RequestedAt,
		)
	}
	return err
}

func (a *PaymentProcessorAdapter) Summary(from, to string) (SummaryResponse, error) {
	return a.repo.Summary(from, to)
}

func (a *PaymentProcessorAdapter) Purge(token string) error {
	if err := a.repo.Purge(); err != nil {
		return err
	}
	if err := a.purge(a.defaultUrl+"/admin/purge-payments", token); err != nil {
		return err
	}
	if err := a.purge(a.fallbackUrl+"/admin/purge-payments", token); err != nil {
		return err
	}

	return nil
}

func (a *PaymentProcessorAdapter) purge(url string, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Rinha-Token", token)

	res, err := a.client.Do(req)
	if err != nil {
		slog.Error("failed to purge the api", "error", err, "url", url)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return ErrInvalidRequest
	}

	return nil
}
