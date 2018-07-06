package p2p

import (
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"

	// cid "gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	cid "github.com/ipfs/go-cid"
)

// Interface ...
type Interface interface {
	Props() Props
	Set(v interface{}) (*cid.Cid, error)
	SetMainchainBlock(block *mainchain.Block) (*cid.Cid, error)
	SetStatechainBlock(block *statechain.Block) (*cid.Cid, error)
	SetStateChainTransaction(tx *statechain.Transaction) (*cid.Cid, error)
	// TODO: how to do a generic get?
	//Get(c *cid.Cid) (interface{}, error) {
	GetMainchainBlock(c *cid.Cid) (*mainchain.Block, error)
	GetStatechainBlock(c *cid.Cid) (*statechain.Block, error)
	GetStatechainTransaction(c *cid.Cid) (*statechain.Transaction, error)
}
