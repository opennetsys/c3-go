package badgerstore

import (
	badger "github.com/dgraph-io/badger"
	badgerds "github.com/ipfs/go-ds-badger"
	bstore "gx/ipfs/QmTVDM4LCSUMFNQzbDLL9zQwp8usE6QHymFdh3h8vL9v6b/go-ipfs-blockstore"
)

// Options are params for creating DB object.
//
// note: DO NOT set the Dir and/or ValuePath fields of opt, they will be set for you.
type Options struct {
	badger.Options
}

func New(path string, options *Options) (bstore.Blockstore, error) {
	return badgerds.NewDatastore(path, &badgerds.Options{
		options,
	}), nil
}
