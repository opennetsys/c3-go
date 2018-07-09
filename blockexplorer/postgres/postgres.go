package postgres

import (
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"

	cid "github.com/c3systems/go-cid"
)

// Props ...
type Props struct{}

// Service ...
type Service struct {
	props Props
}

// New ...
func New(props *Props) *Service {
	if props == nil {
		return &Service{}
	}

	return &Service{
		props: *props,
	}
}

// FetchCIDByMainBlockHash ...
func (s Service) FetchCIDByMainBlockHash(hexHash string) (*cid.Cid, error) {
	return nil, nil
}

// FetchCIDByStateBlockHash ...
func (s Service) FetchCIDByStateBlockHash(hexHash string) (*cid.Cid, error) {
	return nil, nil
}

// FetchCIDByTransactionHash ...
func (s Service) FetchCIDByTransactionHash(hexhHash string) (*cid.Cid, error) {
	return nil, nil
}

// FetchCIDByImageHashAndBlockNumber ...
func (s Service) FetchCIDByImageHashAndBlockNumber(imageHash, blockNumber string) (*cid.Cid, error) {
	return nil, nil
}

// FetchMainHashByStateBlockHash ...
func (s Service) FetchMainHashByStateBlockHash(hexHash string) (string, error) {
	return "", nil
}

// FetchStateBlockHashByTransactionHash ...
func (s Service) FetchStateBlockHashByTransactionHash(hexHash string) (string, error) {
	return "", nil
}

// FetchTransactionBySenderAddress ...
func (s Service) FetchTransactionsBySenderAddress(address string, skip, limit uint64) ([]*cid.Cid, error) {
	return nil, nil
}

// StoreMainBlockMeta ...
func (s *Service) StoreMainBlockMeta(block *mainchain.Block) error {
	return nil
}

// StoreStateBlockMeta ...
func (s *Service) StoreStateBlockMeta(block *statechain.Block) error {
	return nil
}

// StoreTransactionMeta ...
func (s *Service) StoreTransactionMeta(tx *statechain.Transaction) error {
	return nil
}

// StoreMainBlockCID ...
func (s *Service) StoreMainBlockCID(hexHash string, c *cid.Cid) error {
	return nil
}

// StoreStateBlockCID ...
func (s *Service) StoreStateBlockCID(hexHash string, c *cid.Cid) error {
	return nil
}

// StoreTransactionCID ...
func (s *Service) StoreTransactionCID(hexHash string, c *cid.Cid) error {
	return nil
}
