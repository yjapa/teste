package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/klever.io/getchain-raw-worker/cliFlags"
	"github.com/klever.io/getchain-raw-worker/handler"
	"github.com/klever.io/getchain-raw-worker/models"
)

func main() {
	mode, coin, topic, direction := cliFlags.GetCliFlags()

	handler, err := handler.NewHandler(coin, topic)
	if err != nil {
		log.Fatalln(err)
	}

	var options models.Options

	options.Mode = mode

	switch direction {
	case cliFlags.Forward:
		options.Forward = true
	case cliFlags.Backward:
		options.Backward = true
	case cliFlags.Both:
		options.Forward = true
		options.Backward = true
	}

	err = handler.Start(options)
	if err != nil {
		log.Fatalln(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

}
