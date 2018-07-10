package registry

import (
	"fmt"
	"log"

	"github.com/c3systems/c3/core/docker"
)

// Registry ...
type Registry struct {
	client *docker.Client
	host   string
}

// Config ...
type Config struct {
	Host string
}

// NewRegistry ...
func NewRegistry(config *Config) *Registry {
	client := docker.NewClient()
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
