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

// NewSandbox ...
func NewSandbox(config *Config) *Sandbox {
	docker := dockerclient.NewClient()
	dit := ditto.NewDitto(&ditto.Config{})
	sb := &Sandbox{
		docker:            docker,
		ditto:             dit,
		sock:              "/var/run/docker.sock",
		runningContainers: map[string]bool{},
	}

	go sb.cleanupOnExit()

	return sb
}

// PlayConfig ...
type PlayConfig struct {
	ImageID string // ipfs hash
	Payload []byte
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

	info, err := s.docker.InspectContainer(containerID)
	if err != nil {
		return err
	}

	hostPort := info.NetworkSettings.Ports["3333/tcp"][0].HostPort

	done := make(chan bool)
	timedout := make(chan bool)

	go func() {
		// Wait for application to start up
		// TODO: optimize
		time.Sleep(1 * time.Second)
		err := s.sendMessage(config.Payload, hostPort)
		if err != nil {
			log.Fatal(err)
		}

		done <- true
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

func (s *Sandbox) sendMessage(msg []byte, port string) error {
	host := fmt.Sprintf("localhost:%s", port)
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Printf("writing to %s", host)
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
