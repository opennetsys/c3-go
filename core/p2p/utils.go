package p2p

import (
	"context"
	"errors"

	"github.com/c3systems/c3/config"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/merkle"
	"github.com/c3systems/c3/core/chain/statechain"

	bfmt "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	bserv "github.com/ipfs/go-ipfs/blockservice"
	mh "github.com/multiformats/go-multihash"
)

// GetCIDByHash ...
func GetCIDByHash(hash string) (*cid.Cid, error) {
	multiHash, err := mh.Sum([]byte(hash), mhCode, -1)
	if err != nil {
		return nil, err
	}

	return cid.NewCidV1(mhCode, multiHash), nil
}

// GetCID ...
func GetCID(v interface{}) (*cid.Cid, error) {
	if v == nil {
		return nil, errors.New("argument cannot be nil")
	}

	switch v.(type) {
	case *mainchain.Block:
		block, _ := v.(*mainchain.Block)
		return GetMainchainBlockCID(block)

	case *statechain.Block:
		block, _ := v.(*statechain.Block)
		return GetStatechainBlockCID(block)

	case *statechain.Transaction:
		tx, _ := v.(*statechain.Transaction)
		return GetStatechainTransactionCID(tx)

	case *statechain.Diff:
		d, _ := v.(*statechain.Diff)
		return GetStatechainDiffCID(d)

	case *merkle.Tree:
		tree, _ := v.(*merkle.Tree)
		return GetMerkleTreeCID(tree)

	default:
		return nil, errors.New("type must be one of pointer to mainchain block, statechain block, statechain tx, statechain diff, or merkle tree")

	}
}

// GetMainchainBlockCID ...
func GetMainchainBlockCID(block *mainchain.Block) (*cid.Cid, error) {
	if block == nil {
		return nil, errors.New("input cannot be nil")
	}
	if block.Props().BlockHash == nil {
		return nil, errors.New("hash cannot be nil")
	}

	return GetCIDByHash(*block.Props().BlockHash)
}

// GetStatechainBlockCID ...
func GetStatechainBlockCID(block *statechain.Block) (*cid.Cid, error) {
	if block == nil {
		return nil, errors.New("input cannot be nil")
	}
	if block.Props().BlockHash == nil {
		return nil, errors.New("hash cannot be nil")
	}

	return GetCIDByHash(*block.Props().BlockHash)
}

// GetStatechainTransactionCID ...
func GetStatechainTransactionCID(tx *statechain.Transaction) (*cid.Cid, error) {
	if tx == nil {
		return nil, errors.New("input cannot be nil")
	}
	if tx.Props().TxHash == nil {
		return nil, errors.New("hash cannot be nil")
	}

	return GetCIDByHash(*tx.Props().TxHash)
}

// GetStatechainDiffCID ...
func GetStatechainDiffCID(d *statechain.Diff) (*cid.Cid, error) {
	if d == nil {
		return nil, errors.New("input cannot be nil")
	}
	if d.Props().DiffHash == nil {
		return nil, errors.New("hash cannot be nil")
	}

	return GetCIDByHash(*d.Props().DiffHash)
}

// GetMerkleTreeCID ...
func GetMerkleTreeCID(tree *merkle.Tree) (*cid.Cid, error) {
	if tree == nil {
		return nil, errors.New("input cannot be nil")
	}
	if tree.Props().MerkleTreeRootHash == nil {
		return nil, errors.New("hash cannot be nil")
	}

	return GetCIDByHash(*tree.Props().MerkleTreeRootHash)
}

// note: generic fetch won't work bc we have to know what data type to deserialize into
// Fetch ...
//func Fetch(bs bserv.BlockService, c *cid.Cid) (interface{}, error) {
//if bs == nil || c == nil {
//return nil, errors.New("arguments cannot be nil")
//}

//ctx, cancel := context.WithTimeout(context.Background(), config.IPFSTimeout)
//defer cancel()

//data, err := bs.GetBlock(ctx, c)
//if err != nil {
//return nil, err
//}

//var out interface{}
//if err := cbor.DecodeInto(data.RawData(), &out); err != nil {
//return nil, err
//}

//return &out, nil
//}

// FetchMainchainBlock ...
func FetchMainchainBlock(bs bserv.BlockService, c *cid.Cid) (*mainchain.Block, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.IPFSTimeout)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	block := new(mainchain.Block)
	if err := block.Deserialize(data.RawData()); err != nil {
		return nil, err
	}

	return block, nil
}

// FetchStateChainBlock ...
func FetchStateChainBlock(bs bserv.BlockService, c *cid.Cid) (*statechain.Block, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.IPFSTimeout)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	block := new(statechain.Block)
	if err := block.Deserialize(data.RawData()); err != nil {
		return nil, err
	}

	return block, nil
}

// FetchStateChainTransaction ...
func FetchStateChainTransaction(bs bserv.BlockService, c *cid.Cid) (*statechain.Transaction, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.IPFSTimeout)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	tx := new(statechain.Transaction)
	if err := tx.Deserialize(data.RawData()); err != nil {
		return nil, err
	}

	return tx, nil
}

// FetchStateChainDiff ...
func FetchStateChainDiff(bs bserv.BlockService, c *cid.Cid) (*statechain.Diff, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.IPFSTimeout)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	d := new(statechain.Diff)
	if err := d.Deserialize(data.RawData()); err != nil {
		return nil, err
	}

	return d, nil
}

// FetchMerkleTree ...
func FetchMerkleTree(bs bserv.BlockService, c *cid.Cid) (*merkle.Tree, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.IPFSTimeout)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	tree := new(merkle.Tree)
	if err := tree.Deserialize(data.RawData()); err != nil {
		return nil, err
	}

	return tree, nil
}

// Put ...
func Put(bs bserv.BlockService, v interface{}) (*cid.Cid, error) {
	if bs == nil || v == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	switch v.(type) {
	case *mainchain.Block:
		block, _ := v.(*mainchain.Block)
		return PutMainchainBlock(bs, block)

	case *statechain.Block:
		block, _ := v.(*statechain.Block)
		return PutStatechainBlock(bs, block)

	case *statechain.Transaction:
		tx, _ := v.(*statechain.Transaction)
		return PutStatechainTransaction(bs, tx)

	case *statechain.Diff:
		d, _ := v.(*statechain.Diff)
		return PutStatechainDiff(bs, d)

	case *merkle.Tree:
		tree, _ := v.(*merkle.Tree)
		return PutMerkleTree(bs, tree)

	default:
		return nil, errors.New("type must be one of pointer to mainchain block, statechain block, statechain tx, or statechain diff")

	}
}

// PutMainchainBlock ...
func PutMainchainBlock(bs bserv.BlockService, block *mainchain.Block) (*cid.Cid, error) {
	if bs == nil || block == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	c, err := GetMainchainBlockCID(block)
	if err != nil {
		return nil, err
	}

	bytes, err := block.Serialize()
	if err != nil {
		return nil, err
	}

	basicIPFSBlock, err := bfmt.NewBlockWithCid(bytes, c)
	if err != nil {
		return nil, err
	}

	if err := bs.AddBlock(basicIPFSBlock); err != nil {
		return nil, err
	}

	return c, nil
}

// PutStatechainBlock ...
func PutStatechainBlock(bs bserv.BlockService, block *statechain.Block) (*cid.Cid, error) {
	if bs == nil || block == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	c, err := GetStatechainBlockCID(block)
	if err != nil {
		return nil, err
	}

	bytes, err := block.Serialize()
	if err != nil {
		return nil, err
	}

	basicIPFSBlock, err := bfmt.NewBlockWithCid(bytes, c)
	if err != nil {
		return nil, err
	}

	if err := bs.AddBlock(basicIPFSBlock); err != nil {
		return nil, err
	}

	return c, nil
}

// PutStatechainTransaction ...
func PutStatechainTransaction(bs bserv.BlockService, tx *statechain.Transaction) (*cid.Cid, error) {
	if bs == nil || tx == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	c, err := GetStatechainTransactionCID(tx)
	if err != nil {
		return nil, err
	}

	bytes, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	basicIPFSBlock, err := bfmt.NewBlockWithCid(bytes, c)
	if err != nil {
		return nil, err
	}

	if err := bs.AddBlock(basicIPFSBlock); err != nil {
		return nil, err
	}

	return c, nil
}

// PutStatechainDiff ...
func PutStatechainDiff(bs bserv.BlockService, d *statechain.Diff) (*cid.Cid, error) {
	if bs == nil || d == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	c, err := GetStatechainDiffCID(d)
	if err != nil {
		return nil, err
	}

	bytes, err := d.Serialize()
	if err != nil {
		return nil, err
	}

	basicIPFSBlock, err := bfmt.NewBlockWithCid(bytes, c)
	if err != nil {
		return nil, err
	}

	if err := bs.AddBlock(basicIPFSBlock); err != nil {
		return nil, err
	}

	return c, nil
}

// PutMerkleTree ...
func PutMerkleTree(bs bserv.BlockService, tree *merkle.Tree) (*cid.Cid, error) {
	if bs == nil || tree == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	c, err := GetMerkleTreeCID(tree)
	if err != nil {
		return nil, err
	}

	bytes, err := tree.Serialize()
	if err != nil {
		return nil, err
	}

	basicIPFSBlock, err := bfmt.NewBlockWithCid(bytes, c)
	if err != nil {
		return nil, err
	}

	if err := bs.AddBlock(basicIPFSBlock); err != nil {
		return nil, err
	}

	return c, nil
}
