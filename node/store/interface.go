package store

import (
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/statechain"
)

// Interface ...
type Interface interface {
	HasTx(hash string) (bool, error)
	GetTx(hash string) (*statechain.Transaction, error)
	GetTxs(hashes []string) ([]*statechain.Transaction, error)
	RemoveTx(hash string) error
	RemoveTxs(hashes []string) error
	AddTx(tx *statechain.Transaction) error
	GatherPendingTransactions() ([]*statechain.Transaction, error)
	GetHeadBlock() (mainchain.Block, error)
	SetHeadBlock(block *mainchain.Block) error
	SetPendingMainchainBlock(block *mainchain.Block) error
	GetPendingMainchainBlocks() ([]*mainchain.Block, error)
	RemovePendingMainchainBlock(blockHash string) error
	RemovePendingMainchainBlocks(blockHashes []string) error
}
