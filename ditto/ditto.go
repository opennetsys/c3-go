package ditto

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Ditto ...
type Ditto struct {
}

// Config ...
type Config struct {
}

// New ...
func New(config *Config) *Ditto {
	return &Ditto{}
}

// UploadImage uploads Docker image to IPFS
func (s Ditto) UploadImage(reader io.Reader) error {
	tmp := mktemp()
	fmt.Println("temp:", tmp)

	if err := untar(reader, tmp); err != nil {
		return err
	}

	if err := process(tmp); err != nil {
		return err
	}

	return nil
}

func mktemp() string {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	return tmp
}

func process(tmp string) error {
	root := mktemp()
	workdir := root
	fmt.Println("preparing image in:", workdir)
	reposJSON, err := readJSON(tmp + "/" + "repositories")
	if err != nil {
		return err
	}
	if len(reposJSON) != 1 {
		return errors.New("only one repository expected in input file")
	}
	var name string
	for imageName, tags := range reposJSON {
		fmt.Println(imageName, tags)
		if len(tags) != 1 {
			return fmt.Errorf("only one tag expected for %s", imageName)
		}
		for tag, hash := range tags {
			name = normalizeImageName(imageName)
			fmt.Printf("processing image:%s tag:%s hash:256:%s", name, tag, hash)
		}
	}

	return nil
}

func readJSON(filepath string) (map[string]map[string]string, error) {
	body, _ := ioutil.ReadFile(filepath)
	var data map[string]map[string]string
	err := json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func untar(reader io.Reader, dst string) error {
	tr := tar.NewReader(reader)

	for {
		header, err := tr.Next()
		switch {
		// no more files
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		// create directory if doesn't exit
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// create file
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			// copy contents to file
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

func normalizeImageName(name string) string {
	// TODO
	return name
}
