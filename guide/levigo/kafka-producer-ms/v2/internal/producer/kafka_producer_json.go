package producer

import (
	"fmt"

	"github.com/MateusLeviDev/internal/models"
	"github.com/MateusLeviDev/pkg/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pquerna/ffjson/ffjson"
)

type Producer struct {
	producer   *kafka.Producer
	topic      string
	deliveryCh chan kafka.Event
}

func NewProducer(bootstrapServers, topic string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "myProducer",
		"acks":              "all"})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	logger.InfoAsync("Kafka Producer created successfully")
	return &Producer{
		producer:   p,
		topic:      topic,
		deliveryCh: make(chan kafka.Event, 100),
	}, nil
}

func (p *Producer) ProduceBatch(users []models.User, correlationID string) error {
	logger.InfoAsync("Starting Batch Production")
	for _, u := range users {
		payload, err := ffjson.Marshal(&u)
		if err != nil {
			logger.ErrorAsync("Failed to serialize payload:", err)
			return fmt.Errorf("failed to serialize payload: %w", err)
		}
		err = p.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
			Value:          payload,
			Headers:        []kafka.Header{{Key: "correlation-id", Value: []byte(correlationID)}},
		}, p.deliveryCh)
		if err != nil {
			logger.ErrorAsync("Produce failed:", err)
			return fmt.Errorf("produce failed: %w", err)
		}
	}
	// Wait for all delivery reports
	for range users {
		if err := p.waitForDeliveryReport(); err != nil {
			return err
		}
	}
	logger.InfoAsync("Batch production completed")
	return nil
}

func (p *Producer) waitForDeliveryReport() error {
	e := <-p.deliveryCh
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		logger.ErrorAsync("Delivery failed:", m.TopicPartition.Error)
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}

	// logger.InfoAsync(fmt.Sprintf("Delivered message to topic %s [%d] at offset %v",
	// 	*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset))
	return nil
}

func (p *Producer) Close() {
	logger.InfoAsync("Closing producer")
	close(p.deliveryCh)
	p.producer.Close()
}
