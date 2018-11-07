package p2p

import (
	"sync"

	bstore "github.com/ipfs/go-ipfs-blockstore"
	bserv "github.com/ipfs/go-blockservice"
	host "github.com/libp2p/go-libp2p-host"
	routing "github.com/libp2p/go-libp2p-routing"
	mh "github.com/multiformats/go-multihash"
)

const (
	mhCode = mh.SHA2_256
)

var (
	service *Service
	once    sync.Once
)

// Props ...
type Props struct {
	// NOTE: https://github.com/ipfs/go-ipfs/blob/master/docs/datastores.md
	BlockStore bstore.Blockstore
	Host       host.Host
	Router     routing.ContentRouting
}

// Service ...
type Service struct {
	props Props
	// Peers is a block store that will fetch blocks from other connected nodes
	peersOrLocal bserv.BlockService
	// Local is a block store that will only fetch data locally
	local bstore.Blockstore
}
