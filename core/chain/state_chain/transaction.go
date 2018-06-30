package statechain

import "sync"

// Transaction ...
type Transaction struct{}

// TransactionsMutex ...
type TransactionsMutex struct {
	Mut sync.Mutex
	Txs []*Transaction
}
