package sandbox

import (
	"errors"
	"log"
	"time"

	"github.com/c3systems/c3/core/dockerclient"
	"github.com/c3systems/c3/ditto"
)

// Sandbox ...
type Sandbox struct {
	docker *dockerclient.Client
	ditto  *ditto.Ditto
	sock   string
}

// Config ...
type Config struct {
}

// PlayConfig ...
type PlayConfig struct {
	ImageID    string // ipfs hash
	StateBlock string
}

// NewSandbox ...
func NewSandbox(config *Config) *Sandbox {
	dckr := dockerclient.New()
	dit := ditto.New(&ditto.Config{})

	return &Sandbox{
		docker: dckr,
		ditto:  dit,
		sock:   "/var/run/docker.sock",
	}
}

/*
[The] simplest way is to just expose the Docker socket to your CI container, by bind-mounting it with the -v flag.

Simply put, when you start your CI container (Jenkins or other), instead of hacking something together with Docker-in-Docker, start it with:

docker run -v /var/run/docker.sock:/var/run/docker.sock ...
Now this container will have access to the Docker socket, and will therefore be able to start containers. Except that instead of starting "child" containers, it will start "sibling" containers.
*/

// docker run -v /var/run/docker.sock:/var/run/docker.sock ...

// Play ...
func (s *Sandbox) Play(config *PlayConfig) error {
	dockerImageID, err := s.ditto.PullImage(config.ImageID)
	if err != nil {
		return err
	}

	containerID, err := s.docker.RunContainer(dockerImageID, []string{}, &dockerclient.RunContainerConfig{})
	if err != nil {
		return err
	}

	log.Printf("running container %s", containerID)

	done := make(chan bool)
	timedout := make(chan bool)

	go func() {
		select {
		case <-time.After(10 * time.Second):
			err := s.docker.StopContainer(containerID)
			if err != nil {
				log.Fatal(err)
			}

			timedout <- true
		}
	}()

	select {
	case <-timedout:
		return errors.New("timedout")
	case <-done:
		log.Println("done")
		return nil
	}
}
