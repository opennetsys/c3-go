package fileutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// CreateTempFile ...
func CreateTempFile(filename string) (*os.File, error) {
	paths := strings.Split(filename, "/")

	// NOTE: does not like slashes for some reason, hence using underscore
	prefix := strings.Join(paths[:len(paths)-1], "_")
	filename = strings.Join(paths[len(paths)-1:len(paths)], "")

	tmpdir, err := ioutil.TempDir("/tmp", prefix)
	if err != nil {
		return nil, err
	}

	filepath := fmt.Sprintf("%s/%s", tmpdir, filename)

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// RemoveFiles ...
func RemoveFiles(fileNames *[]string) error {
	if fileNames == nil {
		return nil
	}

	for idx := range *fileNames {
		if err := os.Remove((*fileNames)[idx]); err != nil {
			return fmt.Errorf("err cleaning up file; %s; %v", (*fileNames)[idx], err)
		}
	}

	return nil
}
