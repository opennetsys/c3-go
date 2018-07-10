package node

import (
	"context"

	"github.com/c3systems/c3/core/chain/mainchain/miner"
	nodestore "github.com/c3systems/c3/node/store"
	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
)

// Props ...
type Props struct {
	// Note: it's preferable to not use abbreviations if it means uppercasing all the letters. Acronyms are OK.
	Context context.Context
	Channel chan interface{}
	Host    host.Host
	Store   nodestore.Interface // store is used to temporarily store blocks and txs for mining and verification
	Miner   miner.Interface
	Pubsub  *floodsub.PubSub // note: how to make this into an interface?
	//Wallet     wallet.Interface
}

// Service ...
type Service struct {
	props Props
}
