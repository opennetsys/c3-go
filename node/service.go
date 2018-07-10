package node

import (
	"errors"
	"log"

	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/mainchain/miner"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/p2p"
	nodetypes "github.com/c3systems/c3/node/types"
	peer "github.com/libp2p/go-libp2p-peer"
)

// New ...
func New(props *Props) (*Service, error) {
	if props == nil {
		return nil, errors.New("props cannot be nil")
	}
	if props.Store == nil ||
		props.Miner == nil || props.Pubsub == nil {
		return nil, errors.New("p2p node, store, miner and pubsub are required")
	}

	return &Service{
		props: *props,
	}, nil
}

func (s Service) startMiner(p2pSvc p2p.Interface) error {
	if err := s.props.Miner.Start(); err != nil {
		return err
	}

	return s.spawnMinerListener()
}

func (s Service) spawnMinerListener() error {
	go func() {
		for {
			switch v := <-s.props.Channel; v.(type) {
			case error:
				err, _ := v.(error)
				log.Printf("received an error from the miner\n%v", err)

			case *mainchain.Block:
				block, _ := v.(*mainchain.Block)
				// TODO: do something with the block
				log.Println(block)

			default:
				log.Printf("received message of unknown type from the miner\ntype %T\n%v", v, v)
			}
		}
	}()

	return nil
}

func (s Service) listenForEvents() error {
	if err := s.spawnBlocksListener(); err != nil {
		return err
	}

	return s.spawnTransactionsListener()
}

func (s Service) spawnBlocksListener() error {
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
		}
	}()

	return nil
}

func (s Service) spawnTransactionsListener() error {
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

//// GetInfo ...
//func (s Service) GetInfo() (*nodetypes.GetInfoResponse, error) {
//var res nodetypes.GetInfoResponse

//head, err := s.props.Blockchain.MainHead()
//if err != nil {
//return nil, err
//}

//res.BlockHeight = head.Props().BlockNumber

//return &res, err
//}

func (s Service) handleReceiptOfMainchainBlock(block *mainchain.Block) {
	if block == nil {
		log.Println("[node] received nil block")
		return
	}

	// TODO: check the block explorer to be sure that we haven't already received this block
	// TODO: check that this is the next block
	// TODO: stop the miner

	// TODO: handle this err better?
	//  1) try again?
	//  2) ping the network to see if other nodes have accepted?
	ok, err := s.props.Miner.VerifyMainchainBlock(block)
	if err != nil {
		log.Printf("[node] received err while verifying mainchain block\nblock: %v\nerr: %v", *block, err)
		return
	}

	// note: ping the other nodes to tell them we didn't accept the block? See if they did?
	if !ok {
		log.Printf("[node] received invalid mainnchain block\nblock: %v\nerr: %v", *block, err)
		return
	}

	// note: just want to set locally?
	_, err = s.props.Miner.Props().P2P.Set(block)
	if err != nil {
		// TODO: need to handle this err better
		log.Printf("[node] err setting main chain block: %v\nerr: %v", *block, err)
		return
	}

	// TODO: add cid to the block explorer
	// TODO: do the inner loop in a go function; and generally move these all to functions
	for _, stateblockHash := range block.Props().StateBlockHashes {
		if stateblockHash == nil {
			log.Println("got nil stateblock hash")
			return
		}

		stateblockCid, err := p2p.GetCIDByHash(*stateblockHash)
		if err != nil {
			// TODO: handle this err better
			log.Printf("[node] err getting statechain block cid\n%v", err)
			continue
		}

		stateblock, err := s.props.Miner.Props().P2P.GetStatechainBlock(stateblockCid)
		if err != nil {
			// TODO: handle this err better
			log.Printf("[node] err getting statechain block\n%v", err)
			continue
		}

		// note: just want to set locally?
		if _, err := s.props.Miner.Props().P2P.Set(stateblock); err != nil {
			// TODO: handle this err better
			log.Printf("[node] err setting statechain block\n%v", err)
			// note: don't continue, here, because we want to get and set all the tx's and diffs
		}

		// TODO: do in a goroutine; and in a new func
		for _, txHash := range stateblock.Props().TxHashes {
			ok, err := s.props.Store.HasTx(txHash)
			if err == nil && ok {
				if err := s.props.Store.RemoveTx(txHash); err != nil {
					// TODO: need to handle this err better
					log.Printf("[node] err removing tx from store\n%v", err)
				}
			}
			if err != nil {
				log.Printf("[node] err checking if store hash tx hash: %s\nerr: %v", txHash, err)
			}

			txCid, err := p2p.GetCIDByHash(txHash)
			if err != nil {
				// TODO: handle this err better
				log.Printf("[node] err getting statechain transaction cid\n%v", err)
				continue
			}

			tx, err := s.props.Miner.Props().P2P.GetStatechainTransaction(txCid)
			if err != nil {
				// TODO: handle this err better
				log.Printf("[node] err getting statechain transaction\n%v", err)
				continue
			}

			// note: just want to set locally?
			if _, err := s.props.Miner.Props().P2P.Set(tx); err != nil {
				// TODO: handle this err better
				log.Printf("[node] err setting statechain tx\n%v", err)
				// note: don't continue, here, because we want to remove the tx from our available pool
			}
		}

		// TODO: do in a goroutine and in a new func
		diffCid, err := p2p.GetCIDByHash(stateblock.Props().StatePrevDiffHash)
		if err != nil {
			// TODO: handle this err better
			log.Printf("[node] err getting statechain diff cid\n%v", err)
			continue
		}

		d, err := s.props.Miner.Props().P2P.GetStatechainDiff(diffCid)
		if err != nil {
			// TODO: handle this err better
			log.Printf("[node] err getting statechain diff\n%v", err)
			continue
		}

		// note: just want to set locally?
		if _, err := s.props.Miner.Props().P2P.Set(d); err != nil {
			// TODO: handle this err better
			log.Printf("[node] err setting statechain diff\n%v", err)
			continue
		}
	}
}

func (s Service) handleReceiptOfStatechainTransaction(tx *statechain.Transaction) {
	if tx == nil {
		return
	}

	ok, err := miner.VerifyTransaction(tx)
	if err != nil {
		log.Printf("[node] err verifying tx: %v\nerr: %v", *tx, err)
		return
	}
	if !ok {
		log.Printf("[node] received an invalid tx\n%v", *tx)
		return
	}

	// TODO: also check the block explorer to be sure this tx isn't already in a block
	// TODO: check the miner to see if it needs to stop
	// note: verify tx checks that TxHash is not nil
	ok, err = s.props.Store.HasTx(*tx.Props().TxHash)
	if err == nil && ok {
		if err := s.props.Store.AddTx(tx); err != nil {
			// TODO: need to handle this err better
			log.Printf("[node] err adding tx to store\n%v", err)
			return
		}
	}
	if err != nil {
		log.Printf("[node] err checking if store has tx\n%v", err)
	}

	// note: just want to set locally?
	_, err = s.props.Miner.Props().P2P.SetStatechainTransaction(tx)
	if err != nil {
		// TODO: need to handle this err better
		log.Printf("[node] err setting tx: %v\nerr: %v", *tx, err)
		return
	}

	// TODO: add cid to the block explorer
}

// TODO: everything
func (s Service) fetchHeadBlock() (*mainchain.Block, error) {
	return nil, nil
}
