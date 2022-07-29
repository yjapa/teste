package worker

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/klever-io/getchain-sdk/interfaces"
)

type Worker struct {
	Explorer                 interfaces.Coin
	KafkaTopic               string
	PendingTransactionsTopic string
	KafkaProducer            *kafka.Producer
}
