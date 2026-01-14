package ports

import "fn-contract-settled/internal/models"

type ContractProducer interface {
	ProduceContractSettled(
		eventData models.ContractSettled,
		correlationID string,
	) error

	Close()
}