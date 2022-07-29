package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	sdk "github.com/klever-io/getchain-sdk/interfaces"
)

type ConsumerWorker struct {
	Explorer                 sdk.Coin
	KafkaTopic               string
	PendingTransactionsTopic string
	KafkaConsumer            *kafka.Consumer
	KafkaProducer            *kafka.Producer
}
