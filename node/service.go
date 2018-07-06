package node

import (
	"errors"

	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
	nodetypes "github.com/c3systems/c3/node/types"
	//"github.com/c3systems/c3/node/wallet"
	//ipfsaddr "github.com/ipfs/go-ipfs-addr"
	//libp2p "github.com/libp2p/go-libp2p"
	//host "github.com/libp2p/go-libp2p-host"
	//peerstore "github.com/libp2p/go-libp2p-peerstore"
	//floodsub "github.com/libp2p/go-floodsub"
)

// New ...
func New(props *nodetypes.Props) (*nodetypes.Service, error) {
	if props == nil {
		return nil, errors.New("props cannot be nil")
	}
	if props.P2pNode == nil || props.Store == nil ||
		props.Blockchain == nil || props.Pubsub == nil {
		return nil, errors.New("p2p node, store, blockchain and pubsub are required")
	}

	return &nodetyps.Service{
		props: *props,
	}, nil
}

// Start ...
func (s nodetypes.Service) Start() error {
	if err := s.listenBlocks(); err != nil {
		return nil, err
	}
	if err := s.listenTransactions(); err != nil {
		return nil, err
	}

	return nil
}

func (s nodetypes.Service) listenBlocks() error {
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

func (s nodetypes.Service) listenTransactions() error {
	sub, err := s.pubsub.Subscribe("transactions")
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
func (s nodetypes.Service) CreateNewBlock() (*mainchain.Block, error) {
	//var blk *mainchan.Block
	//blk.PrevHash = s.props.blockchain.MainHead().Props().BlockHash
	//blk.Transactions, err = s.props.store.SelectTransactions()
	//if err != nil {
	//return nil, err
	//}
	//blk.Height = s.props.blockchain.Head.Height + 1
	//blk.Time = uint64(time.Now().Unix())

	//return &blk

	return nil, nil
}

// BroadcastBlock ...
// note: only mainchain blocks get broadcast
func (s nodetypes.Service) BroadcastBlock(block *mainchain.Block) error {
	if block == nil {
		return errors.New("cannot broadcast nil block")
	}

	data := block.Serialize()
	s.props.pubsub.Publish("blocks", data) // note: no err?

	return nil
}

// BroadcastTransaction ...
func (s nodetypes.Service) BroadcastTransaction(tx *statechain.Transaction) (*nodetypes.SendTxResponse, error) {
	if tx == nil {
		return nil, errors.New("cannot broadcast nil transaction")
	}

	var res types.SendTxResponse

	data, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	s.props.pubsub.Publish("transactions", data)

	res.TxHash = tx.Props().TxHash

	return &res
}

// GetInfo ...
func (s nodetypes.Service) GetInfo() *nodetypes.GetInfoResponse {
	var res types.GetInfoResponse
	res.BlockHeight = s.props.blockchain.MainHead().Props().BlockNumber

	return &res
}
