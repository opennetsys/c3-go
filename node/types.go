package node

import (
	"context"
	"crypto/ecdsa"

	"github.com/c3systems/c3-go/core/p2p"
	"github.com/c3systems/c3-go/core/p2p/protobuff"
	nodestore "github.com/c3systems/c3-go/node/store"

	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
)

// Keys ...
// note: any concern keeping these in memory? Maybe only fetch when needed?
type Keys struct {
	Priv *ecdsa.PrivateKey
	Pub  *ecdsa.PublicKey
}

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
	Keys                Keys
	Protobyff           protobuff.Interface
	BlockDifficulty     int
}

// Service ...
type Service struct {
	props Props
}
