package memstore

import datastore "github.com/ipfs/go-datastore"

func New() *datastore.Datastore {
	return datastore.NewMapDatastore()
}
