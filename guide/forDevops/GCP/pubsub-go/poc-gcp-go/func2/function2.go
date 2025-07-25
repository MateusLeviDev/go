package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

const maxRetries = 3

func main() {
	ctx := context.Background()

	// Configurações do cliente Pub/Sub
	projectId := os.Getenv("GCP_PROJECT_ID")
	emulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")

	subscriptionID := "example-subscription3"

	// Criando o cliente Pub/Sub
	client, err := pubsub.NewClient(ctx, projectId, option.WithEndpoint(emulatorHost))
	if err != nil {
		log.Fatalf("Erro ao criar o cliente Pub/Sub: %v", err)
	}
	defer client.Close()

	// Obtendo a subscription
	subscription := client.Subscription(subscriptionID)

	// Canal para tratar sinais do sistema
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	var mu sync.Mutex
	received := 0

	fmt.Println("Consumindo mensagens...")

	// Função de callback para processamento de mensagens
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-sigchan
		fmt.Println("Recebido sinal de interrupção, cancelando...")
		cancel()
	}()

	url := "http://localhost:3000/func2"

	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		received++
		messageID := msg.ID
		fmt.Printf("Mensagem recebida: %s, ID: %s\n", string(msg.Data), messageID)

		// Fazendo POST com a mensagem recebida
		success := false
		for i := 0; i < maxRetries; i++ {
			err := postMessage(url, msg.Data, messageID)
			if err != nil {
				fmt.Printf("Erro ao fazer o POST (tentativa %d), ID: %s: %v\n", i+1, messageID, err)
				time.Sleep(2 * time.Second) // Espera antes de tentar novamente
			} else {
				fmt.Println("POST realizado com sucesso, ID:", messageID)
				success = true
				break
			}
		}
		if success {
			msg.Ack()
			fmt.Printf("Confirmando mensagem (Ack), ID: %s...\n", messageID)
		} else {
			fmt.Printf("Falha ao processar a mensagem após várias tentativas, ID: %s\n", messageID)
		}
	})
	if err != nil {
		log.Fatalf("Erro ao consumir mensagens: %v", err)
	}

	fmt.Printf("Recebidas %d mensagens\n", received)
}

// Função para fazer POST com a mensagem recebida
func postMessage(url string, message []byte, messageID string) error {
	// Criação do payload
	payload := map[string]string{
		"message": string(message),
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao criar o payload JSON, ID: %s: %v", messageID, err)
	}

	// Fazendo a requisição POST
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("erro ao fazer a requisição POST, ID: %s: %v", messageID, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler o corpo da resposta, ID: %s: %v", messageID, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("recebido código de status %v, ID: %s, resposta: %s", resp.StatusCode, messageID, string(body))
	}

	fmt.Printf("Corpo da resposta, ID: %s: %s\n", messageID, string(body))
	return nil
}
