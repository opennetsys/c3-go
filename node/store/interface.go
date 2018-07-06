package store

import "github.com/c3systems/c3/core/chain/statechain"

// Interface ...
type Interface interface {
	HasTx(hash string) (bool, error)
	GetTx(hash string) (*statechain.Transaction, error)
	GetTxs(hashes []string) ([]*statechain.Transaction, error)
	RemoveTx(hash string) error
	RemoveTxs(hashes []string) error
	AddTx(tx *statechain.Transaction) error
	GatherTransactions() ([]*statechain.Transaction, error)
}
