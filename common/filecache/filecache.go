package filecache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// WIP

const cachedir = "/tmp"

var fmutex sync.RWMutex

// Set writes item to cache
func Set(key string, data interface{}, expire time.Duration) error {
	key = regexp.MustCompile("[^a-zA-Z0-9_-]").ReplaceAllLiteralString(key, "")
	file := fmt.Sprintf("filecache.%s.%v", key, strconv.FormatInt(time.Now().Add(expire).Unix(), 10))
	fpath := filepath.Join(cachedir, file)

	clean(key)

	serialized, err := serialize(data)
	if err != nil {
		return err
	}

	fmutex.Lock()
	defer fmutex.Unlock()
	fp, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer fp.Close()
	if _, err = fp.Write(serialized); err != nil {
		return err
	}

	return nil
}

// Get reads item from cache
func Get(key string, dst interface{}) (bool, error) {
	key = regexp.MustCompile("[^a-zA-Z0-9_-]").ReplaceAllLiteralString(key, "")
	pattern := filepath.Join(cachedir, fmt.Sprintf("filecache.%s.*", key))
	files, err := filepath.Glob(pattern)
	if err != nil {
		return false, err
	}
	if len(files) < 1 {
		return false, nil
	}

	if _, err = os.Stat(files[0]); err != nil {
		return false, err
	}

	fp, err := os.OpenFile(files[0], os.O_RDONLY, 0400)
	if err != nil {
		return false, err
	}
	defer fp.Close()

	var serialized []byte
	buf := make([]byte, 1024)
	for {
		var n int
		n, err = fp.Read(buf)
		serialized = append(serialized, buf[0:n]...)
		if err != nil || err == io.EOF {
			break
		}
	}

	if err = deserialize(serialized, dst); err != nil {
		return false, err
	}

	for _, file := range files {
		exptime, err := strconv.ParseInt(strings.Split(file, ".")[2], 10, 64)
		if err != nil {
			return false, err
		}

		if exptime < time.Now().Unix() {
			if _, err = os.Stat(file); err == nil {
				os.Remove(file)
			}
		}
	}

	return true, nil
}

// clean removes item from cache
func clean(key string) error {
	fmutex.Lock()
	defer fmutex.Unlock()
	pattern := filepath.Join(cachedir, fmt.Sprintf("filecache.%s.*", key))
	files, _ := filepath.Glob(pattern)
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
		}
	}

	return nil
}

// serialize encodes a value using binary
func serialize(src interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// deserialize decodes a value using binary
func deserialize(src []byte, dst interface{}) error {
	buf := bytes.NewReader(src)
	err := gob.NewDecoder(buf).Decode(dst)
	return err
}
