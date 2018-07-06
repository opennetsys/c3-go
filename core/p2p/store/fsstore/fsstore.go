package fsstore

import (
	"fmt"
	"log"
	"os/user"

	flatfs "github.com/ipfs/go-ds-flatfs"
)

// New ...
func New(path string) (*flatfs.Datastore, error) {
	dir := path
	if dir == "~/c3-data/" {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}

		dir = fmt.Sprintf("%s/c3-data", usr.HomeDir)
	}

	var (
		shardFn *flatfs.ShardIdV1
		err     error
	)

	shardFn, err = flatfs.ReadShardFunc(dir)
	if shardFn == nil || err != nil {
		log.Printf("err reading shardfn\n%v", err)
		shardFn = flatfs.Prefix(4)
		if err := flatfs.WriteShardFunc(dir, shardFn); err != nil {
			return nil, err
		}
	}
	log.Printf("shard func: %v\nshard string: %s\n", shardFn.Func(), shardFn.String())
	return flatfs.CreateOrOpen(dir, shardFn, true)
}
