package node

import (
	"context"

	"github.com/c3systems/c3/core/p2p"
	nodestore "github.com/c3systems/c3/node/store"
	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
)

// Props ...
type Props struct {
	// Note: it's preferable to not use abbreviations if it means uppercasing all the letters. Acronyms are OK.
	Context             context.Context
	SubscriberChannel   chan interface{}
	CancelMinersChannel chan struct{}
	Host                host.Host
	Store               nodestore.Interface // store is used to temporarily store blocks and txs for mining and verification
	Pubsub              *floodsub.PubSub    // note: how to make this into an interface?
	P2P                 p2p.Interface
	//Wallet     wallet.Interface
}

// Service ...
type Service struct {
	props Props
}
