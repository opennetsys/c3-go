package safemempool

import (
	"errors"
	"sync"
)

const (
	transactionsMembersName = "transactions"
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
	props   Props
	poolMut poolMut
}

// New ...
func New(props *Props) (*Service, error) {
	// 1. check props
	if props == nil {
		return nil, errors.New("props cannot be nil")
	}

	// 2. build the mut
	pool := make(map[string]string)
	pMut := poolMut{
		mut:  sync.Mutex{},
		pool: pool,
	}

	// 3. return service
	return &Service{
		props:   *props,
		poolMut: pMut,
	}, err
}

// Props ...
func (s Service) Props() Props {
	return s.props
}

func (s Service) HasTx(hash string) (bool, error) {
	s.poolMut.mut.Lock()
	_, ok := s.poolMut.pool[buildKey(hash)]
	s.poolMut.mut.Unlock()

	return ok, nil
}

func (s Service) GetTx(hash string) (*statechain.Transaction, error) {
	s.poolMut.mut.Lock()
	byteStr := s.poolMut.pool[buildKey(hash)]
	s.poolMut.mut.Unlock()

	if byteStr == "" {
		return nil, nil
	}

	var tx statechain.Transaction
	err := tx.DeserializeString(bytesStr)

	return &tx, err
}

func (s Service) GetTxs(hashes []string) ([]*statechain.Transaction, error) {
	var txs []*statechain.Transaction

	s.poolMut.mut.Lock()
	defer s.poolMut.mut.Unlock()

	keys := buildKeys(hashes)

	for _, key := range keys {
		byteStr := s.poolMut.pool[key]
		if byteStr != "" {
			var tx statechain.Transaction
			if err := tx.DeserializeString(bytesStr); err != nil {
				return nil, err
			}

			txs = append(txs, &tx)
		}
	}

	return txs, nil
}

func (s Service) RemoveTx(hash string) error {
	s.poolMut.mut.Lock()
	delete(s.poolMut.pool[buildKey(hash)])
	s.poolMut.mut.Unlock()

	return nil
}

// RemoveTxs ...
func (s Service) RemoveTxs(hashes []string) error {
	s.poolMut.mut.Lock()
	defer s.poolMut.mut.Unlock()

	keys := buildKeys(hashes)

	for _, key := range keys {
		delete(s.poolMut.pool[key])
	}

	return nil
}

// AddTx ...
func (s Service) AddTx(tx *statechain.Transaction) error {
	if tx == nil {
		return errors.New("cannot add a nil transaction")
	}

	s.poolMut.mut.Lock()
	defer s.poolMut.mut.Unlock()

	bytesStr := tx.SerializeString()
	s.poolMut.pool[buildKey(hash)] = bytesStr

	return nil
}

// GatherTransactions ...
func (s Service) GatherTransactions() (*[]statechain.Transaction, error) {
	s.poolMut.mut.Lock()
	defer s.poolMut.mut.Unlock()

	txs := make([]*statechain.Transaction, len(s.poolMut.pool), len(s.poolMut.mut))
	idx := 0
	for _, byteStr := range s.poolMut.pool {
		var tx statechain.Transaction
		if err := tx.DeserializeString(byteStr); err != nil {
			return nil, err
		}

		txs[idx] = &tx
		idx++
	}

	return txs, nil
}
