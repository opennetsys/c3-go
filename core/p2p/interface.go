package p2p

import (
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/merkle"
	"github.com/c3systems/c3-go/core/chain/statechain"

	cid "github.com/ipfs/go-cid"
)

// Interface ...
type Interface interface {
	Props() Props
	Set(v interface{}) (*cid.Cid, error)
	SetMainchainBlock(block *mainchain.Block) (*cid.Cid, error)
	SetStatechainBlock(block *statechain.Block) (*cid.Cid, error)
	SetStatechainTransaction(tx *statechain.Transaction) (*cid.Cid, error)
	SetStatechainDiff(d *statechain.Diff) (*cid.Cid, error)
	SetMerkleTree(tree *merkle.Tree) (*cid.Cid, error)
	SetBytes(b []byte) (*cid.Cid, error)
	SetLatestBlock(block *mainchain.Block) (*cid.Cid, error)
	//SaveLocal(v interface{}) (*cid.Cid, error)
	//SaveLocalMainchainBlock(block *mainchain.Block) (*cid.Cid, error)
	//SaveLocalStatechainBlock(block *statechain.Block) (*cid.Cid, error)
	//SaveLocalStatechainTransaction(tx *statechain.Transaction) (*cid.Cid, error)
	//SaveLocalStatechainDiff(d *statechain.Diff) (*cid.Cid, error)
	//Get(c *cid.Cid) (interface{}, error)
	GetMainchainBlock(c *cid.Cid) (*mainchain.Block, error)
	GetStatechainBlock(c *cid.Cid) (*statechain.Block, error)
	GetStatechainTransaction(c *cid.Cid) (*statechain.Transaction, error)
	GetStatechainDiff(c *cid.Cid) (*statechain.Diff, error)
	GetMerkleTree(c *cid.Cid) (*merkle.Tree, error)
	GetBytes(c *cid.Cid) ([]byte, error)
	GetLatestBlock() (*mainchain.Block, error)
	FetchMostRecentStateBlock(imageHash string, block *mainchain.Block) (*statechain.Block, error)
}
