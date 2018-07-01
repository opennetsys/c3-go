package ditto

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
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

	if err := ipfsPrep(tmp); err != nil {
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

func ipfsPrep(tmp string) error {
	root := mktemp()
	workdir := root
	fmt.Println("preparing image in:", workdir)
	reposJSON, err := readJSON(tmp + "/repositories")
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

	workdir = workdir + "/" + name
	mkdir(workdir)
	mkdir(workdir + "/manifests")
	mkdir(workdir + "/blobs")
	manifestJSON, err := readJSONArray(tmp + "/manifest.json")
	if err != nil {
		return err
	}

	if len(manifestJSON) == 0 {
		return errors.New("expected manifest to contain data")
	}

	manifest := manifestJSON[0]
	configFile, ok := manifest["Config"].(string)
	if !ok {
		return errors.New("image archive must be produced by docker > 1.10")
	}

	configDest := workdir + "/blobs/sha256:" + string(configFile[:len(configFile)-5])
	fmt.Println("\nDIST", configDest)
	mkdir(configDest)
	copyfile(tmp+"/"+configFile, configDest+"/"+configFile)

	mf := makeV2Manifest(manifest, configFile, configDest, tmp, workdir)
	spew.Dump(mf)

	//writeJSON(mf, workdir, "manifests", "latest-v2")

	//proc = subprocess.Popen(['ipfs', 'add', '-r', '-q', root], stdout=subprocess.PIPE)

	return nil
}

func writeJSON() {

}

// produce v2 manifest of type/application/vnd.docker.distribution.manifest.v2+json
func makeV2Manifest(manifest map[string]interface{}, configFile, configDest, tmp, workdir string) map[string]interface{} {
	v2manifest := prepareV2Manifest(manifest, tmp, workdir+"/blobs")
	config := make(map[string]interface{})
	config["digest"] = "sha256:" + string(configFile[:len(configFile)-5])
	config["size"] = fileSize(configDest)
	conf, ok := v2manifest["config"].(map[string]interface{})
	if !ok {
	}
	v2manifest["config"] = mergemap(conf, config)
	return v2manifest
}

func mergemap(a, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func prepareV2Manifest(mf map[string]interface{}, tmp, blobDir string) map[string]interface{} {
	res := make(map[string]interface{})
	res["schemaVersion"] = "application/vnd.docker.distribution.manifest.v2+json"
	config := make(map[string]interface{})
	res["config"] = config
	var layers []map[string]interface{}
	mediaType := "application/vnd.docker.image.rootfs.diff.tar.gzip"
	ls, ok := mf["Layers"].([]interface{})
	if !ok {
		log.Fatal("expected layers")
	}
	for _, ifc := range ls {
		layer, ok := ifc.(string)
		if !ok {
			log.Fatal("expected string")
		}
		obj := make(map[string]interface{})
		obj["mediaType"] = mediaType
		size, digest := compressLayer(tmp+"/"+layer, blobDir)
		obj["size"] = size
		obj["digest"] = "sha256:" + digest
		layers = append(layers, obj)
	}
	res["layers"] = layers
	return res
}

func compressLayer(path, blobDir string) (int64, string) {
	log.Printf("compressing layer: %s", path)
	tmp := blobDir + "/layer.tmp.tgz"

	gzipfile(path, tmp)

	digest := sha256File(tmp)
	size := fileSize(tmp)
	renameFile(tmp, blobDir+"/sha256:"+digest)

	return size, digest
}

func gzipfile(src, dst string) {
	data, _ := ioutil.ReadFile(src)
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(data)
	w.Close()

	err := ioutil.WriteFile(dst, b.Bytes(), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func fileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	return fi.Size()
}

func sha256File(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := sha256.New()

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

func renameFile(src, dst string) {
	if err := os.Rename(src, dst); err != nil {
		log.Fatal(err)
	}
}

func mkdir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func copyfile(src, dst string) {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
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

func readJSON(filepath string) (map[string]map[string]string, error) {
	body, _ := ioutil.ReadFile(filepath)
	var data map[string]map[string]string
	err := json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func readJSONArray(filepath string) ([]map[string]interface{}, error) {
	body, _ := ioutil.ReadFile(filepath)
	var data []map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func normalizeImageName(name string) string {
	// TODO
	return name
}
