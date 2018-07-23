package badgerstore

import (
	"path/filepath"
	"strings"

	"github.com/c3systems/c3-go/core/p2p/store"

	badger "github.com/dgraph-io/badger"
	ds "github.com/ipfs/go-datastore"
	badgerds "github.com/ipfs/go-ds-badger"
)

func New(path string, options *badger.Options) (ds.Batching, error) {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(store.UserHomeDir(), path[2:])
	}

	if err := store.CreateDirIfNotExist(path); err != nil {
		return nil, err
	}

	return badgerds.NewDatastore(path, options), nil
}
