package node

import (
	"errors"

	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
	nodetypes "github.com/c3systems/c3/node/types"
	peer "github.com/libp2p/go-libp2p-peer"
	//"github.com/c3systems/c3/node/wallet"
	//ipfsaddr "github.com/ipfs/go-ipfs-addr"
	//libp2p "github.com/libp2p/go-libp2p"
	//host "github.com/libp2p/go-libp2p-host"
	//peerstore "github.com/libp2p/go-libp2p-peerstore"
	//floodsub "github.com/libp2p/go-floodsub"
)

// New ...
func New(props *Props) (*Service, error) {
	if props == nil {
		return nil, errors.New("props cannot be nil")
	}
	if props.Store == nil ||
		props.Blockchain == nil || props.Pubsub == nil {
		return nil, errors.New("p2p node, store, blockchain and pubsub are required")
	}

	return &Service{
		props: *props,
	}, nil
}

// Start ...
func (s Service) Start() error {
	if err := s.listenBlocks(); err != nil {
		return err
	}

	return s.listenTransactions()
}

func (s Service) listenBlocks() error {
	sub, err := s.props.Pubsub.Subscribe("blocks")
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(s.props.Context)
			if err != nil {
				s.props.Channel <- err
				continue
			}

			if peer.ID(msg.GetFrom()).Pretty() == s.props.Host.ID().Pretty() {
				// note: received a message from ourselves
				continue
			}

			var block mainchain.Block
			if err := block.Deserialize(msg.GetData()); err != nil {
				s.props.Channel <- err
				continue
			}

			s.props.Channel <- &block

			//log.Println("Block received over network, blockhash", block.Props().BlockHash)
			//cid := node.blockchain.AddMainBlock(&block)
			//if cid != nil {
			//log.Println("Block added, cid:", cid)
			//node.store.RemoveTxs(block.Transactions())
			//}
		}
	}()

	return nil
}

func (s Service) listenTransactions() error {
	sub, err := s.props.Pubsub.Subscribe("transactions")
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(s.props.Context)
			if err != nil {
				s.props.Channel <- err
				continue
			}

			if peer.ID(msg.GetFrom()).Pretty() == s.props.Host.ID().Pretty() {
				// note: received a message from ourselves
				continue
			}

			var tx statechain.Transaction
			if err := tx.Deserialize(msg.GetData()); err != nil {
				s.props.Channel <- err
				continue
			}
			if err := s.props.Store.AddTx(&tx); err != nil {
				s.props.Channel <- err
				continue
			}

			s.props.Channel <- &tx
		}
	}()

	return nil
}

// CreateNewBlock ...
func (s Service) CreateNewBlock() (*mainchain.Block, error) {
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
func (s Service) BroadcastBlock(block *mainchain.Block) error {
	if block == nil {
		return errors.New("cannot broadcast nil block")
	}

	data, err := block.Serialize()
	if err != nil {
		return err
	}

	return s.props.Pubsub.Publish("blocks", data)
}

// BroadcastTransaction ...
func (s Service) BroadcastTransaction(tx *statechain.Transaction) (*nodetypes.SendTxResponse, error) {
	if tx == nil {
		return nil, errors.New("cannot broadcast nil transaction")
	}

	var res nodetypes.SendTxResponse

	data, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	if err := s.props.Pubsub.Publish("transactions", data); err != nil {
		return nil, err
	}

	res.TxHash = tx.Props().TxHash

	return &res, nil
}

// GetInfo ...
func (s Service) GetInfo() (*nodetypes.GetInfoResponse, error) {
	var res nodetypes.GetInfoResponse

	head, err := s.props.Blockchain.MainHead()
	if err != nil {
		return nil, err
	}

	res.BlockHeight = head.Props().BlockNumber

	return &res, err
}
