package leveldbstore

import (
	leveldbds "github.com/ipfs/go-ds-leveldb"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	opt "github.com/syndtr/goleveldb/leveldb/opt"
	//bstore "gx/ipfs/QmTVDM4LCSUMFNQzbDLL9zQwp8usE6QHymFdh3h8vL9v6b/go-ipfs-blockstore"
)

type Options opt.Options

func New(path string, options *Options) (bstore.Blockstore, error) {
	return leveldbds.NewDatastore(path, options)
}
