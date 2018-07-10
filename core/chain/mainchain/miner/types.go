package miner

import (
	"errors"
	"sync"

	"github.com/c3systems/c3/core/chain/mainchain"
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

type stateBlocksMutex struct {
	mut    sync.Mutex
	blocks []*statechain.Block
}

// Props is passed to the new function
type Props struct {
	Difficulty         uint64
	Channel            chan interface{}
	Async              bool // note: build state blocks asynchronously?
	P2P                p2p.Interface
	GatherTransactions func() ([]*statechain.Transaction, error)
	FetchHeadBlock     func() (*mainchain.Block, error)
}

// Service ...
type Service struct {
	props Props
}
