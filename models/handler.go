package models

import "github.com/klever.io/getchain-raw-worker/cliFlags"

type Options struct {
	Mode     cliFlags.Mode
	Forward  bool
	Backward bool
}
