package node

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/config"
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/statechain"
	"github.com/c3systems/c3-go/core/eosclient"
	"github.com/c3systems/c3-go/core/ethereumclient"
	"github.com/c3systems/c3-go/core/miner"
	"github.com/c3systems/c3-go/core/p2p"
	"github.com/c3systems/c3-go/core/p2p/protobuff"
	"github.com/c3systems/c3-go/core/p2p/store/leveldbstore"
	"github.com/c3systems/c3-go/core/sandbox"
	colorlog "github.com/c3systems/c3-go/log/color"
	loghooks "github.com/c3systems/c3-go/log/hooks"
	nodestore "github.com/c3systems/c3-go/node/store"
	"github.com/c3systems/c3-go/node/store/redisstore"
	"github.com/c3systems/c3-go/node/store/safemempool"
	nodetypes "github.com/c3systems/c3-go/node/types"
	"github.com/c3systems/c3-go/rpc"
	redis "github.com/gomodule/redigo/redis"
	ipfsaddr "github.com/ipfs/go-ipfs-addr"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	lCrypt "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pstoremem "github.com/libp2p/go-libp2p-peerstore/pstoremem"
	floodsub "github.com/libp2p/go-libp2p-pubsub"
	swarm "github.com/libp2p/go-libp2p-swarm"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	tcp "github.com/libp2p/go-tcp-transport"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
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
	EOSClient           *eosclient.CheckpointClient
	EthereumClient      *ethereumclient.CheckpointClient
}

// Service ...
type Service struct {
	props Props
}

// newNode ...
func newNode(props *Props) (*Service, error) {
	if props == nil {
		return nil, errors.New("props cannot be nil")
	}
	if props.Store == nil || props.Pubsub == nil || props.P2P == nil {
		return nil, errors.New("p2p node, store, and pubsub are required")
	}

	return &Service{
		props: *props,
	}, nil
}

// NewFullNode ...
func NewFullNode(cfg *nodetypes.Config) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("config is required to start the node")
	}

	var pwd *string
	if cfg.Keys.Password != "" {
		pwd = &cfg.Keys.Password
	}

	priv, err := c3crypto.ReadPrivateKeyFromPem(cfg.Keys.PEMFile, pwd)
	if err != nil {
		return nil, fmt.Errorf("err reading pem file\n%v", err)
	}
	pub := &priv.PublicKey

	wPriv, wPub, err := lCrypt.ECDSAKeyPairFromKey(priv)
	if err != nil {
		return nil, fmt.Errorf("err generating key pairs\n%v", err)
	}

	pid, err := peer.IDFromPublicKey(wPub)
	if err != nil {
		return nil, fmt.Errorf("err generating pid from public key\n%v", err)
	}

	listen, err := ma.NewMultiaddr(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("err parsing ipfs uri\n%v", err)
	}

	ps := pstoremem.NewPeerstore()
	if err := ps.AddPrivKey(pid, wPriv); err != nil {
		return nil, fmt.Errorf("err adding priv key\n%v", err)
	}
	if err := ps.AddPubKey(pid, wPub); err != nil {
		return nil, fmt.Errorf("err adding pub key\n%v", err)
	}

	ctx := context.Background()
	swarmNet := swarm.NewSwarm(ctx, pid, ps, nil)
	tcpTransport := tcp.NewTCPTransport(genUpgrader(swarmNet))
	if err := swarmNet.AddTransport(tcpTransport); err != nil {
		return nil, fmt.Errorf("err adding tcp transport\n%v", err)
	}
	if err := swarmNet.AddListenAddr(listen); err != nil {
		return nil, fmt.Errorf("err adding swam listen addr\n%v", err)
	}
	bNode := bhost.New(swarmNet)

	dhtSvc, err := dht.New(ctx, bNode)
	if err != nil {
		return nil, fmt.Errorf("err building dht svc\n%v", err)
	}
	if err := dhtSvc.Bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("err bootstrapping dht\n%v", err)
	}

	newNode := rhost.Wrap(bNode, dhtSvc)
	h = newNode

	discoverySvc, err := discovery.NewMdnsService(ctx, newNode, time.Second, "c3")
	if err != nil {
		return nil, fmt.Errorf("error starting discovery service\n%v", err)
	}
	discoverySvc.RegisterNotifee(&DiscoveryNotifee{newNode})

	pubsub, err := floodsub.NewFloodSub(ctx, newNode)
	if err != nil {
		return nil, fmt.Errorf("err building new pubsub service\n%v", err)
	}

	for i, addr := range newNode.Addrs() {
		log.Printf(colorlog.Green("[node] %d: %s/ipfs/%s\n", i, addr, newNode.ID().Pretty()))
	}

	if cfg.Peer != "" {
		addr, err := ipfsaddr.ParseString(cfg.Peer)
		if err != nil {
			return nil, fmt.Errorf("err parsing node uri flag: %s\n%v", cfg.URI, err)
		}

		pinfo, err := peerstore.InfoFromP2pAddr(addr.Multiaddr())
		if err != nil {
			return nil, fmt.Errorf("err getting info from peerstore\n%v", err)
		}

		log.Println("[node] FULL", addr.String())
		log.Println("[node] PIN INFO", pinfo)

		if err := newNode.Connect(ctx, *pinfo); err != nil {
			return nil, fmt.Errorf("[node] bootstrapping a peer failed\n%v", err)
		}

		newNode.Peerstore().AddAddrs(pinfo.ID, pinfo.Addrs, peerstore.PermanentAddrTTL)
	}

	var memPool nodestore.Interface

	mempoolType := strings.ToLower(cfg.MempoolType)
	switch mempoolType {
	case "redis":
		log.Println(`[node] mempool type is "redis"`)
		// TODO: move to config
		redisaddr := "localhost:6379"
		redispool := &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisaddr) },
		}
		memPool, err = redisstore.New(&redisstore.Props{
			Pool: redispool,
		})
		if err != nil {
			return nil, fmt.Errorf("[node] err initializing redisstore\n%v", err)
		}
	case "memory":
		fallthrough
	default:
		log.Println(`[node] mempool type is "memory"`)
		memPool, err = safemempool.New(&safemempool.Props{})
		if err != nil {
			return nil, fmt.Errorf("[node] err initializing mempool\n%v", err)
		}
	}

	// TODO: add cli flags for different types
	diskStore, err := leveldbstore.New(cfg.DataDir, nil)
	if err != nil {
		return nil, fmt.Errorf("[node] err building disk store\n%v", err)
	}
	// wrap the datastore in a 'content addressed blocks' layer
	// TODO: implement metrics? https://github.com/ipfs/go-ds-measure
	blocks := bstore.NewBlockstore(diskStore)

	p2pSvc, err := p2p.New(&p2p.Props{
		BlockStore: blocks,
		Host:       newNode,
		Router:     dhtSvc,
	})
	if err != nil {
		return nil, fmt.Errorf("error starting ipfs p2p network\n%v", err)
	}

	n := new(Service)
	pBuff, err := protobuff.NewNode(&protobuff.Props{
		Host:                   newNode,
		GetHeadBlockFN:         memPool.GetHeadBlock,
		BroadcastTransactionFN: n.BroadcastTransaction,
		AddPendingTxFN:         memPool.AddTx,
	})
	if err != nil {
		return nil, fmt.Errorf("error starting protobuff node\n%v", err)
	}

	initialBlock := &mainchain.GenesisBlock

	// set head block to last mainchain block that was stored
	cachedLatestBlock := make(chan *mainchain.Block, 1)
	go func() {
		latestBlock, err := p2pSvc.GetLatestBlock()
		if err == nil {
			cachedLatestBlock <- latestBlock
		}
	}()

	select {
	case latestBlock := <-cachedLatestBlock:
		initialBlock = latestBlock
	case <-time.After(2 * time.Second):
	}

	c, err := p2pSvc.SetMainchainBlock(initialBlock)
	if err != nil {
		log.Errorf("[miner] error setting mainchain genesis block; error %s", err)
		return nil, err
	}

	log.Printf("[miner] set mainchain genesis block with cid %v", c)

	nextBlock := initialBlock

	peers := newNode.Peerstore().Peers()
	if len(peers) > 1 {
		if err := sendEcho(newNode.ID(), peers, pBuff); err != nil {
			log.Errorln("error echoing peer; is peer online?")
			return nil, fmt.Errorf("err echoing peer\n%v", err)
		}
		if err := fetchHeadBlock(newNode.ID(), nextBlock, peers, pBuff); err != nil {
			return nil, fmt.Errorf("err fetching headblock\n%v", err)
		}
	}

	if err := memPool.SetHeadBlock(nextBlock); err != nil {
		return nil, fmt.Errorf("err setting head block\n%v", err)
	}

	nb := &net.NotifyBundle{
		ConnectedF: onConn,
	}
	newNode.Network().Notify(nb)

	n.props = Props{
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
		BlockDifficulty: cfg.BlockDifficulty,
		EOSClient:       cfg.EOSClient,
		EthereumClient:  cfg.EthereumClient,
	}

	if err := n.listenForEvents(); err != nil {
		return nil, fmt.Errorf("error starting listener\n%v", err)
	}
	// TODO: add a cli flag to determine if the node mines
	if err := n.spawnNextBlockMiner(nextBlock); err != nil {
		return nil, fmt.Errorf("error starting miner in main start method\n%v", err)
	}
	log.Printf("[node] started %s", newNode.ID().Pretty())

	if cfg.RPCHost != "" {
		// start rpc service
		go rpc.New(&rpc.Config{
			Mempool: memPool,
			P2P:     p2pSvc,
			RPCHost: cfg.RPCHost,
		})
	}

	return n, nil
}

func (s *Service) spawnNextBlockMiner(prevBlock *mainchain.Block) error {
	pendingTransactions, err := s.props.Store.GatherPendingTransactions()
	if err != nil {
		log.Errorf("[node] error gathering pending transactions; %v", err)
		return err
	}

	encMinerAddr, err := c3crypto.EncodeAddress(s.props.Keys.Pub)
	if err != nil {
		log.Errorf("[node] error encoding miner address; %v", err)
		return err
	}

	log.Printf("[node] pending tx count: %v", len(pendingTransactions))

	// TODO: need to get this from the network
	blockDifficulty := s.props.BlockDifficulty

	var simulated bool

	// NOTE: if block difficulty is set to 0 than we simulate block hashing (used for testing)
	if s.props.BlockDifficulty == 0 {
		simulated = true
	}

	log.Printf("[node] block difficult level: %v", blockDifficulty)

	ch := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	minerSvc, err := miner.New(&miner.Props{
		Context:             ctx,
		PreviousBlock:       prevBlock,
		Difficulty:          uint64(blockDifficulty),
		Channel:             ch,
		Async:               true, // TODO: need to make this a cli flag
		P2P:                 s.props.P2P,
		Sandbox:             sandbox.New(nil),
		EncodedMinerAddress: encMinerAddr,
		PendingTransactions: pendingTransactions,
		RemoveTx:            s.props.Store.RemoveTx,
		Simulated:           simulated,
	})
	if err != nil {
		log.Errorf("[node] err building miner\n%v", err)
		cancel()

		return err
	}

	if err := minerSvc.SpawnMiner(); err != nil {
		log.Errorf("[node] err spawning miner\n%v", err)
		cancel()

		return err
	}

	return s.spawnMinerListener(cancel, ch)
}

func (s *Service) spawnMinerListener(cancel context.CancelFunc, minerChan chan interface{}) error {
	log.Println("[node] spawned miner listener")

	go func() {
		select {
		case v := <-minerChan:
			{
				switch v.(type) {
				case error:
					err, _ := v.(error)
					log.Errorf("[node] received an error from the miner\n%v", err)

					// just to be safe
					cancel()

					return

				case *miner.MinedBlock:
					log.Println("[node] block mined")
					// just to be safe
					cancel()

					// note: no matter what happens, mine the next block...
					defer func() {
						go func() {
							// TODO: make this recursive and keep trying on err
							nextBlock, err := s.props.Store.GetHeadBlock()
							if err != nil {
								log.Errorf("[node] err getting head block for miner\n%v", err)
								return
							}

							if err := s.spawnNextBlockMiner(&nextBlock); err != nil {
								log.Errorf("[node] error starting miner\n%v", err)
								return
							}
						}()
					}()

					minedBlock, _ := v.(*miner.MinedBlock)

					pendingBlocks, err := s.props.Store.GetPendingMainchainBlocks()
					if err != nil {
						log.Errorf("[node] err checking pending mainchain blocks\n%v", err)
						return
					}

					// TODO: wait until there are no more pending blocks, but for now, just assume all blocks received will be added to the chain
					// TODO: check that all pending blocks have block #'s larger than the one we just mined
					if pendingBlocks != nil && len(pendingBlocks) > 0 {
						log.Printf("[node] we mined a block, but there are mined blocks pending, just abort and wait")
						return
					}

					currentBlock, err := s.props.Store.GetHeadBlock()
					if err != nil {
						log.Errorf("[node] err getting head block\n%v", err)
						return
					}

					eq, err := currentBlock.Equals(minedBlock.PreviousBlock)
					if err != nil {
						log.Errorf("[node] err checking if current block head was the one mined\n%v", err)
						return
					}

					if !eq {
						log.Println(currentBlock, *minedBlock.PreviousBlock, *minedBlock.NextBlock, *minedBlock.NextBlock.Props().BlockHash)
						log.Errorf("[node] the block mined is not built from the current head block\n%v", err)
						return
					}

					sigR, sigS, err := c3crypto.Sign(s.props.Keys.Priv, []byte(*minedBlock.NextBlock.Props().BlockHash))
					if err != nil {
						log.Errorf("[node] err signing mined block\n%v", err)
						return
					}

					nextProps := minedBlock.NextBlock.Props()
					nextProps.MinerSig = &mainchain.MinerSig{
						R: hexutil.EncodeBigInt(sigR),
						S: hexutil.EncodeBigInt(sigS),
					}
					nextBlock := mainchain.New(&nextProps)
					minedBlock.NextBlock = nextBlock

					if err := s.BroadcastMinedBlock(minedBlock); err != nil {
						log.Errorf("[node] err broadcasting mined block\n%s", err)
						return
					}

					go s.checkpointBlock(minedBlock)

					go func() {
						if err := s.setMinedBlockData(minedBlock); err != nil {
							log.Errorf("[node] err setting mined block data\n%v", err)
						}
					}()

					if err := s.props.Store.SetHeadBlock(minedBlock.NextBlock); err != nil {
						log.Errorf("[node] err setting the head block\n%v", err)
						return
					}

					/*
						TODO
						_, err = s.props.P2P.SetLatestBlock(minedBlock.NextBlock)
						if err != nil {
							log.Errorf("[node] err storing the head block\n%v", err)
							return
						}
					*/

					if err := s.removeMinedTxs(minedBlock); err != nil {
						log.Errorf("[node] err removing mined txs\n%v", err)
						return
					}

				default:
					log.Errorf("[node] received message of unknown type from the miner\ntype %T\n%v", v, v)
					// just to be safe
					cancel()

					return
				}
			}
		case <-s.props.CancelMinersChannel:
			{
				// TODO: check if any of the transactions or state image hashes we're mining were included in the new block. If not, we can largely continue
				// with an updated prev bloch hash, and number, etc.
				cancel()

				return
			}
		}
	}()

	return nil
}

func (s *Service) listenForEvents() error {
	if err := s.spawnBlocksListener(); err != nil {
		return err
	}

	return s.spawnTransactionsListener()
}

func (s *Service) spawnBlocksListener() error {
	sub, err := s.props.Pubsub.Subscribe("blocks")
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(s.props.Context)
			if err != nil {
				s.props.SubscriberChannel <- err
				continue
			}

			if peer.ID(msg.GetFrom()).Pretty() == s.props.Host.ID().Pretty() {
				// note: received a message from ourselves
				continue
			}

			block := new(miner.MinedBlock)
			if err := block.Deserialize(msg.GetData()); err != nil {
				s.props.SubscriberChannel <- err
				continue
			}

			s.props.SubscriberChannel <- block
		}
	}()

	return nil
}

func (s *Service) spawnTransactionsListener() error {
	sub, err := s.props.Pubsub.Subscribe("transactions")
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(s.props.Context)
			if err != nil {
				s.props.SubscriberChannel <- err
				continue
			}

			if peer.ID(msg.GetFrom()).Pretty() == s.props.Host.ID().Pretty() {
				// note: received a message from ourselves
				continue
			}

			tx := new(statechain.Transaction)
			if err := tx.Deserialize(msg.GetData()); err != nil {
				s.props.SubscriberChannel <- err
				continue
			}

			s.props.SubscriberChannel <- tx
		}
	}()

	return nil
}

// BroadcastMinedBlock ...
// note: only mainchain blocks get broadcast
func (s *Service) BroadcastMinedBlock(minedBlock *miner.MinedBlock) error {
	if minedBlock == nil {
		return errors.New("cannot broadcast nil block")
	}

	log.Printf("[node] broadcasting the block %s", minedBlock.NextBlock.Props().BlockNumber)
	data, err := minedBlock.Serialize()
	if err != nil {
		return err
	}

	return s.props.Pubsub.Publish("blocks", data)
}

// BroadcastTransaction ...
func (s *Service) BroadcastTransaction(tx *statechain.Transaction) (*nodetypes.SendTxResponse, error) {
	if tx == nil {
		log.Errorln("error; cannot broadcast nil transaction")
		return nil, errors.New("cannot broadcast nil transaction")
	}

	var res nodetypes.SendTxResponse

	data, err := tx.Serialize()
	if err != nil {
		log.Errorf("[node] error serializing transaction; %v", err)
		return nil, err
	}

	if err := s.props.Pubsub.Publish("transactions", data); err != nil {
		log.Errorf("[node] error publishing transaction; %v", err)
		return nil, err
	}

	res.TxHash = tx.Props().TxHash

	log.Printf("[node] transaction %s broadcasted", *tx.Props().TxHash)
	return &res, nil
}

//// GetInfo ...
//func (s *Service) GetInfo() (*nodetypes.GetInfoResponse, error) {
//var res nodetypes.GetInfoResponse

//head, err := s.props.Blockchain.MainHead()
//if err != nil {
//return nil, err
//}

//res.BlockHeight = head.Props().BlockNumber

//return &res, err
//}

func (s *Service) handleReceiptOfMinedBlock(minedBlock *miner.MinedBlock) {
	log.Println("[node] handling receipt of mined block")

	if minedBlock == nil {
		log.Error("[node] received nil block")
		return
	}
	if minedBlock.NextBlock == nil {
		log.Error("[nored] received nil next block")
		return
	}
	if minedBlock.NextBlock.Props().BlockHash == nil {
		log.Error("[node] received block with nil hash")
		return
	}

	log.Println(colorlog.Yellow("[node] received mined block on the channel\nblock number: %s", minedBlock.NextBlock.Props().BlockNumber))

	if err := s.props.Store.SetPendingMainchainBlock(minedBlock.NextBlock); err != nil {
		log.Errorf("[node] err setting pending mainchain block\n%v", err)
		return
	}
	defer func() {
		if err := s.props.Store.RemovePendingMainchainBlock(*minedBlock.NextBlock.Props().BlockHash); err != nil {
			log.Errorf("[node] err removing pending mainchain block\n%v", err)
		}
	}()

	// TODO: check the block explorer to be sure that we haven't already received this block
	// TODO: handle this (and generally all of these) err(ors) better?
	//  1) try again?
	//  2) ping the network to see if other nodes have accepted?
	// note: timeout should be a cli flag
	ctx, cancel := context.WithTimeout(context.Background(), config.MinedBlockVerificationTimeout)
	defer cancel()
	ok, err := miner.VerifyMinedBlock(ctx, s.props.P2P, sandbox.New(nil), minedBlock)
	if err != nil {
		log.Errorf("[node] received err while verifying mined block\nblock: %v\nerr: %v", *minedBlock.NextBlock, err)
		return
	}

	// note: ping the other nodes to tell them we didn't accept the block? See if they did?
	if !ok {
		log.Error("[node] received invalid mined block")
		return
	}
	log.Println("[node] mined block was validated")

	// compare it to the block head that we have
	localHeadBlock, err := s.props.Store.GetHeadBlock()
	if err != nil {
		log.Errorf("[node] err getting our head block\n%v", err)
		return
	}

	localBlockHeight, err := hexutil.DecodeUint64(localHeadBlock.Props().BlockNumber)
	if err != nil {
		log.Errorf("[node] err decoding head block height\n%v", err)
		return
	}
	receivedBlockHeight, err := hexutil.DecodeUint64(minedBlock.NextBlock.Props().BlockNumber)
	if err != nil {
		log.Errorf("[node] err decoding received block height\n%v", err)
		return
	}

	// TODO: if delta(local, received) > 1 then we need to backfill our missing blocks
	if localBlockHeight >= receivedBlockHeight {
		log.Warnf("[node] local block height is %v and received is %v, therefore, not adding block to chain", localBlockHeight, receivedBlockHeight)
		return
	}

	// note: block is valid, keep it
	s.props.CancelMinersChannel <- struct{}{}

	if err := s.props.Store.SetHeadBlock(minedBlock.NextBlock); err != nil {
		log.Errorf("[node] err setting head block in node store\n%v", err)
		return
	}
	if err := s.props.Store.RemovePendingMainchainBlock(*minedBlock.NextBlock.Props().BlockHash); err != nil {
		log.Errorf("[node] err removing pending mainchain block\n%v", err)
		return
	}

	go func() {
		if err := s.setMinedBlockData(minedBlock); err != nil {
			log.Errorf("[node] err setting mined block data\n%v", err)
		}
	}()

	if err := s.removeMinedTxs(minedBlock); err != nil {
		log.Errorf("[node] err removing mined txs\n%v", err)
		return
	}

	// note: start mining the next block, but don't start if there are still pending blocks
	// TODO: if any of the above fails, we may never get here and may be stuck!
	pendingBlocks, err := s.props.Store.GetPendingMainchainBlocks()
	if err != nil {
		log.Errorf("[node] err checking pending mainchain blocks\n%v", err)
		return
	}

	// TODO: check that all pending blocks have block #'s larger than the one we just mined
	if pendingBlocks != nil && len(pendingBlocks) > 0 {
		log.Errorf("[node] blocks pending, don't start mining new block, yet")
		return
	}

	if err := s.spawnNextBlockMiner(minedBlock.NextBlock); err != nil {
		log.Errorf("err starting miner\n%v", err)
		return
	}
}

func (s *Service) handleReceiptOfStatechainTransaction(tx *statechain.Transaction) {
	if tx == nil {
		log.Errorln("[node] received nil tx")
		return
	}

	ok, err := miner.VerifyTransaction(tx)
	if err != nil {
		log.Errorf("[node] err verifying tx: %v\nerr: %v", *tx, err)
		return
	}

	if !ok {
		log.Errorf("[node] received an invalid tx\n%v", *tx)
		return
	}

	// TODO: also check the block explorer to be sure this tx isn't already in a block
	// TODO: check the miner to see if it needs to stop
	// note: verify tx checks that TxHash is not nil
	ok, err = s.props.Store.HasTx(*tx.Props().TxHash)

	if ok {
		log.Printf("[node] tx already in mempool; tx hash: %s", *tx.Props().TxHash)
	}

	if err == nil && !ok {
		if err := s.props.Store.AddTx(tx); err != nil {
			// TODO: need to handle this err better
			log.Errorf("[node] err adding tx to store\n%v", err)
			return
		}

		log.Printf(colorlog.Magenta("[node] tx new added to mempool; tx hash: %s", *tx.Props().TxHash))
	}

	if err != nil {
		log.Errorf("[node] err checking if store has tx\n%v", err)
	}

	if _, err = s.props.P2P.SetStatechainTransaction(tx); err != nil {
		// TODO: need to handle this err better
		log.Errorf("[node] err setting tx: %v\nerr: %v", *tx, err)
		return
	}
}

func (s *Service) setMinedBlockData(minedBlock *miner.MinedBlock) error {
	if minedBlock == nil {
		return errors.New("nil mined block")
	}
	if minedBlock.NextBlock == nil {
		return errors.New("nil next block")
	}
	if minedBlock.NextBlock.Props().BlockHash == nil {
		return errors.New("nil next block block hash")
	}

	blk := minedBlock.NextBlock.Props()

	log.Println(colorlog.Green("[node] storing mined block data\nblock number: %s\nblock hash: %s\nstate blocks merkle hash: %s\nstate chain blocks: %v\ntransactions: %v", blk.BlockNumber, *blk.BlockHash, blk.StateBlocksMerkleHash, len(minedBlock.StatechainBlocksMap), len(minedBlock.TransactionsMap)))

	for _, statechainBlock := range minedBlock.StatechainBlocksMap {
		if statechainBlock == nil {
			log.Errorf("[node] mined block state chain block is nil, continuing")
			continue
		}

		if _, err := s.props.P2P.SetStatechainBlock(statechainBlock); err != nil {
			log.Errorf("[node] error setting state chain block; %v", err)
			return err
		}

		log.Println(colorlog.Green("[node] storing state chain block\nstate chain block number: %s\nstate chain block hash: %s\nstate current hash: %s\ntx hash: %s\nprev state block hash: %s\nprev state diff hash: %s", statechainBlock.Props().BlockNumber, *statechainBlock.Props().BlockHash, statechainBlock.Props().StateCurrentHash, statechainBlock.Props().TxHash, statechainBlock.Props().PrevBlockHash, statechainBlock.Props().StatePrevDiffHash))
	}

	for _, transaction := range minedBlock.TransactionsMap {
		if transaction == nil {
			log.Errorf("[node] mined block transaction is nil, continuing")
			continue
		}

		if _, err := s.props.P2P.SetStatechainTransaction(transaction); err != nil {
			log.Errorf("[node] error setting state chain transaction; %v", err)
			return err
		}
	}

	for _, diff := range minedBlock.DiffsMap {
		if diff == nil {
			log.Errorf("[node] mined block diff is nil, continuing")
			continue
		}

		if _, err := s.props.P2P.SetStatechainDiff(diff); err != nil {
			log.Errorf("[node] error setting state chain diff diff; %v", err)
			return err
		}
	}

	for _, tree := range minedBlock.MerkleTreesMap {
		if tree == nil {
			log.Println("[node] mined block merkle tree is nil, continuing")
			continue
		}

		if _, err := s.props.P2P.SetMerkleTree(tree); err != nil {
			log.Errorf("[node] error setting main chain merkle tree; %v", err)
			return err
		}
	}

	if _, err := s.props.P2P.SetMainchainBlock(minedBlock.NextBlock); err != nil {
		log.Errorf("[node] error setting main chain block; %v", err)
		return err
	}

	return nil
}

func (s *Service) removeMinedTxs(minedBlock *miner.MinedBlock) error {
	log.Println("[node] removing mined transactions for block")
	var txs []string
	for txHash := range minedBlock.TransactionsMap {
		txs = append(txs, txHash)
	}

	return s.props.Store.RemoveTxs(txs)
}

// CheckpointBlock ...
func (s *Service) checkpointBlock(minedBlock *miner.MinedBlock) error {
	if minedBlock == nil {
		return errors.New("cannot checkpoint nil block")
	}

	if s.props.EOSClient != nil {
		blockHash := *minedBlock.NextBlock.Props().BlockHash
		_, err := s.props.EOSClient.CheckpointRoot(blockHash)
		if err != nil {
			return err
		}
	}

	if s.props.EthereumClient != nil {
		blockHash := *minedBlock.NextBlock.Props().BlockHash
		_, err := s.props.EthereumClient.CheckpointRoot(blockHash)
		if err != nil {
			return err
		}
	}

	return nil
}

// Start ...
// NOTE: start is called from cmd package
func (s *Service) Start() error {
	for {
		switch v := <-s.props.SubscriberChannel; v.(type) {
		case error:
			err, _ := v.(error)
			log.Errorf("[node] received an error on the channel %s", err)

		case *miner.MinedBlock:
			log.Print("[node] received mined block")
			b, _ := v.(*miner.MinedBlock)
			go s.handleReceiptOfMinedBlock(b)

		case *statechain.Transaction:
			log.Print("[node] received statechain transaction")
			tx, _ := v.(*statechain.Transaction)
			go s.handleReceiptOfStatechainTransaction(tx)

		default:
			log.Errorf("[node] received an unknown message on channel of type %T\n%v", v, v)
		}
	}
}

// Props ...
func (s *Service) Props() Props {
	return s.props
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
