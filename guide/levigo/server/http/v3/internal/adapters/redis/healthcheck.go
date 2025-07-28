package redis

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProcessorStatus struct {
	URL     string
	Service string
}

type HealthCheckService struct {
	RedisClient *redis.Client
	Healthy     atomic.Value // ProcessorStatus
}

func NewHealthCheckService(client *redis.Client) *HealthCheckService {
	h := &HealthCheckService{RedisClient: client}
	h.Healthy.Store(ProcessorStatus{URL: os.Getenv("PROCESSOR_DEFAULT_URL"), Service: "default"})
	return h
}

func (h *HealthCheckService) Start() {
	go func() {
		ticker := time.NewTicker(6 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if h.acquireLock() {
				h.updateHealthyProcessor()
				h.releaseLock()
			} else {
				h.readHealthyFromRedis()
			}
		}
	}()
}

func (h *HealthCheckService) acquireLock() bool {
	ctx := context.Background()
	ok, _ := h.RedisClient.SetNX(ctx, "health_check_lock", "locked", 10*time.Second).Result()
	return ok
}

func (h *HealthCheckService) releaseLock() {
	h.RedisClient.Del(context.Background(), "health_check_lock")
}

func (h *HealthCheckService) updateHealthyProcessor() {
	defaultURL := os.Getenv("PROCESSOR_DEFAULT_URL")
	fallbackURL := os.Getenv("PROCESSOR_FALLBACK_URL")
	if h.isHealthy(defaultURL + "/service-health") {
		h.Healthy.Store(ProcessorStatus{URL: defaultURL + "/payments", Service: "default"})
		h.saveHealthyToRedis("default")
		return
	}
	if h.isHealthy(fallbackURL + "/service-health") {
		h.Healthy.Store(ProcessorStatus{URL: fallbackURL + "/payments", Service: "fallback"})
		h.saveHealthyToRedis("fallback")
		return
	}
	cur := h.Healthy.Load().(ProcessorStatus)
	h.saveHealthyToRedis(cur.Service)
}

func (h *HealthCheckService) isHealthy(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var res struct{ Failing bool }
	_ = json.NewDecoder(resp.Body).Decode(&res)
	return !res.Failing
}

func (h *HealthCheckService) saveHealthyToRedis(service string) {
	ctx := context.Background()
	h.RedisClient.HMSet(ctx, "healthy_processor_status", map[string]interface{}{
		"service":   service,
		"timestamp": time.Now().Unix(),
	})
}

func (h *HealthCheckService) readHealthyFromRedis() {
	ctx := context.Background()
	m, _ := h.RedisClient.HGetAll(ctx, "healthy_processor_status").Result()
	if m["service"] == "default" {
		h.Healthy.Store(ProcessorStatus{URL: os.Getenv("PROCESSOR_DEFAULT_URL") + "/payments", Service: "default"})
	} else if m["service"] == "fallback" {
		h.Healthy.Store(ProcessorStatus{URL: os.Getenv("PROCESSOR_FALLBACK_URL") + "/payments", Service: "fallback"})
	}
}

func (h *HealthCheckService) GetCurrent() ProcessorStatus {
	return h.Healthy.Load().(ProcessorStatus)
}
