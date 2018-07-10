package miner

import (
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
)

// Interface ...
// TODO: finish
type Interface interface {
	Props() Props
	Start() error
	VerifyMainchainBlock(block *mainchain.Block) (bool, error)
	VerifyStatechainBlock(block *statechain.Block) (bool, error)
}
