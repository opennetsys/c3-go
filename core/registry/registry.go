package registry

import (
	"fmt"
	"log"

	"github.com/c3systems/c3/core/dockerclient"
)

// Registry ...
type Registry struct {
	client *dockerclient.Client
	host   string
}

// Config ...
type Config struct {
	Host string
}

// New ...
func New(config *Config) *Registry {
	client := dockerclient.New()
	return &Registry{
		client: client,
		host:   config.Host,
	}
}

// PullImage ...
func (s Registry) PullImage(imageID string) error {
	err := s.client.PullImage(fmt.Sprintf("%s/%s", s.host, imageID))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// PushImage ...
func (s Registry) PushImage(imageID string) error {
	err := s.client.PushImage(fmt.Sprintf("%s/%s", s.host, imageID))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
