package fsstore

import (
	flatfs "github.com/ipfs/go-ds-flatfs"
)

// New ...
func New(path string, fun *flatfs.ShardIdV1, sync bool) (*flatfs.Datastore, error) {
	return flatfs.CreateOrOpen(path, fun, sync)
}
