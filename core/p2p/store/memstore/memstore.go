package memstore

import (
	datastore "github.com/ipfs/go-datastore"
	ds "github.com/ipfs/go-datastore"
)

// New ...
func New() ds.Batching {
	return datastore.NewMapDatastore()
}
