package handler

import (
	sdk "github.com/klever-io/getchain-sdk/interfaces"
	"github.com/klever-io/klever-chains/common"
	"github.com/klever.io/getchain-raw-worker/interfaces"
	"github.com/klever.io/getchain-raw-worker/parser/coin/eth"
)

type Handler struct {
	Worker         interfaces.Worker
	ConsumerWorker interfaces.ConsumerWorker
	Redis          sdk.Storage
	Coin           common.CoinType
	KafkaTopic     string
}

var parsers = map[common.CoinType]func(kafkaTopic string) interfaces.Parser{
	common.ETH: eth.NewETHParser,
}
