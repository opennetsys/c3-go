package memstore

import datastore "github.com/ipfs/go-datastore"

// New ...
func New() datastore.Datastore {
	return datastore.NewMapDatastore()
}
