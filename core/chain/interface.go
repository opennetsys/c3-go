package chain

import (
	mainchain "github.com/c3systems/c3/core/chain/mainchain"
	statechain "github.com/c3systems/c3/core/chain/statechain"

	cid "gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
)

// Interface ...
type Interface interface {
	AddMainBlock(block *mainchain.Block) *cid.Cid
	Transactions() []*statechain.Transaction
	MainHead() (*mainchain.Block, error)
	StateHead(hash string) (*statechain.Block, error)
}
