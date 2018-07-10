package miner

import (
	"crypto/rand"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/p2p"
)

// New returns a new service
func New(props *Props) (*Service, error) {
	if props == nil {
		return nil, errors.New("props are required")
	}

	return &Service{
		props: *props,
	}, nil
}

// Props returns the props
func (s Service) Props() Props {
	return s.props
}

// Start ...
// note: this is blocking and needs to be run in a go routine!
func (s Service) Start() error {
	// TODO: interrupt if block received
	// TODO: write to store and publish on ipfs
	// TODO: reward ourselves with some coin
	go func() {
		for {
			var (
				block *mainchain.Block
				err   error
			)

			switch s.props.Async {
			case true:
				block, err = s.buildMainchainBlockAsync()
				if err != nil {
					s.props.Channel <- err
					continue
				}

			default:
				block, err = s.buildMainchainBlock()
				if err != nil {
					s.props.Channel <- err
					continue
				}
			}

			if block == nil {
				s.props.Channel <- errors.New("built a nil block")
				continue
			}

			s.props.Channel <- block
			//s.props.Node.BroadcastBlock(block)
		}
	}()

	return nil
}

func (s Service) buildMainchainBlockAsync() (*mainchain.Block, error) {
	var (
		wg             sync.WaitGroup
		stateBlocksMut stateBlocksMutex
	)

	// 1. gather tx's
	// TODO: only choose high value tx's to mine
	txs, err := s.props.GatherTransactions()
	if err != nil {
		return nil, err
	}

	txsMap := BuildTxsMap(txs)

	// 2. apply txs
	for imageHash, transactions := range txsMap {
		wg.Add(1)
		go func(iHash string, txs []*statechain.Transaction) {
			defer wg.Done()

			block, err := BuildNextState(iHash, txs)
			if err != nil {
				log.Printf("[miner] err mining state block for hash %s transactions %v: %v", iHash, txs, err)
				return
			}

			stateBlocksMut.mut.Lock()
			stateBlocksMut.blocks = append(stateBlocksMut.blocks, block)
			stateBlocksMut.mut.Unlock()
		}(imageHash, transactions)
	}
	wg.Wait()

	// 3. mine main block
	return s.mineBlock(stateBlocksMut.blocks)
}

func (s Service) buildMainchainBlock() (*mainchain.Block, error) {
	var stateBlocks []*statechain.Block

	// 1. gather tx's
	// TODO: only choose high value tx's to mine
	txs, err := s.props.GatherTransactions()
	if err != nil {
		return nil, err
	}

	txsMap := BuildTxsMap(txs)

	// 2. apply txs
	for imageHash, transactions := range txsMap {
		block, err := BuildNextState(imageHash, transactions)
		if err != nil {
			log.Printf("[miner] err mining state block for hash %s transactions %v: %v", imageHash, transactions, err)
			continue
		}

		stateBlocks = append(stateBlocks, block)
	}

	// 3. mine main block
	return s.mineBlock(stateBlocks)
}

func (s Service) mineBlock(stateBlocks []*statechain.Block) (*mainchain.Block, error) {
	// TODO: kill if next block is received from network
	// TODO: timeout
	for {
		block, err := NewFromStateBlocks(stateBlocks)
		if err != nil {
			return nil, err
		}

		hash, nonce, err := s.generateHashAndNonce(block)
		if err != nil {
			return nil, err
		}

		check, err := CheckHashAgainstDifficulty(hash, block.Props().Difficulty)
		if err != nil {
			return nil, err
		}

		if check {
			return s.buildNextBlock(block, hash, nonce)
		}
	}
}

func (s Service) generateHashAndNonce(block *mainchain.Block) (string, string, error) {
	nonce, err := s.generateNonce()
	if err != nil {
		return "", "", err
	}

	tmpProps := block.Props()
	tmpProps.Nonce = nonce
	tmpBlock := mainchain.New(&tmpProps)

	hash, err := tmpBlock.CalcHash()
	return hash, nonce, err
}

func (s Service) generateNonce() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(bytes)), nil
}

func (s Service) buildNextBlock(block *mainchain.Block, hash, nonce string) (*mainchain.Block, error) {
	nextProps := block.Props()

	prevBlock, err := s.props.FetchHeadBlock()
	if err != nil {
		return nil, err
	}
	if prevBlock == nil {
		return nil, errors.New("received nil block head")
	}

	prevProps := prevBlock.Props()
	if prevProps.BlockHash == nil {
		return nil, errors.New("previous block's block hash is nil")
	}

	nextProps.BlockHash = &hash
	// note: checked for nil block hash, above
	nextProps.PrevBlockHash = *prevProps.BlockHash
	blockHeight, err := hexutil.DecodeUint64(prevProps.BlockNumber)
	if err != nil {
		return nil, err
	}
	nextProps.BlockNumber = hexutil.EncodeUint64(blockHeight + 1)
	nextProps.BlockTime = hexutil.EncodeUint64(uint64(time.Now().Unix()))
	nextProps.Nonce = nonce
	nextProps.Difficulty = hexutil.EncodeUint64(s.props.Difficulty)

	return mainchain.New(&nextProps), nil
}

// VerifyMainchainBlock verifies a mainchain block
// TODO: check block time
// TODO: fetch and check previous block hash
func (s Service) VerifyMainchainBlock(block *mainchain.Block) (bool, error) {
	if block == nil {
		return false, errors.New("block is nil")
	}

	if block.Props().BlockHash == nil {
		return false, errors.New("block hash is nil")
	}

	if mainchain.ImageHash != block.Props().ImageHash {
		return false, nil
	}

	if block.Props().MinerSig == nil {
		return false, nil
	}

	ok, err := CheckBlockHashAgainstDifficulty(block)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	// hash must verify
	tmpHash, err := block.CalcHash()
	if err != nil {
		return false, err
	}
	// note: already checked for nil hash
	if *block.Props().BlockHash != tmpHash {
		return false, nil
	}

	// the sig must verify
	pub, err := PubFromBlock(block)
	if err != nil {
		return false, err
	}

	// note: checked for nil sig, above
	sigR, err := hexutil.DecodeBigInt(block.Props().MinerSig.R)
	if err != nil {
		return false, err
	}
	sigS, err := hexutil.DecodeBigInt(block.Props().MinerSig.S)
	if err != nil {
		return false, err
	}

	// note: nil blockhash was checked, above
	ok, err = c3crypto.Verify(pub, []byte(*block.Props().BlockHash), sigR, sigS)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	// TODO: do in go funcs
	for _, stateblockHash := range block.Props().StateBlockHashes {
		if stateblockHash == nil {
			return false, nil
		}

		stateblockCid, err := p2p.GetCIDByHash(*stateblockHash)
		if err != nil {
			return false, err
		}

		stateblock, err := s.props.P2P.GetStatechainBlock(stateblockCid)
		if err != nil {
			return false, err
		}

		ok, err := s.VerifyStatechainBlock(stateblock)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

// VerifyStatechainBlock verifies a block
// TODO: check timestamp?
func (s Service) VerifyStatechainBlock(block *statechain.Block) (bool, error) {
	if block == nil {
		return false, ErrNilBlock
	}

	// 1. block must have a hash
	if block.Props().BlockHash == nil {
		return false, nil
	}

	// TODO: check the block # and StatePrevDiffHash

	// 2. verify the block hash
	tmpHash, err := block.CalcHash()
	if err != nil {
		return false, err
	}
	// note: checked nil BlockHash, above
	if tmpHash != *block.Props().BlockHash {
		return false, nil
	}

	// 3. verify each tx in the block
	// TODO: do in go funcs
	var txs []*statechain.Transaction
	for _, txHash := range block.Props().TxHashes {
		txCid, err := p2p.GetCIDByHash(txHash)
		if err != nil {
			return false, err
		}

		tmpTx, err := s.props.P2P.GetStatechainTransaction(txCid)
		if err != nil {
			return false, err
		}

		ok, err := VerifyTransaction(tmpTx)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
		txs = append(txs, tmpTx)
	}

	// note: just printing to keep the txs var alive
	log.Println(txs)

	// 4. run the txs through the container
	// TODO: step #4

	// 5. verify the statehash and prev diff hash
	// TODO step #5

	return true, nil
}
