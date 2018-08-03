package fsstore

import (
	"path/filepath"
	"strings"

	"github.com/c3systems/c3-go/core/p2p/store"

	flatfs "github.com/ipfs/go-ds-flatfs"
)

/*
 *
 * BEWARE: https://github.com/ipfs/go-ds-flatfs/issues/44
 *
 */

// New ...
func New(path string) (*flatfs.Datastore, error) {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(store.UserHomeDir(), path[2:])
	}

	var (
		shardFn *flatfs.ShardIdV1
		err     error
	)

	if err := store.CreateDirIfNotExist(path); err != nil {
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
