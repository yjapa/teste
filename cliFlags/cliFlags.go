package cliFlags

import (
	"flag"
	"log"
	"os"

	"github.com/klever-io/klever-chains/common"
)

var (
	help = flag.Bool("help", false, "Show help")
)

func GetCliFlags() (Mode, common.CoinType, string, Direction) {
	mode := flag.String("mode", "listener", "if container must run as raw worker, listener or parser")
	coin := flag.String("coin", "", "coin symbol in klever-chains")
	kafkaTopic := flag.String("kafka-topic", "", "kafka topic to which worker will send informations")
	direction := flag.String("direction", "forward", "handler should work forward, backward or both")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if coin == nil || *coin == "" {
		log.Fatalln("coin parameter is required")
	}

	coinType, err := common.CoinTypeByName(*coin)
	if err != nil {
		log.Fatalln(err)
	}

	if kafkaTopic == nil || *kafkaTopic == "" {
		log.Fatalln("kafka topic is required")
	}

	if *mode != "raw" && *direction == "backward" {
		log.Fatalln("listener and parser must work forwardly")
	}
	return Mode(*mode), coinType, *kafkaTopic, Direction(*direction)
}
