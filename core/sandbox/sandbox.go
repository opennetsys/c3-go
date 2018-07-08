package sandbox

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/c3systems/c3/core/dockerclient"
	"github.com/c3systems/c3/ditto"
)

// Sandbox ...
type Sandbox struct {
	docker            *dockerclient.Client
	ditto             *ditto.Ditto
	sock              string
	runningContainers map[string]bool
}

// Config ...
type Config struct {
}

// PlayConfig ...
type PlayConfig struct {
	ImageID string // ipfs hash
	Payload []byte
}

// NewSandbox ...
func NewSandbox(config *Config) *Sandbox {
	dckr := dockerclient.New()
	dit := ditto.New(&ditto.Config{})
	sb := &Sandbox{
		docker:            dckr,
		ditto:             dit,
		sock:              "/var/run/docker.sock",
		runningContainers: map[string]bool{},
	}

	go sb.cleanupOnExit()

	return sb
}

// TODO: include transaction inputs

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
	s.runningContainers[containerID] = true

	done := make(chan bool)
	timedout := make(chan bool)

	go func() {
		// Wait for application to start up
		// TODO: optimize
		time.Sleep(1 * time.Second)
		err := s.sendMessage(config.Payload)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		select {
		case <-time.After(1 * time.Minute):
			err := s.docker.StopContainer(containerID)
			if err != nil {
				log.Fatal(err)
			}

			delete(s.runningContainers, containerID)
			timedout <- true
		}
	}()

	// TODO: return new state

	select {
	case <-timedout:
		return errors.New("timedout")
	case <-done:
		log.Println("done")
		return nil
	}
}

func (s *Sandbox) sendMessage(msg []byte) error {
	// TODO: use dynamic port
	conn, err := net.Dial("tcp", "localhost:3333")
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.Write(msg)
	conn.Write([]byte("\n"))
	return nil
}

func (s *Sandbox) cleanupOnExit() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	fmt.Printf("caught sig: %+v", sig)
	s.cleanup()
	os.Exit(0)
}

func (s *Sandbox) cleanup() {
	for cid := range s.runningContainers {
		err := s.docker.StopContainer(cid)
		if err != nil {
			log.Println("error", err)
		}
	}
}
