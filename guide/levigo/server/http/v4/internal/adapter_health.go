package internal

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/MateusLeviDev/pkg/constants"
	"github.com/bytedance/sonic"
)

func (a *PaymentProcessorAdapter) EnableHealthCheck(should string) {
	if should != "true" {
		return
	}

	go func() {
		ticker := time.NewTicker(constants.HealthCheckTicker)
		defer ticker.Stop()

		for range ticker.C {
			if err := a.storeHealthStatus(a.defaultUrl+constants.PaymentServiceHealthUrl, constants.HealthCheckKeyDefault); err != nil {
				slog.Debug("failed to update the health check", "err", err)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(constants.HealthCheckTicker)
		defer ticker.Stop()

		for range ticker.C {
			if err := a.storeHealthStatus(a.fallbackUrl+constants.PaymentServiceHealthUrl, constants.HealthCheckKeyFallback); err != nil {
				slog.Debug("failed to update the health check", "err", err)
			}
		}
	}()
}

func (a *PaymentProcessorAdapter) storeHealthStatus(url string, key string) error {
	resDefault, err := a.retrieveHealth(url)
	if err != nil {
		return err
	}

	reqbody := HealthCheckResponse{
		Failing:         resDefault.Failing,
		MinResponseTime: resDefault.MinResponseTime,
	}
	rawBody, err := sonic.Marshal(reqbody)
	if err != nil {
		slog.Debug("failed to encode the json object for redis", "err", err)
		return err
	}

	if err := a.db.Set(context.Background(), key, rawBody, 0).Err(); err != nil {
		slog.Debug("failed to save health check in redis", "err", err)
		return err
	}

	slog.Debug("updating the health check", "healthCheckStatus", reqbody, "key", key)
	return nil
}

func (a *PaymentProcessorAdapter) retrieveHealth(url string) (HealthCheckResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return HealthCheckResponse{}, err
	}

	res, err := a.client.Do(req)
	if res == nil || err != nil || res.StatusCode != 200 {
		slog.Debug("failed to health check", "url", url)
		return HealthCheckResponse{}, err
	}

	var respBody HealthCheckResponse
	decoder := sonic.ConfigFastest.NewDecoder(res.Body)
	if err := decoder.Decode(&respBody); err != nil {
		slog.Debug("failed to parse the response", "url", url)
		return HealthCheckResponse{}, err
	}

	return respBody, nil
}

func (a *PaymentProcessorAdapter) StartWorkers() {
	for range a.workers {
		go a.retryWorkers()
	}

	go func() {
		for {
			slog.Debug("Status of queue", "lenRetryQueue", len(a.retryQueue))
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		ticker := time.NewTicker(constants.HealthCheckTicker)
		defer ticker.Stop()

		for range ticker.C {
			if err := a.syncHealthStatus(constants.HealthCheckKeyDefault); err != nil {
				slog.Debug("failed update the health check", "err", err)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(constants.HealthCheckTicker)
		defer ticker.Stop()

		for range ticker.C {
			if err := a.syncHealthStatus(constants.HealthCheckKeyFallback); err != nil {
				slog.Debug("failed update the health check", "err", err)
			}
		}
	}()
}

func (a *PaymentProcessorAdapter) syncHealthStatus(key string) error {
	resBody, err := a.db.Get(context.Background(), key).Result()
	if err != nil {
		slog.Debug("failed to get the health check", "err", err)
		return err
	}

	var healthCheckStatus HealthCheckResponse
	if err := sonic.ConfigFastest.Unmarshal([]byte(resBody), &healthCheckStatus); err != nil {
		slog.Debug("failed to unmarshal the health check from redis", "err", err)
		return err
	}

	switch key {
	case constants.HealthCheckKeyDefault:
		a.healthStatusDefault.Store(healthCheckStatus)
	case constants.HealthCheckKeyFallback:
		a.healthStatusFallback.Store(healthCheckStatus)
	}

	return nil
}

func (a *PaymentProcessorAdapter) retryWorkers() {
	for payment := range a.retryQueue {
		time.Sleep(time.Millisecond * 10) // wait before retry
		a.Process(payment)
	}
}
