package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MateusLeviDev/config"
	"github.com/MateusLeviDev/internal/adapters/processor"
	"github.com/MateusLeviDev/internal/adapters/redis"
	"github.com/MateusLeviDev/internal/application"
	payment "github.com/MateusLeviDev/internal/payment/delivery/http"

	redisv9 "github.com/redis/go-redis/v9"
)

type server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *server {
	return &server{cfg: cfg}
}

func (s *server) Run() error {
	fmt.Printf("Listening on Port: %s\n", s.cfg.Http.Port)

	redisClient := redisv9.NewClient(&redisv9.Options{Addr: s.cfg.Redis.Address})
	repo := redis.NewPaymentRepository(redisClient)
	processorClient := processor.NewClient(s.cfg.Processor.DefaultURL, s.cfg.Processor.FallbackURL)
	paymentService := &application.PaymentService{Repo: repo, Processor: processorClient}
	summaryService := &application.SummaryService{Repo: repo}

	healthCheck := redis.NewHealthCheckService(redisClient)
	healthCheck.Start()
	for i := 0; i < s.cfg.Workers; i++ {
		worker := &redis.Worker{Client: redisClient, Health: healthCheck, Repo: repo, WorkerNum: i}
		go worker.Start(context.Background())
	}

	mux := http.NewServeMux()
	paymentHandlers := payment.NewHandler(paymentService, summaryService, s.cfg)
	paymentHandlers.MapRoutes(mux)

	return http.ListenAndServe(s.cfg.Http.Port, mux)
}
