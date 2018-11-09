package types

import (
	"github.com/c3systems/c3-go/core/eosclient"
	"github.com/c3systems/c3-go/core/ethereumclient"
)

// NewAddressResponse ...
type NewAddressResponse struct {
	Address string
}

// SendTxResponse ...
type SendTxResponse struct {
	TxHash *string
}

// GetInfoResponse ...
type GetInfoResponse struct {
	BlockHeight string
}

// Keys ...
type Keys struct {
	PEMFile  string
	Password string
}

// Config ...
type Config struct {
	URI             string
	Peer            string
	DataDir         string
	Keys            Keys
	BlockDifficulty int
	MempoolType     string
	RPCHost         string
	EOSClient       *eosclient.CheckpointClient
	EthereumClient  *ethereumclient.CheckpointClient
}
