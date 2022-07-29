package parser

import (
	"fmt"

	"github.com/klever-io/klever-chains/common"
	"github.com/klever.io/getchain-raw-worker/indexer/data"
	"github.com/klever.io/getchain-raw-worker/interfaces"
	"github.com/klever.io/getchain-raw-worker/utils"
)

func NewParser(chain common.CoinType, kafkaTopic string) interfaces.Parser {
	return &Parser{
		KafkaTopic:               kafkaTopic,
		PendingTransactionsTopic: fmt.Sprintf(`PENDING_TXS_%s`, kafkaTopic),
	}
}

func (s *Parser) ParseBlock(kafkaMsg []byte) (data.Block, error) {
	return data.Block{}, utils.ErrNotImplemented
}
func (s *Parser) ParseTransaction(tx interface{}) (data.Transaction, error) {
	return data.Transaction{}, utils.ErrNotImplemented
}
func (s *Parser) ParseInteraction(itx interface{}) (data.Interaction, error) {
	return data.Interaction{}, utils.ErrNotImplemented
}
