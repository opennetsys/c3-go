package coder

import "errors"

const (
	// Proto3 is the protobuf syntax=3; https://developers.google.com/protocol-buffers/docs/proto3
	Proto3 byte = iota
)

// CurrentCode is the code currently used by this package
const CurrentCode = Proto3

var (
	// ErrNilBytes ...
	ErrNilBytes = errors.New("nil bytes")
)
