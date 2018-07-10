package sandbox

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/c3systems/c3/common/network"
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
	if config == nil {
		config = &Config{}
	}
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

// Play in the sandbox
func (s *Sandbox) Play(config *PlayConfig) ([]byte, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	dockerImageID, err := s.ditto.PullImage(config.ImageID)
	if err != nil {
		return nil, err
	}

	hp, err := network.GetFreePort()
	if err != nil {
		return nil, err
	}

	hostPort := strconv.Itoa(hp)

	containerID, err := s.docker.RunContainer(dockerImageID, []string{}, &dockerclient.RunContainerConfig{
		Volumes: map[string]string{
		// sock binding will be required for spawning sibling containers
		//"/var/run/docker.sock": "/var/run/docker.sock",
		//"/tmp":                 "/tmp",
		},
		Ports: map[string]string{
			"3333": hostPort,
		},
	})
	if err != nil {
		return nil, err
	}

	s.runningContainers[containerID] = true

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

	timer := time.NewTimer(20 * time.Second)

	go func() {
		select {
		case <-timer.C:
			timedout <- true
			err := s.docker.StopContainer(containerID)
			if err != nil {
				log.Fatal(err)
			}

			delete(s.runningContainers, containerID)
		}
	}()

	// TODO: return new state

	select {
	case <-timedout:
		close(timedout)
		return nil, errors.New("timedout")
	case <-done:
		log.Println("reading new state...")
		result, err := s.docker.ContainerExec(containerID)
		if err != nil {
			return nil, err
		}

		log.Println("done")
		close(timedout)
		timer.Stop()
		err = s.docker.StopContainer(containerID)
		if err != nil {
			return nil, err
		}

		delete(s.runningContainers, containerID)

		return result, nil
	}
}

func (s *Sandbox) sendMessage(msg []byte, port string) error {
	log.Printf("sending message %s", msg)
	// TODO: communicate over IPC
	host := fmt.Sprintf("localhost:%s", port)
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("writing to %s", host)
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
