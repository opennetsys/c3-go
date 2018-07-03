package p2p

import (
	"context"
	"errors"
	"time"

	cid "github.com/ipfs/go-cid"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	nonerouting "github.com/ipfs/go-ipfs-routing"
	bserv "github.com/ipfs/go-ipfs/blockservice"
	"github.com/ipfs/go-ipfs/exchange/bitswap"
	"github.com/ipfs/go-ipfs/exchange/bitswap/network"
	cbor "github.com/ipfs/go-ipld-cbor"
	//bserv "github.com/ipfs/go-ipfs/blockservice"
	//"github.com/ipfs/go-ds-flatfs"
	//"github.com/ipfs/go-ipfs/exchange/bitswap"
	//"github.com/ipfs/go-ipfs/exchange/bitswap/network"
	//bstore "github.com/ipfs/go-ipfs-blockstore"
	//nonerouting "github.com/ipfs/go-ipfs-routing"
	//multihash "github.com/multiformats/go-multihash"
	//cid "github.com/ipfs/go-cid"
	//cbor "github.com/ipfs/go-ipld-cbor"
	//host "github.com/libp2p/go-libp2p-host"
)

var service *Service

// Props ...
type Props struct {
	BlockStore bstore.Blockstore
	Host       *Host.host
}

// Service ...
type Service struct {
	props Props
	// Peers is a block store that will fetch blocks from other connected nodes
	peers bserv.BlockService
	// Local is a block store that will only fetch data locally
	local bstore.Blockstore
}

// New ...
func New(props Props) *Service {
	if service != nil {
		return service
	}

	// base backing datastore, currently just on disk, but can be swapped out
	// Register our types with the cbor encoder. This pregenerates serializers
	// for these types.
	cbor.RegisterCborType(mainchain.Block{})
	cbor.RegisterCborType(statechain.Block{})
	cbor.RegisterCborType(statechain.Transaction{})
	// TODO: need to store merkle tree tx's

	// wrap the datastore in a 'content addressed blocks' layer
	blocks := bstore.NewBlockstore(props.BlockStore)

	nr, _ := nonerouting.ConstructNilRouting(nil, nil, nil)
	bsnet := network.NewFromIpfsHost(props.Host, nr)

	bswap := bitswap.New(context.Background(), props.Host.ID(), bsnet, blocks, true)

	// Bitswap only fetches blocks from other nodes, to fetch blocks from
	// either the local cache, or a remote node, we can wrap it in a
	// 'blockservice'
	bservice := bserv.New(blocks, bswap)

	service = &Blockchain{
		peers: bservice,
		local: blocks,
	}

	return service
}

// GetCID ...
func GetCID(v interface{}) (*cid.Cid, error) {
	switch v.(type) {
	case *mainchain.Block:
		block, _ := v.(*mainchain.Block)
		if block == nil {
			return nil, errors.New("input cannot be nil")
		}

		nd, err := cbor.WrapObject(block, mh.SHA2_256, -1)
		if err != nil {
			return nil, err
		}

		return nd.Cid(), nil

	case *statechain.Block:
		block, _ := v.(*statechain.Block)
		if block == nil {
			return nil, errors.New("input cannot be nil")
		}

		nd, err := cbor.WrapObject(block, mh.SHA2_256, -1)
		if err != nil {
			return nil, err
		}

		return nd.Cid(), nil

	case *statechain.Transaction:
		tx, _ := v.(*statechain.Transaction)
		if tx == nil {
			return nil, errors.New("input cannot be nil")
		}

		nd, err := cbor.WrapObject(tx, mh.SHA2_256, -1)
		if err != nil {
			return nil, err
		}

		return nd.Cid(), nil

	default:
		return nil, errors.New("type must be one of pointer to mainchain block, statechain block or statechain tx")

	}
}

// LoadMainchainBlock ...
func LoadMainchainBlock(bs bserv.BlockService, c *cid.Cid) (*mainchain.Block, error) {
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

// LoadStateChainBlock ...
func LoadStateChainBlock(bs bserv.BlockService, c *cid.Cid) (*state.Block, error) {
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

// LoadStateChainTransaction
func LoadStateChainTransaction(bs bserv.BlockService, c *cid.Cid) (*state.Transaction, error) {
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
