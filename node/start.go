package node

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/c3systems/c3/core/chain"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/p2p"
	"github.com/c3systems/c3/core/p2p/store/fsstore"
	"github.com/c3systems/c3/node/store/safemempool"
	nodetypes "github.com/c3systems/c3/node/types"
	//"github.com/c3systems/c3/node/wallet"

	ipfsaddr "github.com/ipfs/go-ipfs-addr"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	floodsub "github.com/libp2p/go-floodsub"
	libp2p "github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

// Start ...
// note: start is called from cobra
func Start(cfg *nodetypes.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//c := context.Background()

	if cfg == nil {
		// note: is this the correct way to fail an app with cobra?
		return errors.New("config is required to start the node")
	}

	newNode, err := libp2p.New(ctx, libp2p.Defaults, libp2p.ListenAddrStrings(cfg.URI))
	if err != nil {
		return fmt.Errorf("err building libp2p service\n%v", err)
	}

	pubsub, err := floodsub.NewFloodSub(ctx, newNode)
	if err != nil {
		return fmt.Errorf("err building new pubsub service\n%v", err)
	}

	for i, addr := range newNode.Addrs() {
		log.Printf("%d: %s/ipfs/%s\n", i, addr, newNode.ID().Pretty())
	}

	if cfg.Peer != "" {
		addr, err := ipfsaddr.ParseString(cfg.Peer)
		if err != nil {
			return fmt.Errorf("err parsing node uri flag: %s\n%v", cfg.URI, err)
		}

		pinfo, err := peerstore.InfoFromP2pAddr(addr.Multiaddr())
		if err != nil {
			return fmt.Errorf("err getting info from peerstore\n%v", err)
		}

		if err := newNode.Connect(ctx, *pinfo); err != nil {
			log.Printf("bootstrapping a peer failed\n%v", err)
		}
	}

	// TODO: add cli flags for different types
	memPool, err := safemempool.New(&safemempool.Props{})
	if err != nil {
		return fmt.Errorf("err initializing mempool\n%v", err)
	}
	diskStore, err := fsstore.New(cfg.DataDir)
	if err != nil {
		return fmt.Errorf("err building disk store\n%v", err)
	}
	// wrap the datastore in a 'content addressed blocks' layer
	blocks := bstore.NewBlockstore(diskStore)

	p2p, err := p2p.New(&p2p.Props{
		BlockStore: blocks,
		Host:       newNode,
	})
	if err != nil {
		return fmt.Errorf("err starting ipfs p2p network\n%v", err)
	}

	blockchain, err := chain.New(&chain.Props{
		P2P: p2p,
	})
	if err != nil {
		return fmt.Errorf("err building the blockchain\n%v", err)
	}

	ch := make(chan interface{})
	n, err := New(&Props{
		Context:    ctx,
		Channel:    ch,
		Host:       newNode,
		Store:      memPool,
		Blockchain: blockchain,
		Pubsub:     pubsub,
	})
	if err != nil {
		return fmt.Errorf("err building the node\n%v", err)
	}

	if err := n.Start(); err != nil {
		return fmt.Errorf("err starting node\n%v", err)
	}

	//go func() {
	log.Printf("Node %s started", newNode.ID().Pretty())
	hash := "fakeHash"
	tx := statechain.NewTransaction(&statechain.TransactionProps{
		TxHash:  &hash,
		Method:  "foo",
		Payload: "bar",
	})
	res, err := n.BroadcastTransaction(tx)
	if err != nil {
		log.Printf("err broadcasting tx\n%v", err)
	}
	log.Printf("tx resp\n%v", res)

	for {
		switch v := <-ch; v.(type) {
		case error:
			log.Println("[node] received an error on the channel", err)

		case *mainchain.Block, *statechain.Block, *statechain.Transaction:
			// do a stoofs
			log.Printf("received %T\n%v", v, v)

		default:
			log.Printf("[node] received an unknown message on channel of type %T\n%v", v, v)
		}
	}
	//}()

	//return nil
	//blockchain := NewBlockchain(newNode)

	//node.p2pNode = newNode
	//node.mempool = NewMempool()
	//node.pubsub = pubsub
	//node.blockchain = blockchain
	//node.wallet = NewWallet()
}
