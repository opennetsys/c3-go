package fsstore

import (
	"github.com/c3systems/c3-go/common/dirutil"

	flatfs "github.com/ipfs/go-ds-flatfs"
)

/*
 *
 * BEWARE: https://github.com/ipfs/go-ds-flatfs/issues/44
 *
 */

// New ...
func New(path string) (*flatfs.Datastore, error) {
	var (
		shardFn *flatfs.ShardIdV1
		err     error
	)

	if err := dirutil.CreateDirIfNotExist(path); err != nil {
		return nil, err
	}

	shardFn, err = flatfs.ReadShardFunc(path)
	if shardFn == nil || err != nil {
		shardFn = flatfs.Prefix(4)
		if err := flatfs.WriteShardFunc(path, shardFn); err != nil {
			return nil, err
		}
	}

	return flatfs.CreateOrOpen(path, shardFn, true)
}
