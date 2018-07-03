package p2p

import (
	"context"
	"errors"

	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
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
	//mh "github.com/multiformats/go-multihash"
	//cid "github.com/ipfs/go-cid"
	//cbor "github.com/ipfs/go-ipld-cbor"
	//host "github.com/libp2p/go-libp2p-host"
)

// New ...
func New(props *Props) (*Service, error) {
	once.Do(func() {
		if props == nil {
			return nil, errors.New("props cannot be nil")
		}
		if props.Host == nil || props.BlockStore == nil {
			return nil, errors.New("host and blockstore are required")
		}

		// Register our types with the cbor encoder. This pregenerates serializers
		// for these types.
		cbor.RegisterCborType(mainchain.Block{})
		cbor.RegisterCborType(statechain.Block{})
		cbor.RegisterCborType(statechain.Transaction{})
		// TODO: need to store merkle tree tx's

		// wrap the datastore in a 'content addressed blocks' layer
		// TODO: implement metrics? https://github.com/ipfs/go-ds-measure
		blocks := bstore.NewBlockstore(props.BlockStore)

		// TODO: research if this is what we want...
		nr, err := nonerouting.ConstructNilRouting(nil, nil, nil)
		if err != nil {
			return nil, err
		}

		bsnet := network.NewFromIpfsHost(props.Host, nr)
		bswap := bitswap.New(context.Background(), props.Host.ID(), bsnet, blocks, true)

		// Bitswap only fetches blocks from other nodes, to fetch blocks from
		// either the local cache, or a remote node, we can wrap it in a
		// 'blockservice'
		bservice := bserv.New(blocks, bswap)

		service = &Service{
			props:        *props,
			peersOrLocal: bservice,
			local:        blocks,
		}
	})

	return service, nil
}

// Props ...
func (s Service) Props() Props {
	return s.props
}

// Set ...
func (s Service) Set(v interface{}) (*cid.Cid, error) {
	return Put(s.peersOrLocal, v)
}

// SetMainchainBlock ...
// note: this function does not do any validation!
func (s Service) SetMainchainBlock(block *mainchain.Block) (*cid.Cid, error) {
	return PutMainchainBlock(s.peersOrLocal, block)
}

// SetStatechainBlock ...
func (s Service) SetStatechainBlock(block *statechain.Block) (*cid.Cid, error) {
	return PutStatechainBlock(s.peersOrLocal, block)
}

// SetStatechainTransaction ...
func (s Service) SetStatechainTransaction(tx *statechain.Transaction) (*cid.Cid, error) {
	return PutStatechainTransaction(s.peersOrLocal, tx)
}

// Get ...
// TODO: how to do a generic get?
//func (s Service) Get(c *cid.Cid) (interface{}, error) {
//return Fetch(s.peers, c)
//}

// GetMainchainBlock ...
func (s Service) GetMainchainBlock(c *cid.Cid) (*mainchain.Block, error) {
	return FetchMainchainBlock(s.peersOrLocal, c)
}

// GetStatechainBlock ...
func (s Service) GetStatechainBlock(c *cid.Cid) (*statechain.Block, error) {
	return GetStatechainBlock(s.peersOrLocal, c)
}

// GetStatechainTransaction ...
func (s Service) GetStatechainTransaction(c *cid.Cid) (*statechain.Transaction, error) {
	return GetStatechainTransaction(s.peersOrLocal, c)
}
