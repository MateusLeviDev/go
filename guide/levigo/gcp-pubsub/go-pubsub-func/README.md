# GCP event driven rest API

This project demonstrates a "Hello, World!" example of triggering a Cloud Function via Pub/Sub when a REST API endpoint is called, focusing on local testing using emulators.


### Run PubSub emulator 

`$ gcloud beta emulators pubsub start`

* Note: This starts the emulator. You'll need to keep this running in a separate terminal window.

### Crete the topic 

`$ curl -X PUT http://localhost:8085/v1/projects/demo-test/topics/send-event-topic`

* Title: This command creates the Pub/Sub topic named "send-event-topic" in the emulator.
* Note: The demo-test project ID is arbitrary for the emulator.

### Check the Topic (Optional):

`$ curl http://localhost:8085/v1/projects/demo-test/topics`

* Title: This command lists the existing topics in the emulator, allowing you to verify that the topic was created successfully.


### Associate the topic with the Push endpoint:

```
curl -X PUT \
  http://localhost:8085/v1/projects/demo-test/subscriptions/my-subscription \
  -H "Content-Type: application/json" \
  -d '{
      "topic": "projects/demo-test/topics/send-event-topic",
      "pushConfig": {
        "pushEndpoint": "http://localhost:8082/projects/demo-test/topics/send-event-topic"
      }
  }'
```


### Check subscriptions

`$ curl http://localhost:8085/v1/projects/demo-test/subscriptions`

### Start the cloud event function

Navifgate to `/cmd/functions` and run

`$ go run main.go`

You should see the message:
2025/08/02 18:13:23 Listening on port 8082


### Call the REST API to trigger the entire flow:

`curl --location --request POST 'http://localhost:8080/save'`


### Payload to test the cloud event function
```
{
  "subscription": "projects/my-project/subscriptions/my-subscription",
  "message": {
    "@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
    "attributes": {
      "attr1":"attr1-value"
    },
    "data": "dGVzdCBtZXNzYWdlIDM=",
    "messageId": "message-id",
    "publishTime":"2021-02-05T04:06:14.109Z"
  }
}

```

* the request to test the function
```
curl -X POST http://localhost:8082/projects/demo-test/topics/send-event-topic \
  -H "Content-Type: application/json" \
  -H "ce-id: test-123" \
  -H "ce-source: test-source" \
  -H "ce-specversion: 1.0" \
  -H "ce-type: type.googleapis.com/google.pubusb.v1.PubsubMessage" \
  -d '{
  "subscription": "projects/my-project/subscriptions/my-subscription",
  "message": {
    "@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
    "attributes": {
      "attr1":"attr1-value"
    },
    "data": "eyJlbWFpbCI6ImNvbnRhY3RAeWFob28uY29tIn0=",
    "messageId": "message-id",
    "publishTime":"2021-02-05T04:06:14.109Z"
  }
}'
```

---

```
HTTP server com graceful shutdown, Pub/Sub Publisher com Drain via WaitGroup, closing gate no handler, evitar Close() duplicado, e aplicar timeouts no servidor.


package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
)

// Publisher com suporte a Drain
type Publisher struct {
	client *pubsub.Client
	topic  *pubsub.Topic

	wg      sync.WaitGroup
	closing chan struct{}
}

func NewPublisher(ctx context.Context, projectID, topicID string) (*Publisher, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		client:  client,
		topic:   client.Topic(topicID),
		closing: make(chan struct{}),
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, msg *pubsub.Message) error {
	select {
	case <-p.closing:
		// Rejeita novas mensagens após início do shutdown
		return fmt.Errorf("publisher is shutting down")
	default:
	}

	p.wg.Add(1)
	res := p.topic.Publish(ctx, msg)

	// Tratar resultado de forma assíncrona
	go func() {
		defer p.wg.Done()
		_, err := res.Get(ctx)
		if err != nil {
			log.Printf("publish error: %v", err)
		}
	}()
	return nil
}

// Drain aguarda todas as publicações terminarem
func (p *Publisher) Drain() {
	log.Println("Publisher draining...")
	close(p.closing) // fecha o gate
	p.wg.Wait()      // espera terminar
	log.Println("Publisher drained.")
}

func (p *Publisher) Close() error {
	p.Drain()
	p.topic.Stop()
	return p.client.Close()
}

func main() {
	ctx := context.Background()

	// Cria publisher
	pub, err := NewPublisher(ctx, "meu-projeto", "minha-topic")
	if err != nil {
		log.Fatalf("falha ao criar publisher: %v", err)
	}

	// Handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		err := pub.Publish(r.Context(), &pubsub.Message{
			Data: []byte("hello world"),
		})
		if err != nil {
			http.Error(w, "falha ao publicar: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})

	// HTTP Server com timeouts
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	// Run server em goroutine
	go func() {
		log.Println("Servidor iniciado em :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("erro no servidor: %v", err)
		}
	}()

	// Espera sinal de interrupção
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Iniciando graceful shutdown...")

	// Timeout de shutdown total
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 1. Para de aceitar novas requisições
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("erro no shutdown do servidor: %v", err)
	}

	// 2. Fecha publisher após servidor encerrar
	if err := pub.Close(); err != nil {
		log.Printf("erro ao fechar publisher: %v", err)
	}

	log.Println("Shutdown completo.")
}
```
