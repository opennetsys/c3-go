package badgerstore

import (
	badger "github.com/dgraph-io/badger"
	badgerds "github.com/ipfs/go-ds-badger"
	bstore "github.com/ipfs/go-ipfs-blockstore"
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
