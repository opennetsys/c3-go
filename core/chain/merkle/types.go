package merkle

import (
	"errors"
)

var (
	// ErrUnknownKind ...
	ErrUnknownKind = errors.New("unknown kind")
	// ErrNilChainObjects ...
	ErrNilChainObjects = errors.New("nil chain objects")
	// ErrNilChainObject ...
	ErrNilChainObject = errors.New("nil chain object")
	// ErrInconsistentKinds ...
	ErrInconsistentKinds = errors.New("inconsistent kinds")
	// ErrNilMerkleTree ...
	ErrNilMerkleTree = errors.New("merkle tree is nil")
	// ErrNilMerkleTreeRootHash ...
	ErrNilMerkleTreeRootHash = errors.New("merkle tree hash is nil")
	// ErrNilProps ...
	ErrNilProps  = errors.New("props are nil")
	allowedKinds = []string{
		StatechainBlocksKindStr,
		MainchainBlocksKindStr,
		TransactionsKindStr,
		DiffsKindStr,
		MerkleTreesKindStr,
	}
)

const (
	// StatechainBlocksKindStr ...
	StatechainBlocksKindStr = "statechainBlocks"
	// MainchainBlocksKindStr ...
	MainchainBlocksKindStr = "mainchainBlocks"
	// TransactionsKindStr ...
	TransactionsKindStr = "transactions"
	// DiffsKindStr ...
	DiffsKindStr = "diffs"
	// MerkleTreesKindStr ...
	MerkleTreesKindStr = "merkleTrees"
)

// Tree ...
type Tree struct {
	props TreeProps
}

// TreeProps ...
type TreeProps struct {
	MerkleTreeRootHash *string
	Kind               string
	Hashes             []string
}
