package main

import (
	"context"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/MateusLeviDev/functions/event"
)

func main() {
	port := "8082"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.RegisterCloudEventFunctionContext(context.Background(), "/projects/demo-test/topics/send-event-topic", event.ProcessEvent); err != nil {
		log.Fatalf("Failed to register function: %v", err)
	}

	log.Printf("Listening on port %s", port)

	if err := funcframework.StartHostPort("localhost", port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v", err)
	}
}
