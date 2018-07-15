package node

import (
	"errors"
	"log"

	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/mainchain/miner"
	"github.com/c3systems/c3/core/chain/statechain"
	nodetypes "github.com/c3systems/c3/node/types"

	peer "github.com/libp2p/go-libp2p-peer"
)

// New ...
func New(props *Props) (*Service, error) {
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

func (s *Service) setProps(props Props) error {
	if props.Store == nil || props.Pubsub == nil || props.P2P == nil {
		return errors.New("p2p node, store, and pubsub are required")
	}

	s.props = props

	return nil
}

func (s Service) spawnNextBlockMiner(prevBlock *mainchain.Block) error {
	pendingTransactions, err := s.props.Store.GatherPendingTransactions()
	if err != nil {
		return err
	}
	encMinerAddr, err := c3crypto.EncodeAddress(s.props.Keys.Pub)
	if err != nil {
		return err
	}

	isValid := true
	ch := make(chan interface{})
	minerSvc, err := miner.New(&miner.Props{
		IsValid:             &isValid,
		PreviousBlock:       prevBlock,
		Difficulty:          6, // TODO: need to get this from the network
		Channel:             ch,
		Async:               true, // TODO: need to make this a cli flag
		P2P:                 s.props.P2P,
		EncodedMinerAddress: encMinerAddr,
		PendingTransactions: pendingTransactions,
	})
	if err != nil {
		log.Printf("[node] err building miner\n%v", err)
		return err
	}

	if err := minerSvc.SpawnMiner(); err != nil {
		log.Printf("[node] err spawning miner\n%v", err)
		return err
	}

	return s.spawnMinerListener(ch, &isValid)
}

func (s Service) spawnMinerListener(minerChan chan interface{}, isValid *bool) error {
	if isValid == nil {
		return errors.New("nil IsValid")
	}
	if *isValid == false {
		return errors.New("is valid is false")
	}

	go func() {
		select {
		case v := <-minerChan:
			{
				switch v.(type) {
				case error:
					err, _ := v.(error)
					log.Printf("[node] received an error from the miner\n%v", err)
					return

				case *miner.MinedBlock:
					log.Println("[node] block mined")

					// note: no matter what happens, mine the next block...
					defer func() {
						go func() {
							// TODO: make this recursive and keep trying on err
							nextBlock, err := s.props.Store.GetHeadBlock()
							if err != nil {
								log.Printf("[node] err getting head block for miner\n%v", err)
								return
							}

							if err := s.spawnNextBlockMiner(&nextBlock); err != nil {
								log.Printf("[node] err starting miner\n%v", err)
								return
							}
						}()
					}()

					minedBlock, _ := v.(*miner.MinedBlock)

					pendingBlocks, err := s.props.Store.GetPendingMainchainBlocks()
					if err != nil {
						log.Printf("[node] err checking pending mainchain blocks\n%v", err)
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
						log.Printf("[node] err getting head block\n%v", err)
						return
					}

					eq, err := currentBlock.Equals(minedBlock.PreviousBlock)
					if err != nil {
						log.Printf("[node] err checking if current block head was the one mined\n%v", err)
						return
					}

					if !eq {
						log.Println(currentBlock, *minedBlock.PreviousBlock, *minedBlock.NextBlock, *minedBlock.NextBlock.Props().BlockHash)
						log.Printf("[node] the block mined is not built from the current head block\n%v", err)
						return
					}

					sigR, sigS, err := c3crypto.Sign(s.props.Keys.Priv, []byte(*minedBlock.NextBlock.Props().BlockHash))
					if err != nil {
						log.Printf("[node] err signing mined block\n%v", err)
						return
					}
					nextProps := minedBlock.NextBlock.Props()
					nextProps.MinerSig = &mainchain.BlockSig{
						R: hexutil.EncodeBigInt(sigR),
						S: hexutil.EncodeBigInt(sigS),
					}
					nextBlock := mainchain.New(&nextProps)
					minedBlock.NextBlock = nextBlock

					if err := s.BroadcastMinedBlock(minedBlock); err != nil {
						log.Printf("[node] err broadcasting mined block\n%s", err)
						return
					}

					go func() {
						if err := s.setMinedBlockData(minedBlock); err != nil {
							log.Printf("[node] err setting mined block data\n%v", err)
						}
					}()

					if err := s.props.Store.SetHeadBlock(minedBlock.NextBlock); err != nil {
						log.Printf("[node] err setting the head block\n%v", err)
						return
					}

					if err := s.removeMinedTxs(minedBlock); err != nil {
						log.Printf("[node] err removing mined txs\n%v", err)
						return
					}

				default:
					log.Printf("[node] received message of unknown type from the miner\ntype %T\n%v", v, v)
					return
				}
			}
		case <-s.props.CancelMinersChannel:
			{
				// TODO: check if any of the transactions or state image hashes we're mining were included in the new block. If not, we can largely continue
				*isValid = false
				return
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

func (s Service) spawnTransactionsListener() error {
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
func (s Service) BroadcastMinedBlock(minedBlock *miner.MinedBlock) error {
	if minedBlock == nil {
		return errors.New("cannot broadcast nil block")
	}

	log.Println("[node] broadcasting the block")
	data, err := minedBlock.Serialize()
	if err != nil {
		return err
	}

	return s.props.Pubsub.Publish("blocks", data)
}

// BroadcastTransaction ...
func (s *Service) BroadcastTransaction(tx *statechain.Transaction) (*nodetypes.SendTxResponse, error) {
	if tx == nil {
		return nil, errors.New("cannot broadcast nil transaction")
	}

	var res nodetypes.SendTxResponse

	data, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	log.Println("HIIII", s.props, s.props.Pubsub)

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

func (s Service) handleReceiptOfMinedBlock(minedBlock *miner.MinedBlock) {
	if minedBlock == nil {
		log.Println("[node] received nil block")
		return
	}
	if minedBlock.NextBlock == nil {
		log.Println("[nored] received nil next block")
		return
	}
	if minedBlock.NextBlock.Props().BlockHash == nil {
		log.Println("[node] received block with nil hash")
		return
	}

	log.Printf("[node] received mined block on the channel\n%v", *minedBlock)

	if err := s.props.Store.SetPendingMainchainBlock(minedBlock.NextBlock); err != nil {
		log.Printf("[node] err setting pending mainchain block\n%v", err)
		return
	}
	defer func() {
		if err := s.props.Store.RemovePendingMainchainBlock(*minedBlock.NextBlock.Props().BlockHash); err != nil {
			log.Printf("[node] err removing pending mainchain block\n%v", err)
		}
	}()

	// TODO: check the block explorer to be sure that we haven't already received this block
	// TODO: handle this (and generally all of these) err(ors) better?
	//  1) try again?
	//  2) ping the network to see if other nodes have accepted?
	// TODO: implement context.Context rather than a pointer to a bool
	isValid := true
	// TODO: add method to verify mined block
	ok, err := miner.VerifyMinedBlock(s.props.P2P, &isValid, minedBlock)
	if err != nil {
		log.Printf("[node] received err while verifying mined block\nblock: %v\nerr: %v", *minedBlock.NextBlock, err)
		return
	}

	// note: ping the other nodes to tell them we didn't accept the block? See if they did?
	if !ok {
		log.Printf("[node] received invalid mined block\nblock: %v\nerr: %v", *minedBlock, err)
		return
	}
	log.Println("[node] mined block was validated")

	// compare it to the block head that we have
	localHeadBlock, err := s.props.Store.GetHeadBlock()
	if err != nil {
		log.Printf("[node] err getting our head block\n%v", err)
		return
	}

	localBlockHeight, err := hexutil.DecodeUint64(localHeadBlock.Props().BlockNumber)
	if err != nil {
		log.Printf("[node] err decoding head block height\n%v", err)
		return
	}
	receivedBlockHeight, err := hexutil.DecodeUint64(minedBlock.NextBlock.Props().BlockNumber)
	if err != nil {
		log.Printf("[node] err decoding received block height\n%v", err)
		return
	}

	// TODO: if delta(local, received) > 1 then we need to backfill our missing blocks
	if localBlockHeight >= receivedBlockHeight {
		log.Printf("[node] local block height is %v and received is %v, therefore, not adding block to chain", localBlockHeight, receivedBlockHeight)
		return
	}

	// note: block is valid, keep it
	s.props.CancelMinersChannel <- struct{}{}

	if err := s.props.Store.SetHeadBlock(minedBlock.NextBlock); err != nil {
		log.Printf("[node] err setting head block in node store\n%v", err)
		return
	}
	if err := s.props.Store.RemovePendingMainchainBlock(*minedBlock.NextBlock.Props().BlockHash); err != nil {
		log.Printf("[node] err removing pending mainchain block\n%v", err)
		return
	}

	go func() {
		if err := s.setMinedBlockData(minedBlock); err != nil {
			log.Printf("[node] err setting mined block data\n%v", err)
		}
	}()

	if err := s.removeMinedTxs(minedBlock); err != nil {
		log.Printf("[node] err removing mined txs\n%v", err)
		return
	}

	// note: start mining the next block, but don't start if there are still pending blocks
	// TODO: if any of the above fails, we may never get here and may be stuck!
	pendingBlocks, err := s.props.Store.GetPendingMainchainBlocks()
	if err != nil {
		log.Printf("[node] err checking pending mainchain blocks\n%v", err)
		return
	}

	// TODO: check that all pending blocks have block #'s larger than the one we just mined
	if pendingBlocks != nil && len(pendingBlocks) > 0 {
		log.Printf("[node] blocks pending, don't start mining new block, yet")
		return
	}

	if err := s.spawnNextBlockMiner(minedBlock.NextBlock); err != nil {
		log.Printf("err starting miner\n%v", err)
		return
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

	if _, err = s.props.P2P.SetStatechainTransaction(tx); err != nil {
		// TODO: need to handle this err better
		log.Printf("[node] err setting tx: %v\nerr: %v", *tx, err)
		return
	}
}

func (s Service) setMinedBlockData(minedBlock *miner.MinedBlock) error {
	if minedBlock == nil {
		return errors.New("nil mined block")
	}
	if minedBlock.NextBlock == nil {
		return errors.New("nil next block")
	}
	if minedBlock.NextBlock.Props().BlockHash == nil {
		return errors.New("nil next block block hash")
	}

	for _, statechainBlock := range minedBlock.StatechainBlocksMap {
		if statechainBlock == nil {
			continue
		}

		if _, err := s.props.P2P.SetStatechainBlock(statechainBlock); err != nil {
			return err
		}
	}

	for _, transaction := range minedBlock.TransactionsMap {
		if transaction == nil {
			continue
		}

		if _, err := s.props.P2P.SetStatechainTransaction(transaction); err != nil {
			return err
		}
	}

	for _, diff := range minedBlock.DiffsMap {
		if diff == nil {
			continue
		}

		if _, err := s.props.P2P.SetStatechainDiff(diff); err != nil {
			return err
		}
	}

	for _, tree := range minedBlock.MerkleTreesMap {
		if tree == nil {
			continue
		}

		if _, err := s.props.P2P.SetMerkleTree(tree); err != nil {
			return err
		}
	}

	if _, err := s.props.P2P.SetMainchainBlock(minedBlock.NextBlock); err != nil {
		return err
	}

	return nil
}

func (s Service) removeMinedTxs(minedBlock *miner.MinedBlock) error {
	var txs []string
	for txHash := range minedBlock.TransactionsMap {
		txs = append(txs, txHash)
	}

	return s.props.Store.RemoveTxs(txs)
}
