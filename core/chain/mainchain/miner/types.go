package miner

import (
	"errors"
	"sync"

	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/merkle"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/p2p"
)

var (
	// ErrNilBlock ...
	ErrNilBlock = errors.New("block is nil")
	// ErrNoHash ...
	ErrNoHash = errors.New("no hash present")
	// ErrNilTx ...
	ErrNilTx = errors.New("transaction is nil")
	// ErrNoSig ...
	ErrNoSig = errors.New("no signature present")
	// ErrInvalidFromAddress ...
	ErrInvalidFromAddress = errors.New("from address is not valid")
	// ErrNilDiff ...
	ErrNilDiff = errors.New("diff is nil")
)

// Props is passed to the new function
type Props struct {
	IsValid             *bool // TODO: implement better fix than this isValid var
	PreviousBlock       *mainchain.Block
	Difficulty          uint64
	Channel             chan interface{}
	Async               bool // note: build state blocks asynchronously?
	P2P                 p2p.Interface
	PendingTransactions []*statechain.Transaction
}

// Service ...
type Service struct {
	props      Props
	minedBlock *MinedBlock
}

// MinedBlock ...
type MinedBlock struct {
	NextBlock     *mainchain.Block
	PreviousBlock *mainchain.Block

	// map keys are hashes
	mut                 sync.Mutex
	StatechainBlocksMap map[string]*statechain.Block
	TransactionsMap     map[string]*statechain.Transaction
	DiffsMap            map[string]*statechain.Diff
	MerkleTreesMap      map[string]*merkle.Tree
}
