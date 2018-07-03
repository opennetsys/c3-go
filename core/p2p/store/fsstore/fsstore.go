package fsstore

import (
	flatfs "github.com/ipfs/go-ds-flatfs"
)

func New(path string, fun *flatfs.ShardIdV1, sync bool) (*flatfs.Datastore, error) {
	return fs.CreateOrOpen(path, fun, sync)
}
