package p2p

import (
	"context"
	"errors"
	"time"

	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"

	bfmt "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	bserv "github.com/ipfs/go-ipfs/blockservice"
	cbor "github.com/ipfs/go-ipld-cbor"
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

	default:
		return nil, errors.New("type must be one of pointer to mainchain block, statechain block, statechain tx, or statechain diff")

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

	// note: will this work? We may need to pass the same bytes into the basicblock function
	//bytes, err := block.Serialize()
	//if err != nil {
	//return nil, err
	//}
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

	//bytes, err := block.Serialize()
	//if err != nil {
	//return nil, err
	//}
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

	//bytes, err := tx.Serialize()
	//if err != nil {
	//return nil, err
	//}
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

	//bytes, err := d.Serialize()
	//if err != nil {
	//return nil, err
	//}
	return GetCIDByHash(*d.Props().DiffHash)
}

// Fetch ...
func Fetch(bs bserv.BlockService, c *cid.Cid) (interface{}, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	var out interface{}
	if err := cbor.DecodeInto(data.RawData(), &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// FetchMainchainBlock ...
func FetchMainchainBlock(bs bserv.BlockService, c *cid.Cid) (*mainchain.Block, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	var out mainchain.Block
	if err := cbor.DecodeInto(data.RawData(), &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// FetchStateChainBlock ...
func FetchStateChainBlock(bs bserv.BlockService, c *cid.Cid) (*statechain.Block, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	var out statechain.Block
	if err := cbor.DecodeInto(data.RawData(), &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// FetchStateChainTransaction ...
func FetchStateChainTransaction(bs bserv.BlockService, c *cid.Cid) (*statechain.Transaction, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	var out statechain.Transaction
	if err := cbor.DecodeInto(data.RawData(), &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// FetchStateChainDiff ...
func FetchStateChainDiff(bs bserv.BlockService, c *cid.Cid) (*statechain.Diff, error) {
	if bs == nil || c == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := bs.GetBlock(ctx, c)
	if err != nil {
		return nil, err
	}

	var out statechain.Diff
	if err := cbor.DecodeInto(data.RawData(), &out); err != nil {
		return nil, err
	}

	return &out, nil
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
