package leveldbstore

import (
	leveldbds "github.com/ipfs/go-ds-leveldb"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	opt "github.com/syndtr/goleveldb/leveldb/opt"
)

type Options opt.Options

func New(path string, options *Options) (bstore.Blockstore, error) {
	return leveldbds.NewDatastore(path, options)
}
