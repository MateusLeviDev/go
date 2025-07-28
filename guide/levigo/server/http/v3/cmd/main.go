package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/MateusLeviDev/config"
	"github.com/MateusLeviDev/internal/server"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cfg)
	if err := server.NewServer(cfg).Run(); err != nil {
		log.Fatalf("Failed to start http server: %v", err)
	}
}
