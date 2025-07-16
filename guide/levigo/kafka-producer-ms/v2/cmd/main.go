package main

import (
	"sync"
	"time"

	"github.com/MateusLeviDev/internal/producer"
	"github.com/MateusLeviDev/internal/service"
	"github.com/MateusLeviDev/pkg/constants"
	"github.com/MateusLeviDev/pkg/logger"
	"github.com/MateusLeviDev/pkg/utils"
	"github.com/google/uuid"
)

// Design Pattern: FANOUT -<
// The advantage of the Fan-Out pattern is that tasks are executed in parallel by the workers
func main() {
	start := time.Now()

	correlationID := uuid.New().String()

	users, err := service.GetUSers()
	if err != nil {
		logger.ErrorAsync(err)
	}

	elapsed := time.Since(start)
	defer logger.InfoAsync("Reading 1,000,000 users from CSV took ", elapsed)

	kafkaProducerInstance, err := producer.NewProducer(constants.KafkaBootstrapServers, constants.KafkaTopic)
	if err != nil {
		logger.ErrorAsync("Failed to create kafkaProducerInstance:", err)
	}
	defer kafkaProducerInstance.Close()

	var wg sync.WaitGroup

	wg.Add(constants.NumWorkers)

	mainCh := make(chan func())

	channels := utils.Split(mainCh, constants.NumWorkers)

	// Start the workers
	for i := 0; i < constants.NumWorkers; i++ {
		go utils.Worker(channels[i], &wg)
	}

	// Send tasks to the main channel
	go func() {
		// Close the main channel when the function ends
		defer close(mainCh)

		// First task: Write users to JSON file
		mainCh <- func() {
			utils.WriteUsersToJSONFile(users, constants.JSONFileName)
		}

		// Second task: Send users to Kafka in batches
		mainCh <- func() {
			batches := utils.BatchUsers(users, constants.BatchSize)

			// Start time for sending batches to Kafka
			startBatchSend := time.Now()

			for _, batch := range batches {
				err := kafkaProducerInstance.ProduceBatch(batch, correlationID)
				if err != nil {
					logger.ErrorAsync("Error sending batch to Kafka:", err)
					return
				}
			}

			// Calculate elapsed time for sending batches to Kafka
			elapsedBatchSend := time.Since(startBatchSend)
			logger.InfoAsync("Sending batches to Kafka took ", elapsedBatchSend)
		}
	}()

	// Wait for all the workers to finish
	wg.Wait()
}
