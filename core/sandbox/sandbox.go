package sandbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/c3systems/c3/common/network"
	"github.com/c3systems/c3/common/stringutil"
	c3config "github.com/c3systems/c3/config"
	"github.com/c3systems/c3/core/docker"
	"github.com/c3systems/c3/registry"
)

// Sandbox ...
type Sandbox struct {
	docker            *docker.Client
	registry          *registry.Registry
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
	docker := docker.NewClient()
	dit := registry.NewRegistry(&registry.Config{})
	sb := &Sandbox{
		docker:            docker,
		registry:          dit,
		sock:              "/var/run/docker.sock",
		runningContainers: map[string]bool{},
	}

	go sb.cleanupOnExit()

	return sb
}

// PlayConfig ...
type PlayConfig struct {
	ImageID      string // ipfs hash
	Payload      []byte
	InitialState []byte
}

// Play in the sandbox
func (s *Sandbox) Play(config *PlayConfig) ([]byte, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	dockerImageID, err := s.registry.PullImage(config.ImageID)
	if err != nil {
		return nil, err
	}

	hp, err := network.GetFreePort()
	if err != nil {
		return nil, err
	}

	tmpdir, err := ioutil.TempDir("/tmp", "")
	if err != nil {
		return nil, err
	}

	fmt.Println("state loaded in tmp dir", tmpdir)
	hostStateFilePath := fmt.Sprintf("%s/%s", tmpdir, c3config.TempContainerStateFileName)

	err = ioutil.WriteFile(hostStateFilePath, config.InitialState, 0600)
	if err != nil {
		return nil, err
	}

	hostPort := strconv.Itoa(hp)
	containerID, err := s.docker.RunContainer(dockerImageID, []string{}, &docker.RunContainerConfig{
		Volumes: map[string]string{
			// sock binding will be required for spawning sibling containers
			// container:host
			//"/var/run/docker.sock": "/var/run/docker.sock",
			"/tmp": tmpdir,
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
	errEvent := make(chan error)

	go func() {
		// Wait for application to start up
		// TODO: optimize
		time.Sleep(1 * time.Second)
		err := s.sendMessage(config.Payload, hostPort)
		if err != nil {
			errEvent <- err
			return
		}

		done <- true
	}()

	timer := time.NewTimer(20 * time.Second)

	go func() {
		select {
		case <-timer.C:
			timedout <- true
		}
	}()

	select {
	case e := <-errEvent:
		timer.Stop()
		close(timedout)
		close(errEvent)
		err := s.killContainer(containerID)
		if err != nil {
			return nil, err
		}

		return nil, e
	case <-timedout:
		close(timedout)
		close(errEvent)

		err := s.killContainer(containerID)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("timedout")
	case <-done:
		log.Println("reading new state...")
		cmd := []string{"bash", "-c", "cat " + c3config.TempContainerStateFilePath}
		resp, err := s.docker.ContainerExec(containerID, cmd)
		if err != nil {
			return nil, err
		}

		result, err := parseNewState(resp)
		if err != nil {
			return nil, err
		}

		log.Println("done")
		close(timedout)
		close(errEvent)
		timer.Stop()

		err = s.killContainer(containerID)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
}

func (s *Sandbox) killContainer(containerID string) error {
	delete(s.runningContainers, containerID)
	if err := s.docker.StopContainer(containerID); err != nil {
		return err
	}

	return nil
}

func parseNewState(reader io.Reader) ([]byte, error) {
	var state map[string]string

	src, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	log.Println("new state", src, string(src))

	b, err := stringutil.CompactJSON(src)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &state)
	if err != nil {
		return nil, err
	}

	log.Println("new state", state)
	return b, nil
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
