package p2p

import (
	"context"
	"errors"
	"time"

	cid "github.com/ipfs/go-cid"
	bserv "github.com/ipfs/go-ipfs/blockservice"
	cbor "github.com/ipfs/go-ipld-cbor"
)

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

	default:
		return nil, errors.New("type must be one of pointer to mainchain block, statechain block or statechain tx")

	}
}

// GetMainchainBlockCID ...
func GetMainchainBlockCID(block *mainchain.Block) (*cid.Cid, error) {
	if block == nil {
		return nil, errors.New("input cannot be nil")
	}

	nd, err := cbor.WrapObject(block, hashingAlgo, -1)
	if err != nil {
		return nil, err
	}

	return nd.Cid(), nil
}

// GetStatechainBlockCID ...
func GetStatechainBlockCID(block *statechain.Block) (*cid.Cid, error) {
	if block == nil {
		return nil, errors.New("input cannot be nil")
	}

	nd, err := cbor.WrapObject(block, hashingAlgo, -1)
	if err != nil {
		return nil, err
	}

	return nd.Cid(), nil
}

// GetStatechainTransactionCID ...
func GetStatechainTransactionCID(tx *statechain.Transaction) (*cid.Cid, error) {
	if tx == nil {
		return nil, errors.New("input cannot be nil")
	}

	nd, err := cbor.WrapObject(tx, hashingAlgo, -1)
	if err != nil {
		return nil, err
	}

	return nd.Cid(), nil
}

// Fetch ...
// TODO: how to do a generic fetch?
//func Fetch(bs bserv.BlockService, c *cid.Cid) (v interface{}, error) {
//return nil, nil
//}

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
func FetchStateChainBlock(bs bserv.BlockService, c *cid.Cid) (*state.Block, error) {
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

// FetchStateChainTransaction
func FetchStateChainTransaction(bs bserv.BlockService, c *cid.Cid) (*state.Transaction, error) {
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

// Put ...
func Put(bs bserv.BlockService, v interface{}) (*cid.Cid, error) {
	if bs == nil || v == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	switch v.(type) {
	case *mainchain.Block:
		block, _ := v.(*mainchain.Block)
		return PutMainchainBlock(block)

	case *statechain.Block:
		block, _ := v.(*statechain.Block)
		return PutStatechainBlock(block)

	case *statechain.Transaction:
		tx, _ := v.(*statechain.Transaction)
		return PuttStatechainTransaction(tx)

	default:
		return nil, errors.New("type must be one of pointer to mainchain block, statechain block or statechain tx")

	}
}

// PutMainchainBlock ...
func PutMainchainBlock(bs bserv.BlockService, block *mainchain.Block) (*cid.Cid, error) {
	if bs == nil || block == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	nd, err := cbor.WrapObject(block, hashingAlgo, 32)
	if err != nil {
		return nil, err
	}

	if err := bs.AddBlock(nd); err != nil {
		return nil, err
	}

	return nd.Cid(), nil
}

// PutStatechainBlock ...
func PutStatechainBlock(bs bserv.BlockService, block *statechain.Block) (*cid.Cid, error) {
	if bs == nil || block == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	nd, err := cbor.WrapObject(block, hashingAlgo, 32)
	if err != nil {
		return nil, err
	}

	if err := bs.AddBlock(nd); err != nil {
		return nil, err
	}

	return nd.Cid(), nil
}

// PutStatechainTransaction ...
func PutStatechainTransaction(bs bserv.BlockService, tx *statechain.Transaction) (*cid.Cid, error) {
	if bs == nil || tx == nil {
		return nil, errors.New("arguments cannot be nil")
	}

	nd, err := cbor.WrapObject(tx, hashingAlgo, 32)
	if err != nil {
		return nil, err
	}

	if err := bs.AddBlock(nd); err != nil {
		return nil, err
	}

	return nd.Cid(), nil
}
