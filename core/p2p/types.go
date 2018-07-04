package p2p

import (
	"sync"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	bstore "gx/ipfs/QmdpuJBPBZ6sLPj9BQpn3Rpi38BT2cF1QMiUfyzNWeySW4/go-ipfs-blockstore"

	bserv "github.com/ipfs/go-ipfs/blockservice"
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
