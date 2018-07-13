package node

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/mainchain/miner"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/p2p"
	"github.com/c3systems/c3/core/p2p/protobuff"
	pb "github.com/c3systems/c3/core/p2p/protobuff/pb"
	"github.com/c3systems/c3/core/p2p/store/leveldbstore"
	"github.com/c3systems/c3/node/store/safemempool"
	nodetypes "github.com/c3systems/c3/node/types"

	ipfsaddr "github.com/ipfs/go-ipfs-addr"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	csms "github.com/libp2p/go-conn-security-multistream"
	floodsub "github.com/libp2p/go-floodsub"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	secio "github.com/libp2p/go-libp2p-secio"
	swarm "github.com/libp2p/go-libp2p-swarm"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	tcp "github.com/libp2p/go-tcp-transport"
	ma "github.com/multiformats/go-multiaddr"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

// Start ...
// note: start is called from cobra
func Start(cfg *nodetypes.Config) error {
	n := new(Service)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg == nil {
		// note: is this the correct way to fail an app with cobra?
		return errors.New("config is required to start the node")
	}

	var pwd *string
	if cfg.Keys.Password != "" {
		pwd = &cfg.Keys.Password
	}

	priv, err := c3crypto.ReadPrivateKeyFromPem(cfg.Keys.PEMFile, pwd)
	if err != nil {
		return fmt.Errorf("err reading pem file\n%v", err)
	}
	pub := &priv.PublicKey

	wPriv, wPub := c3crypto.NewWrappedKeyPair(priv)

	pid, err := peer.IDFromPublicKey(wPub)
	if err != nil {
		return fmt.Errorf("err generating pid from public key\n%v", err)
	}

	listen, err := ma.NewMultiaddr(cfg.URI)
	if err != nil {
		return fmt.Errorf("err parsing ipfs uri\n%v", err)
	}

	ps := peerstore.NewPeerstore()
	ps.AddPrivKey(pid, wPriv)
	ps.AddPubKey(pid, wPub)

	swarmNet := swarm.NewSwarm(ctx, pid, ps, nil)
	tcpTransport := tcp.NewTCPTransport(genUpgrader(swarmNet))
	if err := swarmNet.AddTransport(tcpTransport); err != nil {
		return fmt.Errorf("err adding tcp transport\n%v", err)
	}
	if err := swarmNet.AddListenAddr(listen); err != nil {
		return fmt.Errorf("err adding swam listen addr\n%v", err)
	}
	newNode := bhost.New(swarmNet)

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

		newNode.Peerstore().AddAddrs(pinfo.ID, pinfo.Addrs, peerstore.PermanentAddrTTL)
	}

	// TODO: add cli flags for different types
	memPool, err := safemempool.New(&safemempool.Props{})
	if err != nil {
		return fmt.Errorf("err initializing mempool\n%v", err)
	}

	// TODO: add cli flags for different types
	// diskStore, err := fsstore.New(cfg.DataDir)
	diskStore, err := leveldbstore.New(cfg.DataDir, nil)
	// diskStore, err := leveldbds.NewDatastore(cfg.DataDir, nil)
	if err != nil {
		return fmt.Errorf("err building disk store\n%v", err)
	}
	// wrap the datastore in a 'content addressed blocks' layer
	// TODO: implement metrics? https://github.com/ipfs/go-ds-measure
	blocks := bstore.NewBlockstore(diskStore)

	p2pSvc, err := p2p.New(&p2p.Props{
		BlockStore: blocks,
		Host:       newNode,
	})
	if err != nil {
		return fmt.Errorf("err starting ipfs p2p network\n%v", err)
	}

	pBuff, err := protobuff.NewNode(&protobuff.Props{
		Host:                   newNode,
		GetHeadBlockFN:         memPool.GetHeadBlock,
		BroadcastTransactionFN: n.BroadcastTransaction,
	})
	if err != nil {
		return fmt.Errorf("err starting protobuff node\n%v", err)
	}

	nextBlock := &mainchain.GenesisBlock
	peers := newNode.Peerstore().Peers()
	if len(peers) > 1 {
		if err := fetchHeadBlock(nextBlock, peers, pBuff); err != nil {
			return fmt.Errorf("err fetching headblock\n%v", err)
		}
	}

	if err := memPool.SetHeadBlock(nextBlock); err != nil {
		return fmt.Errorf("err setting head block\n%v", err)
	}

	n, err = New(&Props{
		Context:             ctx,
		SubscriberChannel:   make(chan interface{}),
		CancelMinersChannel: make(chan struct{}),
		Host:                newNode,
		Store:               memPool,
		Pubsub:              pubsub,
		P2P:                 p2pSvc,
		Protobyff:           pBuff,
		Keys: Keys{
			Priv: priv,
			Pub:  pub,
		},
	})
	if err != nil {
		return fmt.Errorf("err building the node\n%v", err)
	}

	if err := n.listenForEvents(); err != nil {
		return fmt.Errorf("err starting listener\n%v", err)
	}
	// TODO: add a cli flag to determine if the node mines
	if err := n.spawnNextBlockMiner(nextBlock); err != nil {
		return fmt.Errorf("err starting miner\n%v", err)
	}
	log.Printf("Node %s started", newNode.ID().Pretty())

	for {
		switch v := <-n.props.SubscriberChannel; v.(type) {
		case error:
			log.Println("[node] received an error on the channel", err)

		case *miner.MinedBlock:
			log.Print("[node] received mined block")
			b, _ := v.(*miner.MinedBlock)
			go n.handleReceiptOfMainchainBlock(b)

		case *statechain.Transaction:
			log.Print("[node] received statechain transaction")
			tx, _ := v.(*statechain.Transaction)
			go n.handleReceiptOfStatechainTransaction(tx)

		default:
			log.Printf("[node] received an unknown message on channel of type %T\n%v", v, v)
		}
	}
}

// note: https://github.com/libp2p/go-libp2p-swarm/blob/da01184afe4c67bec58c5e73f3350ad80b624c0d/testing/testing.go#L39
func genUpgrader(n *swarm.Swarm) *tptu.Upgrader {
	id := n.LocalPeer()
	pk := n.Peerstore().PrivKey(id)
	secMuxer := new(csms.SSMuxer)
	secMuxer.AddTransport(secio.ID, &secio.Transport{
		LocalID:    id,
		PrivateKey: pk,
	})

	stMuxer := msmux.NewBlankTransport()
	stMuxer.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)

	return &tptu.Upgrader{
		Secure:  secMuxer,
		Muxer:   stMuxer,
		Filters: n.Filters,
	}

}

func fetchHeadBlock(headBlock *mainchain.Block, peers []peer.ID, pBuff protobuff.Interface) error {
	// TODO: pass contexts to pBuff functions
	ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ch := make(chan interface{})

	// note: err == nil, here!
	if err := pBuff.FetchHeadBlock(peers[1], ch); err != nil {
		return err
	}
	select {
	case v := <-ch:
		switch v.(type) {
		case error:
			err, _ := v.(error)
			return err

		case *pb.HeadBlockResponse:
			hb, _ := v.(*pb.HeadBlockResponse)
			if hb == nil {
				return errors.New("received nil headblock")
			}

			block := new(mainchain.Block)
			if err := block.Deserialize(hb.HeadBlockBytes); err != nil {
				return err
			}

			headBlock = block
			return nil

		default:
			return fmt.Errorf("received unknown type %T\n%v", v, v)

		}

	case <-ctx1.Done():
		return errors.New("fetching headblock from peer timedout")

	}
}
