package main

import (
	"flag"
	"fn-contract-settled/config"
	kafka "fn-contract-settled/internal/adapters"
	"log"
	"net/http"

	httpServer "fn-contract-settled/internal/http"
)

func main() {
	log.Println("Starting microservice")

	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	producer, err := kafka.NewProducer(cfg.Producer)
	if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	server := httpServer.NewContractServer(producer)
	log.Printf("HTTP server running on %s\n", ":8080")

	log.Fatal(http.ListenAndServe(":8080", server))
}
