package memstore

import datastore "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore"

// New ...
func New() datastore.Datastore {
	return datastore.NewMapDatastore()
}
