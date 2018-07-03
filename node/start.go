package node

import (
	"log"

	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/p2p"
	"github.com/c3systems/c3/core/p2p/store/fsstore"
	floodsub "github.com/c3systems/c3/node/pubsub/floodsub"
	"github.com/c3systems/c3/node/store/safemempool"
	nodetypes "github.com/c3systems/c3/node/types"
	//"github.com/c3systems/c3/core/node/wallet"

	ipfsaddr "github.com/ipfs/go-ipfs-addr"
	libp2p "github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

// Start ...
// note: start is called from cobra
func Start(cfg *nodetypes.CFG) {
	if cfg == nil {
		// note: is this the correct way to fail an app with cobra?
		log.Fatal("config is required to start the node")
	}

	newNode, err := libp2p.New(props.Ctx, libp2p.Defaults)
	if err != nil {
		log.Fatal("err building libp2p service", err)
	}

	pubsub, err := floodsub.New(ctx, newNode)
	if err != nil {
		log.Fatal("err building new pubsub service", err)
	}

	for i, addr := range newNode.Addrs() {
		log.Printf("%d: %s/ipfs/%s\n", i, addr, newNode.ID().Pretty())
	}

	addr, err := ipfsaddr.ParseString(cfg.URI)
	if err != nil {
		log.Fatalf("err parsing node uri flag: %s\n%v", cfg.URI, err)
	}
	log.Println("Node Address:", addr)

	pinfo, err := peerstore.InfoFromP2pAddr(addr.Multiaddr())
	if err != nil {
		log.Fatal("err getting info from peerstore", err)
	}

	if err := newNode.Connect(ctx, *pinfo); err != nil {
		log.Fatal("bootstrapping a peer failed", err)
	}

	memPool := safemempool.New(&safemempool.Props{})
	diskStore, err := fsstore.New(cfg.DataDir, nil, true)
	if err != nil {
		log.Fatal("err building disk store", err)
	}

	p2p, err := p2p.New(&p2p.Props{
		BlockStore: diskStore,
		Host:       newNode,
	})
	if err != nil {
		log.Fatal("err starting ipfs p2p network", err)
	}

	blockchain, err := chain.New(&chain.Props{
		P2P: p2p,
	})
	if err != nil {
		log.Fatal("err building the blockchain", err)
	}

	c := ctx.Background()
	ch := make(chan interface{})
	n, err := node.New(&node.Props{
		CTX:        c,
		CH:         ch,
		Host:       newNode,
		Store:      memPool,
		Blockchain: blockchain,
		PubSub:     pubsub,
	})
	if err != nil {
		log.Fatal("err building the node", err)
	}

	if err := n.Start(); err != nil {
		log.Fatal("err starting node", err)
	}

	for {
		switch v := <-ch; v.(type) {
		case error:
			log.Println("[node] received an error on the channel", err)

		case *mainchain.Block:
			fallthrough
		case *statechain.Block:
			fallthrough
		case *statechain.Transaction:
			// do a stoofs

		default:
			log.Printf("[node] received an unknown message on channel of type %T\n%v", v, v)
		}
	}
	//blockchain := NewBlockchain(newNode)

	//node.p2pNode = newNode
	//node.mempool = NewMempool()
	//node.pubsub = pubsub
	//node.blockchain = blockchain
	//node.wallet = NewWallet()
}
