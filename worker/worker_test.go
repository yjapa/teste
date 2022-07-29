package worker_test

import (
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/klever-io/klever-chains/common"
	"github.com/klever.io/getchain-raw-worker/interfaces"
	"github.com/klever.io/getchain-raw-worker/worker"
	"github.com/stretchr/testify/require"
)

var s interfaces.Worker

func TestMain(m *testing.M) {
	var err error
	s, err = worker.NewWorker(common.POLYGON, "MATIC_KAFKA_TOPIC")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestSendPendingTransactionToKafka(t *testing.T) {
	tx, _, err := s.GetExplorer().GetTransaction("0x337a4d5878cc61010b6c58ab76a52aaa7ee6c01410158b9b5507dbf4d32fdbea")
	require.Nil(t, err)
	require.NotNil(t, tx)

	err = s.PendingTransactionToKafka(tx)
	require.Nil(t, err)
}

func TestSendBlockToKafka(t *testing.T) {
	blockNumber := big.NewInt(28441986)
	block, err := s.GetExplorer().GetBlock(blockNumber)
	require.Nil(t, err)
	require.NotNil(t, block)

	err = s.BlockToKafka(block)
	require.Nil(t, err)
}

func TestSendBlockBatchToKafka(t *testing.T) {
	from := big.NewInt(30615113)
	to := big.NewInt(30615114)

	blockBatch, err := s.GetExplorer().GetBlockBatch(from, to)
	require.Nil(t, err)
	require.NotNil(t, blockBatch)

	err = s.BlockBatchToKafka(blockBatch)
	require.Nil(t, err)
}
