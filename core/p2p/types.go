package p2p

import (
	"sync"

	// host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	host "gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
	// bstore "gx/ipfs/QmPpdpS9fknTBM3qHDcpayU6nYPZQeVjia2fbNrD8YWDe6/go-ipfs-blockstore"
	// bstore "gx/ipfs/QmTVDM4LCSUMFNQzbDLL9zQwp8usE6QHymFdh3h8vL9v6b/go-ipfs-blockstore"
	// bstore "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/blockstore"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	// bstore "gx/ipfs/QmdpuJBPBZ6sLPj9BQpn3Rpi38BT2cF1QMiUfyzNWeySW4/go-ipfs-blockstore"
	bstore "gx/ipfs/QmaG4DZ4JaqEfvPWt5nPPgoTzhc1tr1T3f4Nu9Jpdm8ymY/go-ipfs-blockstore"
	// mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"

	// bstore "github.com/ipfs/go-ipfs-blockstore"
	// host "github.com/libp2p/go-libp2p-host"
	// mh "github.com/multiformats/go-multihash"

	// bserv "github.com/ipfs/go-ipfs/blockservice"
	bserv "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/blockservice"
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
