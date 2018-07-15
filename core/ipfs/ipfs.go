package ipfs

import (
	"errors"
	"log"
	"os/exec"
	"strings"
	"time"

	api "github.com/ipfs/go-ipfs-api"
)

// Client ...
type Client struct {
	client *api.Shell
}

// NewClient ...
func NewClient() *Client {
	err := RunDaemon()
	if err != nil {
		log.Fatal(err)
	}

	url, err := getIpfsAPIURL()
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewShell(url)
	return &Client{
		client: client,
	}
}

// Get ...
func (client *Client) Get(hash, outdir string) error {
	return client.client.Get(hash, outdir)
}

// AddDir ...
func (client *Client) AddDir(dir string) (string, error) {
	return client.client.AddDir(dir)
}

// Refs ...
func (client *Client) Refs(hash string, recursive bool) (<-chan string, error) {
	return client.client.Refs(hash, recursive)
}

// RunDaemon ...
func RunDaemon() error {
	var err error
	ready := make(chan bool)
	go func() {
		if err = spawnIpfsDaemon(ready); err != nil {
			log.Println(err)
		}
	}()

	if !<-ready {
		return errors.New("failed to start IPFS daemon")
	}

	return nil
}

// hacky way to spawn daemon
// TODO: improve
func spawnIpfsDaemon(ready chan bool) error {
	out, err := exec.Command("pgrep", "ipfs").Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		log.Println("IPFS is not running. Starting...")

		go func() {
			// TODO: detect when running by watching log output
			time.Sleep(10 * time.Second)
			ready <- true
		}()

		err := exec.Command("ipfs", "daemon").Run()
		if err != nil {
			ready <- false
			log.Println(err)
			return errors.New("failed to start IPFS")
		}
	}

	ready <- true
	log.Println("IPFS is running...")
	return nil
}

func getIpfsAPIURL() (string, error) {
	out, err := exec.Command("ipfs", "config", "Addresses.API").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
