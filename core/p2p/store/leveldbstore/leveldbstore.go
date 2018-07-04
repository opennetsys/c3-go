package leveldbstore

import (
	leveldbds "github.com/ipfs/go-ds-leveldb"
	opt "github.com/syndtr/goleveldb/leveldb/opt"
	bstore "gx/ipfs/QmdpuJBPBZ6sLPj9BQpn3Rpi38BT2cF1QMiUfyzNWeySW4/go-ipfs-blockstore"
)

type Options opt.Options

func New(path string, options *Options) (bstore.Blockstore, error) {
	return leveldbds.NewDatastore(path, options)
}
