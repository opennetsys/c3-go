package p2p

import (
	"context"
	"errors"

	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"

	cid "github.com/ipfs/go-cid"
	nonerouting "github.com/ipfs/go-ipfs-routing/none"
	bserv "github.com/ipfs/go-ipfs/blockservice"
	"github.com/ipfs/go-ipfs/exchange/bitswap"
	"github.com/ipfs/go-ipfs/exchange/bitswap/network"
)

// New ...
func New(props *Props) (*Service, error) {
	var err error

	once.Do(func() {
		if props == nil {
			err = errors.New("props cannot be nil")
			return
		}
		if props.Host == nil || props.BlockStore == nil {
			err = errors.New("host and blockstore are required")
			return
		}

		// Register our types with the cbor encoder. This pregenerates serializers
		// for these types.
		// cbor.RegisterCborType(mainchain.Block{})
		// cbor.RegisterCborType(statechain.Block{})
		// cbor.RegisterCborType(statechain.Transaction{})

		// TODO: research if this is what we want...
		nr, err1 := nonerouting.ConstructNilRouting(nil, nil, nil, nil)
		if err1 != nil {
			err = err1
			return
		}

		bsnet := network.NewFromIpfsHost(props.Host, nr)
		bswap := bitswap.New(context.Background(), bsnet, props.BlockStore)

		// Bitswap only fetches blocks from other nodes, to fetch blocks from
		// either the local cache, or a remote node, we can wrap it in a
		// 'blockservice'
		bservice := bserv.New(props.BlockStore, bswap)

		service = &Service{
			props:        *props,
			peersOrLocal: bservice,
			local:        props.BlockStore,
		}
	})

	return service, err
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
	return FetchStateChainBlock(s.peersOrLocal, c)
}

// GetStatechainTransaction ...
func (s Service) GetStatechainTransaction(c *cid.Cid) (*statechain.Transaction, error) {
	return FetchStateChainTransaction(s.peersOrLocal, c)
}
