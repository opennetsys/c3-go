package chain

import (
	mainchain "github.com/c3systems/c3/core/chain/mainchain"
	statechain "github.com/c3systems/c3/core/chain/statechain"

	cid "github.com/ipfs/go-cid"
)

// Interface ...
type Interface interface {
	AddMainBlock(block *mainchain.Block) *cid.Cid
	PendingTransactions() []*statechain.Transaction
	MainHead() (*mainchain.Block, error)
	StateHead(hash string) (*statechain.Block, error)
}
