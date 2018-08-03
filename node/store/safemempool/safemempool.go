package safemempool

import (
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/statechain"
)

type poolMut struct {
	mut  sync.Mutex
	pool map[string]string
}

// Props ...
type Props struct {
}

// Service ...
type Service struct {
	props            Props
	txPoolMut        *poolMut
	pendingBlocksMut *poolMut
	headBlock        *mainchain.Block
}

// New ...
func New(props *Props) (*Service, error) {
	// 1. check props
	if props == nil {
		return nil, errors.New("props cannot be nil")
	}

	// 2. build the mut
	txPool := make(map[string]string)
	txMut := poolMut{
		mut:  sync.Mutex{},
		pool: txPool,
	}
	pendingBlocksPool := make(map[string]string)
	pendingBlocksMut := poolMut{
		mut:  sync.Mutex{},
		pool: pendingBlocksPool,
	}

	// 3. return service
	return &Service{
		props:            *props,
		txPoolMut:        &txMut,
		pendingBlocksMut: &pendingBlocksMut,
	}, nil
}

// Props ...
func (s *Service) Props() Props {
	return s.props
}

// HasTx ...
func (s *Service) HasTx(hash string) (bool, error) {
	s.txPoolMut.mut.Lock()
	defer s.txPoolMut.mut.Unlock()
	_, ok := s.txPoolMut.pool[buildKey(hash)]

	return ok, nil
}

// GetTx ...
func (s *Service) GetTx(hash string) (*statechain.Transaction, error) {
	s.txPoolMut.mut.Lock()
	defer s.txPoolMut.mut.Unlock()
	byteStr := s.txPoolMut.pool[buildKey(hash)]

	if byteStr == "" {
		return nil, nil
	}

	tx := new(statechain.Transaction)
	err := tx.DeserializeString(byteStr)

	return tx, err
}

// GetTxs ...
func (s *Service) GetTxs(hashes []string) ([]*statechain.Transaction, error) {
	var txs []*statechain.Transaction

	s.txPoolMut.mut.Lock()
	defer s.txPoolMut.mut.Unlock()

	keys := buildKeys(hashes)

	for _, key := range keys {
		byteStr := s.txPoolMut.pool[key]
		if byteStr != "" {
			tx := new(statechain.Transaction)
			if err := tx.DeserializeString(byteStr); err != nil {
				return nil, err
			}

			txs = append(txs, tx)
		}
	}

	return txs, nil
}

// RemoveTx ...
func (s *Service) RemoveTx(hash string) error {
	s.txPoolMut.mut.Lock()
	defer s.txPoolMut.mut.Unlock()
	delete(s.txPoolMut.pool, buildKey(hash))

	return nil
}

// RemoveTxs ...
func (s *Service) RemoveTxs(hashes []string) error {
	s.txPoolMut.mut.Lock()
	defer s.txPoolMut.mut.Unlock()

	keys := buildKeys(hashes)

	for _, key := range keys {
		delete(s.txPoolMut.pool, key)
	}

	return nil
}

// AddTx ...
func (s *Service) AddTx(tx *statechain.Transaction) error {
	if tx == nil {
		return errors.New("cannot add a nil transaction")
	}
	if tx.Props().TxHash == nil {
		return errors.New("nil tx hash")
	}

	bytesStr, err := tx.SerializeString()
	if err != nil {
		return err
	}

	hash := tx.Props().TxHash
	s.txPoolMut.mut.Lock()
	defer s.txPoolMut.mut.Unlock()
	s.txPoolMut.pool[buildKey(*hash)] = bytesStr

	return nil
}

// GatherPendingTransactions ...
func (s *Service) GatherPendingTransactions() ([]*statechain.Transaction, error) {
	log.Println("[mempool] gathering pending transactions")
	s.txPoolMut.mut.Lock()
	defer s.txPoolMut.mut.Unlock()

	txs := make([]*statechain.Transaction, len(s.txPoolMut.pool), len(s.txPoolMut.pool))
	idx := 0

	log.Printf("[mempool] tx pool size; %v", len(s.txPoolMut.pool))
	for _, byteStr := range s.txPoolMut.pool {
		log.Printf("[mempool] tx pool byte str size is %v", len(byteStr))
		tx := new(statechain.Transaction)
		if err := tx.DeserializeString(byteStr); err != nil {
			return nil, err
		}

		txs[idx] = tx
		idx++
	}

	return txs, nil
}

// GetHeadBlock ...
func (s *Service) GetHeadBlock() (mainchain.Block, error) {
	if s.headBlock == nil {
		return mainchain.Block{}, errors.New("no headblock")
	}

	return *s.headBlock, nil
}

// SetHeadBlock ...
func (s *Service) SetHeadBlock(block *mainchain.Block) error {
	s.headBlock = block
	return nil
}

// SetPendingMainchainBlock ...
func (s *Service) SetPendingMainchainBlock(block *mainchain.Block) error {
	if block == nil {
		return errors.New("block is nil")
	}

	if block.Props().BlockHash == nil {
		return errors.New("block hash is nil")
	}

	encodedString, err := block.SerializeString()
	if err != nil {
		return err
	}

	s.pendingBlocksMut.mut.Lock()
	defer s.pendingBlocksMut.mut.Unlock()
	// note: already checked for nil has, above
	s.pendingBlocksMut.pool[*block.Props().BlockHash] = encodedString

	return nil
}

// GetPendingMainchainBlocks ...
func (s *Service) GetPendingMainchainBlocks() ([]*mainchain.Block, error) {
	var pendingBlocks []*mainchain.Block
	s.pendingBlocksMut.mut.Lock()
	defer s.pendingBlocksMut.mut.Unlock()
	for _, encodedString := range s.pendingBlocksMut.pool {
		block := new(mainchain.Block)
		if err := block.DeserializeString(encodedString); err != nil {
			return nil, err
		}

		pendingBlocks = append(pendingBlocks, block)
	}

	return pendingBlocks, nil
}

// RemovePendingMainchainBlock ...
func (s *Service) RemovePendingMainchainBlock(blockHash string) error {
	s.pendingBlocksMut.mut.Lock()
	defer s.pendingBlocksMut.mut.Unlock()

	delete(s.pendingBlocksMut.pool, blockHash)

	return nil
}

// RemovePendingMainchainBlocks ...
func (s *Service) RemovePendingMainchainBlocks(blockHashes []string) error {
	s.pendingBlocksMut.mut.Lock()
	defer s.pendingBlocksMut.mut.Unlock()

	for _, blockHash := range blockHashes {
		delete(s.pendingBlocksMut.pool, blockHash)
	}

	return nil
}
