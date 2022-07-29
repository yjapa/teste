package eth

import (
	"encoding/json"
	"math/big"
	"time"

	sdk "github.com/klever-io/getchain-sdk/models"
	"github.com/klever-io/klever-chains/common"
	"github.com/klever.io/getchain-raw-worker/indexer/data"
	"github.com/klever.io/getchain-raw-worker/interfaces"
	"github.com/klever.io/getchain-raw-worker/parser"
	"github.com/klever.io/getchain-raw-worker/utils"
)

type ParserETH struct {
	interfaces.Parser
}

func NewETHParser(kafkaTopic string) interfaces.Parser {
	parser := parser.NewParser(common.ETH, kafkaTopic)

	return &ParserETH{
		Parser: parser,
	}
}

func (s *ParserETH) ParseBlock(kafkaMsg []byte) (data.Block, error) {
	var rawBlock *sdk.RawBlock
	if err := json.Unmarshal(kafkaMsg, &rawBlock); err != nil {
		return data.Block{}, err
	}

	blockHeight, ok := new(big.Int).SetString(rawBlock.Block.Number, 0)
	if !ok {
		return data.Block{}, utils.ErrParseField
	}

	size, ok := new(big.Int).SetString(rawBlock.Block.Size, 0)
	if !ok {
		return data.Block{}, utils.ErrParseField
	}

	timestamp, ok := new(big.Int).SetString(rawBlock.Block.Timestamp, 0)
	if !ok {
		return data.Block{}, utils.ErrParseField
	}

	unixTimestamp := time.Unix(timestamp.Int64(), 0)

	txsMaps := make(map[string]data.Transaction, 0)
	for _, tx := range rawBlock.Block.Transactions {
		parsedTransaction, err := s.ParseTransaction(tx)
		if err != nil {
			return data.Block{}, err
		}
		txsMaps[tx.Hash] = parsedTransaction
	}

	mapLogs, err := s.getInteractionsMap(rawBlock.Logs)
	if err != nil {
		return data.Block{}, err
	}

	var transactions []data.Transaction
	for txHash, tx := range txsMaps {
		tx.Timestamp = unixTimestamp
		tx.Interactions = append(tx.Interactions, mapLogs[txHash]...)
		transactions = append(transactions, tx)
	}

	block := data.Block{
		Hash:         rawBlock.Block.Hash,
		BlockHeight:  blockHeight.String(),
		Timestamp:    unixTimestamp,
		Size:         size.String(),
		Transactions: transactions,
	}

	return block, nil
}

func (s *ParserETH) ParseTransaction(tx interface{}) (data.Transaction, error) {

	ethTransaction := tx.(sdk.ETHTransaction)

	block, ok := new(big.Int).SetString(ethTransaction.BlockNumber, 0)
	if !ok {
		return data.Transaction{}, utils.ErrParseField
	}

	value, ok := new(big.Int).SetString(ethTransaction.Value, 0)
	if !ok {
		return data.Transaction{}, utils.ErrParseField
	}

	gasUsed, ok := new(big.Int).SetString(ethTransaction.Gas, 0)
	if !ok {
		return data.Transaction{}, utils.ErrParseField
	}

	gasPrice, ok := new(big.Int).SetString(ethTransaction.GasPrice, 0)
	if !ok {
		return data.Transaction{}, utils.ErrParseField
	}

	fee := new(big.Int).Mul(gasUsed, gasPrice)

	return data.Transaction{
		TxID:        ethTransaction.Hash,
		Block:       block.String(),
		FromAddress: ethTransaction.From,
		ToAddress:   ethTransaction.To,
		Value:       value.String(),
		Fee:         fee.String(),
		Status:      1,
	}, nil
}

func (s *ParserETH) ParseInteraction(itx interface{}) (data.Interaction, error) {
	ethInteraction := itx.(sdk.ETHLogs)

	from := "0x" + ethInteraction.Topics[1][26:]
	to := "0x" + ethInteraction.Topics[2][26:]
	value, ok := new(big.Int).SetString(ethInteraction.Data, 0)
	if !ok {
		return data.Interaction{}, utils.ErrParseField
	}

	parsedLog := data.Interaction{
		TxID:         ethInteraction.TransactionHash,
		FromAddress:  from,
		ToAddress:    to,
		Value:        value.String(),
		ContractName: ethInteraction.Address,
	}

	return parsedLog, nil
}

func (s *ParserETH) getInteractionsMap(ethLogs []sdk.ETHLogs) (map[string][]data.Interaction, error) {
	mapLogs := make(map[string][]data.Interaction, 0)
	for _, log := range ethLogs {
		if (len(log.Topics) != 3) || (log.Topics[0] != ETH_TOKEN_TRANSFER) {
			continue
		}

		parsedLog, err := s.ParseInteraction(log)
		if err != nil {
			return nil, err
		}

		mapLogs[log.TransactionHash] = append(mapLogs[log.TransactionHash], parsedLog)
	}

	return mapLogs, nil
}
