package leveldbstore

import (
	"github.com/c3systems/c3-go/common/dirutil"

	ds "github.com/ipfs/go-datastore"
	leveldbds "github.com/ipfs/go-ds-leveldb"
)

// New ...
func New(path string, options *leveldbds.Options) (ds.Batching, error) {
	if err := dirutil.CreateDirIfNotExist(path); err != nil {
		return nil, err
	}

	return leveldbds.NewDatastore(path, options)
}
