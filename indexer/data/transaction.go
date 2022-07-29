package data

import "time"

type Block struct {
	Hash         string        `form:"hash" json:"hash" xml:"hash"`
	BlockHeight  string        `form:"blockHeight" json:"blockHeight" xml:"blockHeight"`
	BlockRewards string        `form:"blockRewards" json:"blockRewards" xml:"blockRewards"`
	Size         string        `form:"size" json:"size" xml:"size"`
	Timestamp    time.Time     `form:"timestamp" json:"timestamp" xml:"timestamp"`
	Transactions []Transaction `form:"transactions" json:"transactions" xml:"transactions"`
}

// Transaction generic data type for elastic
type Transaction struct {
	ID           int64         `form:"id" json:"id" xml:"id"`
	TxID         string        `form:"txId" json:"txId" xml:"txId"`
	Block        string        `form:"block" json:"block" xml:"block"`
	FromAddress  string        `form:"fromAddress" json:"fromAddress" xml:"fromAddress"`
	ToAddress    string        `form:"toAddress" json:"toAddress" xml:"toAddress"`
	Timestamp    time.Time     `form:"timestamp" json:"timestamp" xml:"timestamp"`
	Value        string        `form:"value" json:"value" xml:"value"`
	Fee          string        `form:"fee" json:"fee" xml:"fee"`
	Status       int32         `form:"status" json:"status" xml:"status"`
	Specific     interface{}   `form:"specific" json:"specific" xml:"specific"`
	Interactions []Interaction `form:"interactions" json:"interactions"  xml:"interactions"`
}

type Interaction struct {
	TxID         string `form:"txId" json:"txId" xml:"txId"`
	FromAddress  string `form:"fromAddress" json:"fromAddress" xml:"fromAddress"`
	ToAddress    string `form:"toAddress" json:"toAddress" xml:"toAddress"`
	Value        string `form:"value" json:"value" xml:"value"`
	ContractType string `form:"contractType" json:"contractType" xml:"contractType"`
	ContractName string `form:"contractName" json:"contractName" xml:"contractName"`
}
