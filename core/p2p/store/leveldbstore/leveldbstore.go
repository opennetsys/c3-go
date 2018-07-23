package leveldbstore

import (
	"path/filepath"
	"strings"

	"github.com/c3systems/c3-go/core/p2p/store"

	ds "github.com/ipfs/go-datastore"
	leveldbds "github.com/ipfs/go-ds-leveldb"
)

func New(path string, options *leveldbds.Options) (ds.Batching, error) {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(store.UserHomeDir(), path[2:])
	}

	if err := store.CreateDirIfNotExist(path); err != nil {
		return nil, err
	}

	return leveldbds.NewDatastore(path, options)
}
