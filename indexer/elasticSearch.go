package indexer

import (
	"bytes"
	"github.com/klever.io/getchain-raw-worker/indexer/templates"
	"log"
)

const (
	BlockIndex = "blocks"
)

//TODO: create instance of elasticClient

//TODO: create handler for elastic server

func GetTemplates() map[string]*bytes.Buffer {
	temps := make(map[string]*bytes.Buffer)

	blockBuffer, err := templates.Block.ToBuffer()
	if err != nil {
		log.Fatal(err)
	}

	temps[BlockIndex] = blockBuffer

	return temps
}
