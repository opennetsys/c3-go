package node

import (
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
)

type Interface interface {
	Ping() error
	BroadcastBlock(block *mainchain.Block) error
	BroadcastTransaction(tx *statechain.Transaction) error
	ListenTransactions()
	ListenBlocks()
	GatherTransactions() (statechain.TransactionMap, error)
	FetchHeadBlock() (*mainchain.Block, error)
}
