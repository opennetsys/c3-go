package disk

import (
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/statechain"
)

// TODO: everything...

// Props ...
type Props struct {
}

// Service ...
type Service struct {
	props     Props
	headBlock mainchain.Block
}

// New ...
func New(props *Props) (*Service, error) {
	return nil, nil
}

// Props ...
func (s Service) Props() Props {
	return s.props
}

// HasTx ...
func (s Service) HasTx(hash string) (bool, error) {
	return false, nil
}

// GetTx ...
func (s Service) GetTx(hash string) (*statechain.Transaction, error) {
	return nil, nil
}

// GetTxs ...
func (s Service) GetTxs(hashes []string) ([]*statechain.Transaction, error) {
	return nil, nil
}

// RemoveTx ...
func (s Service) RemoveTx(hash string) error {
	return nil
}

// RemoveTxs ...
func (s Service) RemoveTxs(hashes []string) error {
	return nil
}

// AddTx ...
func (s Service) AddTx(tx *statechain.Transaction) error {
	return nil
}

// GatherTransactions ...
func (s Service) GatherTransactions() (*[]statechain.Transaction, error) {
	return nil, nil
}

// GetHeadBlock ...
func (s Service) GetHeadBlock() (mainchain.Block, error) {
	return s.headBlock, nil
}

// SetHeadBlock ...
func (s Service) SetHeadBlock(block mainchain.Block) error {
	s.headBlock = block
	return nil
}
