package interfaces

import (
	"math/big"

	"github.com/klever-io/getchain-sdk/interfaces"
	sdk "github.com/klever-io/getchain-sdk/models"
	"github.com/klever.io/getchain-raw-worker/indexer/data"
	"github.com/klever.io/getchain-raw-worker/models"
)

type Worker interface {
	GetExplorer() interfaces.Coin
	GetKafkaTopic() string
	PendingTransactionToKafka(tx *sdk.Transaction) error
	BlockToKafka(block interface{}) error
	BlockBatchToKafka(blockBatch map[string]interface{}) error
}

type Handler interface {
	Forward() bool
	Backward() bool
	Status() interface{}
	GetHighestBlock() *big.Int
	GetLowestBlock() *big.Int
	GetWorker() Worker
	GetConsumerWorker() ConsumerWorker
	Start(options models.Options) error
}

type ConsumerWorker interface {
	ConsumeFromKafka(parser func([]byte) (data.Block, error)) error
}

type Parser interface {
	ParseBlock(kafkaMsg []byte) (data.Block, error)
	ParseTransaction(tx interface{}) (data.Transaction, error)
	ParseInteraction(itx interface{}) (data.Interaction, error)
}
