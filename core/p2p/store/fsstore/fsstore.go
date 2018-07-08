package fsstore

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	flatfs "github.com/ipfs/go-ds-flatfs"
)

// New ...
func New(path string) (*flatfs.Datastore, error) {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(userHomeDir(), path[2:])
	}

	var (
		shardFn *flatfs.ShardIdV1
		err     error
	)

	if err := createDirIfNotExist(path); err != nil {
		return nil, err
	}

	shardFn, err = flatfs.ReadShardFunc(path)
	if shardFn == nil || err != nil {
		shardFn = flatfs.Prefix(4)
		if err := flatfs.WriteShardFunc(path, shardFn); err != nil {
			return nil, err
		}
	}
	log.Printf("shard func: %v\nshard string: %s\n", shardFn.Func(), shardFn.String())

	return flatfs.CreateOrOpen(path, shardFn, true)
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

func createDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0757)
	}

	return nil
}
