package worker

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/klever-io/getchain-sdk/coin"
	sdk "github.com/klever-io/getchain-sdk/interfaces"
	"github.com/klever-io/getchain-sdk/models"
	"github.com/klever-io/klever-chains/common"
	"github.com/klever.io/getchain-raw-worker/interfaces"
	"github.com/klever.io/getchain-raw-worker/utils"
)

func NewWorker(chain common.CoinType, kafkaTopic string) (interfaces.Worker, error) {
	explorer, err := coin.GetCoin(int32(chain.Int()))
	if err != nil {
		return nil, err
	}

	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": utils.RequireEnv("KAFKA_BOOTSTRAP_SERVERS"),
			"message.max.bytes": 1024 * 1024 * 40,
		})
	if err != nil {
		return nil, err
	}

	return &Worker{
		Explorer:                 explorer,
		KafkaTopic:               kafkaTopic,
		PendingTransactionsTopic: fmt.Sprintf(`PENDING_TXS_%s`, kafkaTopic),
		KafkaProducer:            producer,
	}, nil
}

func (s *Worker) GetExplorer() sdk.Coin {
	return s.Explorer
}

func (s *Worker) GetKafkaTopic() string {
	return s.KafkaTopic
}

func (s *Worker) PendingTransactionToKafka(tx *models.Transaction) error {
	txByte, err := json.Marshal(tx)
	if err != nil {
		return nil
	}

	if err := s.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.PendingTransactionsTopic, Partition: kafka.PartitionAny},
		Value:          txByte,
	}, nil); err != nil {
		return err
	}

	return nil
}

func (s *Worker) BlockToKafka(block interface{}) error {
	blockByte, err := json.Marshal(block)
	if err != nil {
		return err
	}

	if err := s.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.KafkaTopic, Partition: kafka.PartitionAny},
		Value:          blockByte,
	}, nil); err != nil {
		return err
	}

	return nil
}

func (s *Worker) BlockBatchToKafka(blockBatch map[string]interface{}) error {
	var wg sync.WaitGroup
	for _, block := range blockBatch {
		wg.Add(1)
		block := block
		go func() {
			defer wg.Done()
			s.BlockToKafka(block)
		}()
	}

	wg.Wait()
	return nil
}
