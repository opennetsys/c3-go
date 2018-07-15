package merkle

import (
	"crypto/sha256"
	"errors"

	"github.com/c3systems/merkletree"
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

// testContent implements the Content interface provided by merkletree and represents the content stored in the tree.
type testContent struct {
	x string
}

// CalculateHashBytes hashes the values of a TestContent
func (t testContent) CalculateHashBytes() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t testContent) Equals(other merkletree.Content) (bool, error) {
	return t.x == other.(testContent).x, nil
}
