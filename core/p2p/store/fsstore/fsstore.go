package fsstore

import (
	"fmt"
	"log"
	"os"
	"runtime"

	flatfs "github.com/ipfs/go-ds-flatfs"
)

// New ...
func New(path string) (*flatfs.Datastore, error) {
	dir := path
	if dir == "~/c3-data/" {
		dir = fmt.Sprintf("%s/c3-data", userHomeDir())
	}

	var (
		shardFn *flatfs.ShardIdV1
		err     error
	)

	if err := createDirIfNotExist(dir); err != nil {
		return nil, err
	}

	shardFn, err = flatfs.ReadShardFunc(dir)
	if shardFn == nil || err != nil {
		shardFn = flatfs.Prefix(4)
		if err := flatfs.WriteShardFunc(dir, shardFn); err != nil {
			return nil, err
		}
	}
	log.Printf("shard func: %v\nshard string: %s\n", shardFn.Func(), shardFn.String())
	return flatfs.CreateOrOpen(dir, shardFn, true)
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0757)
	}

	return nil
}
