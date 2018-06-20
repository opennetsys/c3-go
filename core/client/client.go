package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
)

// Client ...
type Client struct {
	client *client.Client
}

// New ...
func New() *Client {
	httpclient := &http.Client{}

	if dockerCertPath := os.Getenv("DOCKER_CERT_PATH"); dockerCertPath != "" {
		options := tlsconfig.Options{
			CAFile:             filepath.Join(dockerCertPath, "ca.pem"),
			CertFile:           filepath.Join(dockerCertPath, "cert.pem"),
			KeyFile:            filepath.Join(dockerCertPath, "key.pem"),
			InsecureSkipVerify: os.Getenv("DOCKER_TLS_VERIFY") == "",
		}
		tlsc, err := tlsconfig.Client(options)
		if err != nil {
			log.Fatal(err)
		}

		httpclient.Transport = &http.Transport{
			TLSClientConfig: tlsc,
		}
	}

	host := os.Getenv("DOCKER_HOST")
	version := os.Getenv("DOCKER_VERSION")

	cl, err := client.NewClient(host, version, httpclient, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		client: cl,
	}
}

// ListImages ...
func (cl *Client) ListImages() {
	images, err := cl.client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, img := range images {
		fmt.Printf("ID: %s\n", img.ID)
	}
}
