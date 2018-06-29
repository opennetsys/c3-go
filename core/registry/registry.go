package registry

import (
	"fmt"

	"github.com/miguelmota/c3/core/dockerclient"
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
func (s Registry) PullImage(imageHash string) error {
	s.client.PullImage(fmt.Sprintf("%s/%s", s.host, imageHash))
	return nil
}
