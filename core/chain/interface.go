package chain

import (
	mainchain "github.com/c3systems/c3-go/core/chain/mainchain"
	statechain "github.com/c3systems/c3-go/core/chain/statechain"

	cid "github.com/ipfs/go-cid"
)

// Interface ...
type Interface interface {
	Props() Props
	AddMainBlock(block *mainchain.Block) *cid.Cid
	PendingTransactions() []*statechain.Transaction
	MainHead() (*mainchain.Block, error)
	StateHead(hash string) (*statechain.Block, error)
}
