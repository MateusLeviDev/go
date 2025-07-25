package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type Message struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type PubsubService struct {
	projectId string
	client    *pubsub.Client
}

func NewPubsubService() (*PubsubService, error) {
	projectId := os.Getenv("GCP_PROJECT_ID")
	emulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")

	client, err := pubsub.NewClient(context.Background(), projectId, option.WithEndpoint(emulatorHost))
	if err != nil {
		return nil, err
	}

	return &PubsubService{
		projectId: projectId,
		client:    client,
	}, nil
}

func (s *PubsubService) PublisherMessage(topicID, message string) (string, error) {
	ctx := context.Background()
	topic := s.client.Topic(topicID)
	defer topic.Stop()

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(message),
	})

	id, err := result.Get(ctx)
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("Published message ID: %s", id)
	fmt.Println(msg)
	return msg, nil
}

func main() {
	// URL do endpoint
	url := "http://localhost:3000/posts"

	// Fazendo a requisição GET
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erro ao fazer a requisição: %v", err)
	}
	defer resp.Body.Close()

	// Lendo o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler o corpo da resposta: %v", err)
	}

	// Parse do JSON recebido
	var messages []Message
	if err := json.Unmarshal(body, &messages); err != nil {
		log.Fatalf("Erro ao fazer parse do JSON: %v", err)
	}

	service, err := NewPubsubService()
	if err != nil {
		fmt.Println("Error creating service:", err)
		return
	}

	topicID := "example-topic3"

	// Publicando cada mensagem individualmente
	for _, msg := range messages {
		messageJSON, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Erro ao converter mensagem para JSON: %v", err)
			continue
		}

		_, err = service.PublisherMessage(topicID, string(messageJSON))
		if err != nil {
			fmt.Println("Error publishing message:", err)
			return
		}
	}
}
