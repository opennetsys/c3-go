package ditto

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
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
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	fmt.Println("temp:", tmp)

	if err := untar(reader, tmp); err != nil {
		return err
	}

	return nil
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
