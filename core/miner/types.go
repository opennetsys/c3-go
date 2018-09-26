package miner

import (
	"context"
	"errors"
	"sync"

	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/merkle"
	"github.com/c3systems/c3-go/core/chain/statechain"
	"github.com/c3systems/c3-go/core/p2p"
	"github.com/c3systems/c3-go/core/sandbox"
)

const (
	// StateFileName ...
	StateFileName string = "state.txt"
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
	// ErrInvalidTx ...
	ErrInvalidTx = errors.New("transaction is not valid")
	// ErrNilDiff ...
	ErrNilDiff = errors.New("diff is nil")
)

// Props is passed to the new function
type Props struct {
	Context             context.Context
	PreviousBlock       *mainchain.Block
	Difficulty          uint64
	Channel             chan interface{}
	Async               bool // note: build state blocks asynchronously?
	EncodedMinerAddress string
	P2P                 p2p.Interface
	Sandbox             sandbox.Interface
	PendingTransactions []*statechain.Transaction
	RemoveTx            func(hash string) error
	Simulated           bool
}

// Service ...
type Service struct {
	props      Props
	minedBlock *MinedBlock
}

// MinedBlock ...
type MinedBlock struct {
	NextBlock     *mainchain.Block `json:"nextBlock"`
	PreviousBlock *mainchain.Block `json:"previousBlock"`

	// map keys are hashes
	// TODO: add previous statechain blocks map
	mut                 sync.Mutex
	StatechainBlocksMap map[string]*statechain.Block       `json:"statechainBlocksMap"`
	TransactionsMap     map[string]*statechain.Transaction `json:"transactionsMap"`
	DiffsMap            map[string]*statechain.Diff        `json:"diffsMap"`
	MerkleTreesMap      map[string]*merkle.Tree            `json:"merkleTreesMap"`
}
