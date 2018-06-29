package dockerclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	if host == "" {
		log.Fatal("DOCKER_HOST is required")
	}

	if version == "" {
		version = dockerVersionFromCLI()
		if version == "" {
			log.Fatal("DOCKER_VERSION is required")
		}
	}

	cl, err := client.NewClient(host, version, httpclient, nil)
	//cl, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		client: cl,
	}
}

// ListImages ...
func (s *Client) ListImages() {
	images, err := s.client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, img := range images {
		fmt.Printf("ID: %s\n", img.ID)
	}
}

// PullImage ...
func (s *Client) PullImage(imageURL string) error {
	reader, err := s.client.ImagePull(context.Background(), imageURL, types.ImagePullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, reader)
	return nil
}

// PushImage ...
func (s *Client) PushImage(imageURL string) error {
	reader, err := s.client.ImagePush(context.Background(), imageURL, types.ImagePushOptions{
		RegistryAuth: "123", // if no auth, then any value is required
	})
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, reader)
	return nil
}

// RunContainer ...
func (s *Client) RunContainer(imageID string, cmd []string) {
	/*
		resp, err := s.client.ContainerCreate(context.Background(), &container.Config{
			Image: imageID,
			Cmd:   cmd,
			Tty:   true,
		}, nil, nil, "")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp)
	*/
}

func dockerVersionFromCLI() string {
	cmd := `docker version --format="{{.Client.APIVersion}}"`
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
