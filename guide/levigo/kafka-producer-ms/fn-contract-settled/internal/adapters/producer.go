package kafka

import (
	"encoding/json"
	"fmt"
	"fn-contract-settled/config"
	"fn-contract-settled/internal/events"
	"fn-contract-settled/internal/models"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer   *kafka.Producer
	topic      string
	deliveryCh chan kafka.Event
}

func NewProducer(cfg config.Producer) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.BootstrapServers,
		"client.id":          cfg.ClientID,
		"acks":               "all",
		"linger.ms":          5,
		"retries":            5,
		"enable.idempotence": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	return &Producer{
		producer:   p,
		topic:      cfg.Topic,
		deliveryCh: make(chan kafka.Event, 100),
	}, nil
}

func (p *Producer) ProduceContractSettled(
	eventData models.ContractSettled,
	correlationID string,
) error {
	event, err := events.NewEvent(
		events.MetaData{
			Source:  "fn-contract-settled",
			Type:    "contract.settled",
			Subject: "contrcat/" + eventData.ContractID,
		},
		eventData,
	)
	if err != nil {
		return nil
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.produce(
		payload,
		eventData.ContractID,
		correlationID,
	)
}

func (p *Producer) produce(
	payload []byte,
	key string,
	correlationID string,
) error {

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "correlation-id", Value: []byte(correlationID)},
		},
	}, p.deliveryCh)

	if err != nil {
		return err
	}

	return p.waitForDeliveryReport()
}

func (p *Producer) waitForDeliveryReport() error {
	e := <-p.deliveryCh
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}
	return nil
}

func (p *Producer) Close() {
	p.producer.Flush(5000)
	p.producer.Close()
	close(p.deliveryCh)
}
