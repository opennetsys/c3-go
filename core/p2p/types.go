package p2p

import (
	"sync"

	bstore "github.com/ipfs/go-ipfs-blockstore"
	bserv "github.com/ipfs/go-ipfs/blockservice"
	host "github.com/libp2p/go-libp2p-host"
	mh "github.com/multiformats/go-multihash"
)

const hashingAlgo = mh.SHA2_256

var (
	service *Service
	once    sync.Once
)

// Props ...
type Props struct {
	BlockStore bstore.Blockstore // note: https://github.com/ipfs/go-ipfs/blob/master/docs/datastores.md
	Host       host.Host
}

// Service ...
type Service struct {
	props Props
	// Peers is a block store that will fetch blocks from other connected nodes
	peersOrLocal bserv.BlockService
	// Local is a block store that will only fetch data locally
	local bstore.Blockstore
}
