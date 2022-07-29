package handler_test

import (
	"log"
	"os"
	"testing"

	"github.com/klever-io/klever-chains/common"
	"github.com/klever.io/getchain-raw-worker/handler"
	"github.com/klever.io/getchain-raw-worker/interfaces"
	"github.com/stretchr/testify/require"
)

var s interfaces.Handler

func TestMain(m *testing.M) {
	var err error
	s, err = handler.NewHandler(common.POLYGON, "MATIC_KAFKA_TOPIC")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestHandlerBlockBatch(t *testing.T) {
	err := s.Forward()
	require.NotNil(t, err)
}
