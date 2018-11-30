package fileutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// CreateTempFile ...
func CreateTempFile(filename string) (*os.File, error) {
	paths := strings.Split(filename, "/")

	// NOTE: does not like slashes for some reason, hence using underscore
	prefix := strings.Join(paths[:len(paths)-1], "_")
	filename = strings.Join(paths[len(paths)-1:len(paths)], "")

	// TODO: use os.TempDir()?
	tmpdir, err := ioutil.TempDir("/tmp", prefix)
	if err != nil {
		log.Errorf("err creating temp dir\n%v", err)
		return nil, err
	}

	filepath := fmt.Sprintf("%s/%s", tmpdir, filename)

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Errorf("err opening file: %s\n%v", filepath, err)
		return nil, err
	}

	return f, nil
}

// RemoveFiles ...
func RemoveFiles(fileNames []string) error {
	for idx := range fileNames {
		if err := os.Remove(fileNames[idx]); err != nil {
			log.Errorf("err removing file %s\n%v", fileNames[idx], err)
			return fmt.Errorf("err cleaning up file; %s; %v", fileNames[idx], err)
		}
	}

	return nil
}

// DirsFromFiles ...
func DirsFromFiles(fileNames []string) []string {
	tmp := make(map[string]struct{})
	for idx := range fileNames {
		dir, _ := filepath.Split(fileNames[idx])
		if dir != os.TempDir() {
			tmp[dir] = struct{}{}
		}
	}

	var ret []string
	for k := range tmp {
		ret = append(ret, k)
	}

	return ret
}

// RemoveDirs ...
func RemoveDirs(dirNames []string) error {
	var err error
	for idx := range dirNames {
		// note: do not want to remove the tmp dir
		if strings.ToLower(dirNames[idx]) == strings.ToLower(os.TempDir()) {
			continue
		}

		if err = os.RemoveAll(dirNames[idx]); err != nil {
			log.Errorf("err removing directory: %s\n%v", dirNames[idx], err)
			return err
		}
	}

	return nil
}
