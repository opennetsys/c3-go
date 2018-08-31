package badgerstore

import (
	"github.com/c3systems/c3-go/common/dirutil"

	badger "github.com/dgraph-io/badger"
	ds "github.com/ipfs/go-datastore"
	badgerds "github.com/ipfs/go-ds-badger"
)

// New ...
func New(path string, options *badger.Options) (ds.Batching, error) {
	if err := dirutil.CreateDirIfNotExist(path); err != nil {
		return nil, err
	}

	return badgerds.NewDatastore(path, options), nil
}
