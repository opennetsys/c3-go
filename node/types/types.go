package types

import (
	"context"

	host "github.com/libp2p/go-libp2p-host"

	"github.com/c3systems/c3/core/chain"
	"github.com/c3systems/c3/node/pubsub"
	nodestore "github.com/c3systems/c3/node/store"
)

// Props ...
type Props struct {
	CTX        context.Context
	CH         chan interface{}
	Host       host.Host
	Store      nodestore.Interface // store is used to temporarily store blocks and txs for mining and verification
	Blockchain chain.Interface
	Pubsub     pubsub.Interface
	//Wallet     wallet.Interface
}

// Service ...
type Service struct {
	props Props
}

// NewAddressResponse ...
type NewAddressResponse struct {
	Address string
}

// SendTxResponse ...
type SendTxResponse struct {
	TxHash string
}

// GetInfoResponse ...
type GetInfoResponse struct {
	BlockHeight string
}

// CFG ...
type CFG struct {
	URI     string
	DataDir string
}
