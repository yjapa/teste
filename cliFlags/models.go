package cliFlags

type Direction string
type Mode string

const (
	Forward   Direction = "forward"
	Backward  Direction = "backward"
	Both      Direction = "both"
	RawWorker Mode      = "raw"
	Listener  Mode      = "listener"
	Parser    Mode      = "parser"
)
