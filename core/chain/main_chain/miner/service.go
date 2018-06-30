package miner

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/node"
	"github.com/c3systems/c3/utils/hexutil"
)

type stateBlocksMutex struct {
	mut    sync.Mut
	blocks []*statechain.Block
}

// Props is passed to the new function
type Props struct {
	Node       node.Interface
	Difficulty uint
}

// Service implements the interface
type Service struct {
	props Props
	ch    chan interface{}
}

// New returns a new service
func New(props *Props) (*Service, error) {
	if props == nil {
		return nil, errors.New("props are required")
	}

	return &Service{
		props: *props,
		ch:    make(chan interface{}),
	}, props.Node.Ping()
}

// Props returns the props
func (s Service) Props() Props {
	return s.props
}

// StartMiner ...
func (s Service) StartMiner() {
	// TODO: interrupt if block received
	for {
		block, err := s.buildMainBlock()
		if err != nil {
			log.Println("[miner] err mining: %v", err)
		}

		s.props.Node.BroadcastBlock(block)
	}
}

func (s Service) buildMainBlock() (*mainchain.Block, error) {
	var (
		wg             sync.WaitGroup
		stateBlocksMut stateBlocksMutex
	)

	// 1. gather tx's
	// TODO: only choose high value tx's to mine
	txsMap, err := s.props.Node.GatherTransactions()
	if err != nil {
		return nil, err
	}

	// 2. apply txs
	for imageHash, transactions := range txsMap {
		wg.Add(1)
		go func(iHash, txs) {
			defer wg.Done()

			block, err := statechain.BuildNextState(iHash, txs)
			if err != nil {
				log.Printf("[miner] err mining state block for hash %s transactions %v: %v", iHash, txs, err)
				return
			}

			stateBlocksMut.mut.Lock()
			stateBlocksMut.blocks = append(stateMut.blocks, block)
			stateBlocksMut.mut.Unlock()
		}(imageHash, transactions)
	}
	wg.Wait()

	// 3. mine main block
	return s.mineBlock(stateBlocksMut.blocks)
}

func (s Service) mineBlock(stateBlocks []*statechain.Block) (*mainchain.Block, error) {
	// TODO: kill if next block is received from network
	// TODO: timeout
	for {
		block, err := mainchain.NewFromStateBlocks(stateBlocks)
		if err != nil {
			return nil, err
		}

		hash, nonce, err := s.generateHash(block)
		if err != nil {
			return err
		}

		check, err := s.checkHashAgainstDifficulty(hash)
		if err != nil {
			return nil, err
		}

		if check {
			return s.buildNextBlock(block, hash, nonce)
		}
	}
}

func (s Service) generateHash(block *mainchain.Block) (string, string, error) {
	nonce, err := s.generateNonce()
	if err != nil {
		return "", "", err
	}

	props := block.Props()
	props.Nonce = nonce

	hash, err := mainChain.HashProps(props)
	return hash, props, err
}

func (s Service) generateNonce() (string, error) {
	bytes := make([]byte, 32)
	if err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (s Service) checkHashAgainstDifficulty(hash string) (bool, error) {
	if len(hash) <= s.props.Difficulty {
		return false, errors.New("generated a hash less than difficulty length")
	}

	for i := 0; i < s.props.Difficulty; i++ {
		if hash[i:i+1] != "0" {
			return false, nil
		}
	}

	return true, nil
}

func (s Service) buildNextBlock(block *mainchain.Block, hash, nonce string) (*mainchain.Block, error) {
	nextProps := block.Props()

	prevBlock, err := node.FetchHeadBlock()
	if err != nil {
		return nil, err
	}
	if prevBlock == nil {
		return nil, errors.New("received nil block head")
	}

	prevProps := prevBlock.Props()

	nextProps.BlockHash = hash
	nextProps.PrevBlockHash = prevProps.BlockHash
	blockHeight, err := hexutil.DecodeUint64(prevProps.BlockNumber)
	if err != nil {
		return nil, err
	}
	nextProps.BlockNumber = hexutil.EncodeUint64(blockHeight + 1)
	nextProps.BlockTime = hexutil.EncodeUint64(time.Now.Unix())
	nextProps.Nonce = nonce
	nextProps.Difficulty = s.props.Difficulty

	return mainchain.New(nextProps), nil
}
