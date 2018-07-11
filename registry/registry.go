package registry

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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/c3systems/c3/common/network"
	c3config "github.com/c3systems/c3/config"
	"github.com/c3systems/c3/core/docker"
	"github.com/c3systems/c3/registry/server"
	"github.com/c3systems/c3/registry/util"
	"github.com/davecgh/go-spew/spew"
)

// Registry ...
type Registry struct {
}

// Config ...
type Config struct {
}

// NewRegistry ...
func NewRegistry(config *Config) *Registry {
	return &Registry{}
}

// PushImageByID uploads Docker image by image ID (hash/name) to IPFS
func (registry *Registry) PushImageByID(imageID string) error {
	client := docker.NewClient()
	reader, err := client.ReadImage(imageID)
	if err != nil {
		return err
	}

	return registry.PushImage(reader)
}

// PushImage uploads Docker image to IPFS
func (registry *Registry) PushImage(reader io.Reader) error {
	tmp, err := mktmp()
	if err != nil {
		return err
	}
	fmt.Println("temp:", tmp)

	if err := untar(reader, tmp); err != nil {
		return err
	}

	root, err := ipfsPrep(tmp)
	if err != nil {
		return err
	}

	imageIpfsHash, err := uploadDir(root)
	if err != nil {
		return err
	}

	fmt.Printf("\nuploaded to /ipfs/%s\n", imageIpfsHash)

	fmt.Printf("docker image %s\n", util.DockerizeHash(imageIpfsHash))

	return nil
}

// DownloadImage download Docker image from IPFS
func (registry *Registry) DownloadImage(ipfsHash string) (string, error) {
	tmp, err := mktmp()
	if err != nil {
		return "", err
	}

	path := tmp + "/" + ipfsHash + ".tar"
	outstr, errstr, err := ipfsCmd(fmt.Sprintf("get %s -a -o %s", ipfsHash, path))
	if err != nil {
		return "", err
	}
	_ = outstr
	_ = errstr

	return path, nil
}

// PullImage pull Docker image from IPFS
func (registry *Registry) PullImage(ipfsHash string) (string, error) {
	go server.Run()
	client := docker.NewClient()

	localIP, err := network.LocalIP()
	if err != nil {
		return "", err
	}

	dockerImageID := fmt.Sprintf("%s:%v/%s", localIP.String(), c3config.DockerRegistryPort, util.DockerizeHash(ipfsHash))

	log.Printf("attempting to pull %s", dockerImageID)

	err = client.PullImage(dockerImageID)
	if err != nil {
		return dockerImageID, err
	}

	return dockerImageID, nil
}

func mktmp() (string, error) {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	return tmp, err
}

func ipfsPrep(tmp string) (string, error) {
	root, err := mktmp()
	if err != nil {
		return "", err
	}

	workdir := root
	fmt.Println("preparing image in:", workdir)
	name := "default"

	// read human readable name of image
	if _, err := os.Stat(tmp + "repositories"); err == nil {
		reposJSON, err := readJSON(tmp + "/repositories")
		if err != nil {
			return "", err
		}
		if len(reposJSON) != 1 {
			return "", errors.New("only one repository expected in input file")
		}
		for imageName, tags := range reposJSON {
			fmt.Println(imageName, tags)
			if len(tags) != 1 {
				return "", fmt.Errorf("only one tag expected for %s", imageName)
			}
			for tag, hash := range tags {
				name = normalizeImageName(imageName)
				fmt.Printf("processing image:%s tag:%s hash:256:%s", name, tag, hash)
			}
		}
	}

	workdir = workdir + "/" + name
	mkdir(workdir)
	mkdir(workdir + "/manifests")
	mkdir(workdir + "/blobs")
	manifestJSON, err := readJSONArray(tmp + "/manifest.json")
	if err != nil {
		return "", err
	}

	if len(manifestJSON) == 0 {
		return "", errors.New("expected manifest to contain data")
	}

	manifest := manifestJSON[0]
	configFile, ok := manifest["Config"].(string)
	if !ok {
		return "", errors.New("image archive must be produced by docker > 1.10")
	}

	configDest := workdir + "/blobs/sha256:" + string(configFile[:len(configFile)-5])
	fmt.Println("\nDIST", configDest)
	mkdir(configDest)
	if err := copyFile(tmp+"/"+configFile, configDest+"/"+configFile); err != nil {
		return "", err
	}

	mf, err := makeV2Manifest(manifest, configFile, configDest, tmp, workdir)
	if err != nil {
		return "", err
	}

	spew.Dump(mf)

	err = writeJSON(mf, workdir+"/manifests/latest-v2")
	if err != nil {
		return "", err
	}

	return root, nil
}

func uploadDir(root string) (string, error) {
	outstr, errstr, err := ipfsCmd(fmt.Sprintf("add -r -q %s", root))
	if err != nil {
		return "", err
	}

	if errstr != "" {
		return "", errors.New(errstr)
	}
	if outstr != "" {
		hashes := strings.Split(outstr, "\n")
		imageIpfsHash := hashes[len(hashes)-2 : len(hashes)-1][0]
		return imageIpfsHash, nil
	}

	return "", errors.New("no result")
}

func ipfsCmd(cmdStr string) (string, string, error) {
	path, err := exec.LookPath("ipfs")
	if err != nil {
		return "", "", errors.New("ipfs command was not found. Please install ipfs")
	}
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", path, cmdStr))
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var stdoutBuf, stderrBuf bytes.Buffer
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err = cmd.Start()
	if err != nil {
		return "", "", err
	}

	go copyio(stdoutIn, stdout)
	go copyio(stderrIn, stderr)

	err = cmd.Wait()
	if err != nil {
		return "", "", err
	}

	outstr := strings.TrimSpace(string(stdoutBuf.Bytes()))
	errstr := strings.TrimSpace(string(stderrBuf.Bytes()))

	return outstr, errstr, nil
}

func copyio(out io.Reader, in io.Writer) error {
	_, err := io.Copy(in, out)
	if err != nil {
		return err
	}

	return nil
}

func writeJSON(idate interface{}, path string) error {
	data, err := json.Marshal(idate)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// produce v2 manifest of type/application/vnd.docker.distribution.manifest.v2+json
func makeV2Manifest(manifest map[string]interface{}, configFile, configDest, tmp, workdir string) (map[string]interface{}, error) {
	v2manifest, err := prepareV2Manifest(manifest, tmp, workdir+"/blobs")
	if err != nil {
		return nil, err
	}
	config := make(map[string]interface{})
	config["digest"] = "sha256:" + string(configFile[:len(configFile)-5])
	config["size"], err = fileSize(configDest + "/" + configFile)
	if err != nil {
		return nil, err
	}
	config["mediaType"] = "application/vnd.docker.container.image.v1+json"
	conf, ok := v2manifest["config"].(map[string]interface{})
	if !ok {
		return nil, errors.New("not ok")
	}
	v2manifest["config"] = mergemap(conf, config)
	return v2manifest, nil
}

func mergemap(a, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func prepareV2Manifest(mf map[string]interface{}, tmp, blobDir string) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	res["schemaVersion"] = 2
	res["mediaType"] = "application/vnd.docker.distribution.manifest.v2+json"
	config := make(map[string]interface{})
	res["config"] = config
	var layers []map[string]interface{}
	mediaType := "application/vnd.docker.image.rootfs.diff.tar.gzip"
	ls, ok := mf["Layers"].([]interface{})
	if !ok {
		return nil, errors.New("expected layers")
	}
	for _, ifc := range ls {
		layer, ok := ifc.(string)
		if !ok {
			return nil, errors.New("expected string")
		}
		obj := make(map[string]interface{})
		obj["mediaType"] = mediaType
		size, digest, err := compressLayer(tmp+"/"+layer, blobDir)
		if err != nil {
			return nil, err
		}
		obj["size"] = size
		obj["digest"] = "sha256:" + digest
		layers = append(layers, obj)
	}
	res["layers"] = layers
	return res, nil
}

func compressLayer(path, blobDir string) (int64, string, error) {
	log.Printf("compressing layer: %s", path)
	tmp := blobDir + "/layer.tmp.tgz"

	err := gzipFile(path, tmp)
	if err != nil {
		return int64(0), "", err
	}

	digest, err := sha256File(tmp)
	if err != nil {
		return int64(0), "", err
	}

	size, err := fileSize(tmp)
	if err != nil {
		return int64(0), "", err
	}

	err = renameFile(tmp, blobDir+"/sha256:"+digest)
	if err != nil {
		return int64(0), "", err
	}

	return size, digest, nil
}

func gzipFile(src, dst string) error {
	data, _ := ioutil.ReadFile(src)
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(data)
	w.Close()

	return ioutil.WriteFile(dst, b.Bytes(), os.ModePerm)
}

func fileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return int64(0), err
	}

	return fi.Size(), nil
}

func sha256File(path string) (string, error) {
	// TODO: stream instead of reading whole image in memory
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func renameFile(src, dst string) error {
	if err := os.Rename(src, dst); err != nil {
		return err
	}

	return nil
}

func mkdir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func copyFile(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, 0644)
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