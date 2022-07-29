package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/klever-io/getchain-sdk/coin"
	"github.com/klever-io/klever-chains/common"
	"github.com/klever.io/getchain-raw-worker/indexer/data"
	"github.com/klever.io/getchain-raw-worker/interfaces"
	"github.com/klever.io/getchain-raw-worker/utils"
)

var (
	offsetsMap map[string]kafka.TopicPartition
	count      int
)

func NewConsumerWorker(chain common.CoinType, kafkaTopic string) (interfaces.ConsumerWorker, error) {
	explorer, err := coin.GetCoin(int32(chain.Int()))
	if err != nil {
		return nil, err
	}

	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": utils.RequireEnv("KAFKA_BOOTSTRAP_SERVERS"),
		})
	if err != nil {
		return nil, err
	}

	consumer, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":               utils.RequireEnv("KAFKA_BOOTSTRAP_SERVERS"),
			"broker.address.family":           "v4",
			"group.id":                        "klever-consumer",
			"enable.auto.commit":              false,
			"go.application.rebalance.enable": true,
			"auto.offset.reset":               "earliest",
		})
	if err != nil {
		return nil, err
	}

	return &ConsumerWorker{
		Explorer:                 explorer,
		KafkaTopic:               kafkaTopic,
		PendingTransactionsTopic: fmt.Sprintf(`PENDING_TXS_%s`, kafkaTopic),
		KafkaConsumer:            consumer,
		KafkaProducer:            producer,
	}, nil
}

func (s *ConsumerWorker) ConsumeFromKafka(parser func([]byte) (data.Block, error)) error {
	if err := s.KafkaConsumer.Subscribe(s.KafkaTopic, s.rebalanceCb); err != nil {
		return err
	}

	count = 0
	offsetsMap = make(map[string]kafka.TopicPartition)

	for {
		event := s.KafkaConsumer.Poll(100)
		if event == nil {
			continue
		}

		switch e := event.(type) {
		case *kafka.Message:
			key := fmt.Sprintf("%s[%d]", *e.TopicPartition.Topic, e.TopicPartition.Partition)
			offsetsMap[key] = e.TopicPartition

			if parser != nil {
				block, err := parser(event.(*kafka.Message).Value)
				if err != nil {
					return err
				}

				blockByte, err := json.Marshal(block)
				if err != nil {
					return err
				}

				topic := fmt.Sprintf("%s_PARSED", s.KafkaTopic)
				if err := s.KafkaProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          blockByte,
				}, nil); err != nil {
					return err
				}
			}
			count++
			if count%10 == 0 {
				// 10 messages per commit
				go s.commit(offsetsMap)
			}

		case kafka.Error:
			return e
		default:
			continue
		}
	}

}

func (s *ConsumerWorker) commit(offsets map[string]kafka.TopicPartition) error {
	if len(offsets) == 0 {
		return fmt.Errorf("offsets size zero")
	}

	topics := make([]kafka.TopicPartition, len(offsets))
	index := 0
	for _, topic := range offsets {
		topic.Offset = topic.Offset + 1
		topics[index] = topic
		index++
	}
	if _, err := s.KafkaConsumer.CommitOffsets(topics); err != nil {
		return err
	}

	return nil
}

func (s *ConsumerWorker) rebalanceCb(consumer *kafka.Consumer, event kafka.Event) error {
	switch e := event.(type) {
	case kafka.AssignedPartitions:
		count = 0
		offsetsMap = make(map[string]kafka.TopicPartition)
		s.KafkaConsumer.Assign(e.Partitions)
	case kafka.RevokedPartitions:
		s.commit(offsetsMap)
		s.KafkaConsumer.Unassign()
	}
	return nil
}
