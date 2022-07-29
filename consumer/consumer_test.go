package consumer_test

import (
	"log"
	"os"
	"testing"

	"github.com/klever-io/klever-chains/common"
	"github.com/klever.io/getchain-raw-worker/consumer"
	"github.com/klever.io/getchain-raw-worker/interfaces"
	"github.com/klever.io/getchain-raw-worker/parser/coin/eth"
	"github.com/stretchr/testify/assert"
)

var s interfaces.ConsumerWorker

func TestMain(m *testing.M) {
	var err error
	s, err = consumer.NewConsumerWorker(common.POLYGON, "MATIC_KAFKA_TOPIC")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestGetFromKafka(t *testing.T) {
	err := s.ConsumeFromKafka(nil)
	assert.NoError(t, err)
}

func TestParseFromKafka(t *testing.T) {
	parser := eth.NewETHParser("MATIC_PARSED")
	err := s.ConsumeFromKafka(parser.ParseBlock)
	assert.NoError(t, err)
}
