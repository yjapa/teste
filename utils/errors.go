package utils

import "errors"

var (
	ErrNotImplemented     = errors.New("not implemented yet")
	ErrCoinNotImplemented = errors.New("coin not implemented yet")
	ErrParseField         = errors.New("error parsing field")
)
