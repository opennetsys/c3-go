package node

import (
	"context"
	"errors"

	"github.com/c3systems/c3/core/chain"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/node/p2p"
	"github.com/c3systems/c3/core/node/pubsub"
	nodestore "github.com/c3systems/c3/core/node/store"
	nodetypes "github.com/c3systems/c3/core/node/types"
	"github.com/c3systems/c3/core/node/wallet"
	//ipfsaddr "github.com/ipfs/go-ipfs-addr"
	//libp2p "github.com/libp2p/go-libp2p"
	//host "github.com/libp2p/go-libp2p-host"
	//peerstore "github.com/libp2p/go-libp2p-peerstore"
	//floodsub "github.com/libp2p/go-floodsub"
)

// Props ...
type Props struct {
	CTX        context.Context
	CH         chan interface{}
	P2pNode    p2p.Interface
	Store      nodestore.Interface
	Blockchain chain.Interface
	Pubsub     pubsub.Interface
	Wallet     wallet.Interface
}

// Node ...
type Node struct {
	props Props
}

// New ...
func New(props Props) *Node {
	//newNode, err := libp2p.New(props.Ctx, libp2p.Defaults)
	//if err != nil {
	//return nil, err
	//}

	//pubsub, err := floodsub.NewFloodSub(ctx, newNode)
	//if err != nil {
	//return nil, err
	//}

	//for i, addr := range newNode.Addrs() {
	//log.Printf("%d: %s/ipfs/%s\n", i, addr, newNode.ID().Pretty())
	//}

	//if len(os.Args) > 1 {
	//addrstr := os.Args[1]
	//addr, err := ipfsaddr.ParseString(addrstr)
	//if err != nil {
	//return nil, err
	//}
	//log.Println("Parse Address:", addr)

	//pinfo, _ := peerstore.InfoFromP2pAddr(addr.Multiaddr())

	//if err := newNode.Connect(ctx, *pinfo); err != nil {
	//log.Println("bootstrapping a peer failed", err)
	//return nil, err
	//}
	//}

	//blockchain := NewBlockchain(newNode)

	//node.p2pNode = newNode
	//node.mempool = NewMempool()
	//node.pubsub = pubsub
	//node.blockchain = blockchain
	//node.wallet = NewWallet()

	return &Node{
		props: props,
	}
}

// Start ...
func (n Node) Start() error {
	if err := node.listenBlocks(); err != nil {
		return nil, err
	}
	if err := node.listenTransactions(); err != nil {
		return nil, err
	}

	return nil
}

func (n Node) listenBlocks(ch chan interface{}) error {
	sub, err := n.pubsub.Subscribe("blocks")
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(n.props.CTX)
			if err != nil {
				n.props.CH <- err
				continue
			}

			var block mainchain.Block
			if err := block.Deserialize(msg.GetData()); err != nil {
				n.props.CH <- err
				continue
			}

			ch <- &block

			//log.Println("Block received over network, blockhash", block.Props().BlockHash)
			//cid := node.blockchain.AddMainBlock(&block)
			//if cid != nil {
			//log.Println("Block added, cid:", cid)
			//node.store.RemoveTxs(block.Transactions())
			//}
		}
	}()
}

func (n Node) listenTransactions() error {
	sub, err := n.pubsub.Subscribe("transactions")
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(n.props.CTX)
			if err != nil {
				n.props.ch <- err
				continue
			}

			var tx statechain.Transaction
			if err := tx.Deserialize(msg.GetData()); err != nil {
				node.props.CH <- err
				continue
			}
			if err := n.store.AddTx(&tx); err != nil {
				node.props.CH <- err
				continue
			}

			node.props.CH <= &tx
		}
	}()
}

// CreateNewBlock ...
func (n Node) CreateNewBlock() (*mainchain.Block, error) {
	//var blk *mainchan.Block
	//blk.PrevHash = node.blockchain.MainHead().Props().BlockHash
	//blk.Transactions, err = node.store.SelectTransactions()
	//if err != nil {
	//return nil, err
	//}
	//blk.Height = node.blockchain.Head.Height + 1
	//blk.Time = uint64(time.Now().Unix())

	//return &blk

	return nil, nil
}

// BroadcastBlock ...
func (n Node) BroadcastBlock(block *mainchain.Block) error {
	if block == nil {
		return errors.New("cannot broadcast nil block")
	}

	data := block.Serialize()
	n.pubsub.Publish("blocks", data) // note: no err?

	return nil
}

// BroadcastTransaction ...
func (n Node) BroadcastTransaction(tx *statechain.Transaction) (*nodetypes.SendTxResponse, error) {
	if tx == nil {
		return nil, errors.New("cannot broadcast nil transaction")
	}

	var res types.SendTxResponse

	data, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	n.pubsub.Publish("transactions", data)

	res.TxHash = tx.Props().TxHash

	return &res
}

// GetInfo ...
func (n Node) GetInfo() *nodetypes.GetInfoResponse {
	var res types.GetInfoResponse
	res.BlockHeight = n.blockchain.MainHead().Props().BlockNumber

	return &res
}
