package sandbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3/common/network"
	"github.com/c3systems/c3/common/stringutil"
	c3config "github.com/c3systems/c3/config"
	"github.com/c3systems/c3/core/docker"
	loghooks "github.com/c3systems/c3/logger/hooks"
	"github.com/c3systems/c3/registry"
)

// Sandbox ...
type Sandbox struct {
	docker            docker.Interface
	registry          registry.Interface
	sock              string
	runningContainers map[string]bool
	localIP           string
}

// Config ...
type Config struct {
}

// NewSandbox ...
func NewSandbox(config *Config) Interface {
	if config == nil {
		config = &Config{}
	}

	localIP, err := network.LocalIP()
	if err != nil {
		log.Fatalf("[sandbox] %s", err)
	}

	dockerLocalRegistryHost := os.Getenv("DOCKER_LOCAL_REGISTRY_HOST")
	if dockerLocalRegistryHost == "" {
		dockerLocalRegistryHost = localIP.String()
	}

	docker := docker.NewClient()
	dit := registry.NewRegistry(&registry.Config{
		DockerLocalRegistryHost: dockerLocalRegistryHost,
	})
	sb := &Sandbox{
		docker:            docker,
		registry:          dit,
		sock:              "/var/run/docker.sock",
		runningContainers: map[string]bool{},
		localIP:           localIP.String(),
	}

	//go sb.cleanupOnExit()

	return sb
}

// PlayConfig ...
type PlayConfig struct {
	ImageID      string // can be ipfs hash
	Payload      []byte
	InitialState []byte
}

// Play in the sandbox
func (s *Sandbox) Play(config *PlayConfig) ([]byte, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	var dockerImageID = config.ImageID
	var err error

	// If it's an IPFS hash then pull it from IPFS
	if strings.HasPrefix(config.ImageID, "Qm") {
		dockerImageID, err = s.registry.PullImage(config.ImageID)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("[sandbox] running docker image %s", dockerImageID)

	hp, err := network.GetFreePort()
	if err != nil {
		return nil, err
	}

	tmpdir, err := ioutil.TempDir("/tmp", "")
	if err != nil {
		return nil, err
	}

	hostStateFilePath := fmt.Sprintf("%s/%s", tmpdir, c3config.TempContainerStateFileName)

	// don't write empty file
	if config.InitialState != nil && len(config.InitialState) > 0 {
		err := ioutil.WriteFile(hostStateFilePath, config.InitialState, 0600)
		if err != nil {
			return nil, err
		}
	}

	log.Println("[sandbox] state loaded in tmp dir", tmpdir)

	hostPort := strconv.Itoa(hp)
	containerID, err := s.docker.RunContainer(dockerImageID, nil, &docker.RunContainerConfig{
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
		time.Sleep(3 * time.Second)
		err := s.sendMessage(config.Payload, hostPort)
		if err != nil {
			log.Printf("[sandbox] error sending message; %v", err)
			errEvent <- err
			return
		}

		log.Println("[sandbox] writing to done channel")

		done <- true
	}()

	timer := time.NewTimer(1 * time.Minute)

	go func() {
		select {
		case <-timer.C:
			log.Println("[sandbox] writing to timeout channel")
			timedout <- true
		}
	}()

	select {
	case e := <-errEvent:
		log.Errorf("[sandbox] error; %v", e)
		timer.Stop()
		close(timedout)
		close(errEvent)
		err := s.killContainer(containerID)
		if err != nil {
			log.Printf("[sandbox] error killing container on error event; %v", err)
			return nil, err
		}

		return nil, e
	case <-timedout:
		log.Error("[sandbox] timedout")
		close(timedout)
		close(errEvent)

		err := s.killContainer(containerID)
		if err != nil {
			log.Printf("[sandbox] error killing container after timeout; %v", err)
			return nil, err
		}

		return nil, errors.New("timedout")
	case <-done:
		log.Println("[sandbox] reading new state...")
		cmd := []string{"bash", "-c", "cat " + c3config.TempContainerStateFilePath}
		resp, err := s.docker.ContainerExec(containerID, cmd)
		if err != nil {
			log.Printf("[sandbox] error calling exec on contaienr; %v", err)
			return nil, err
		}

		result, err := parseNewState(resp)
		if err != nil {
			log.Printf("[sandbox] error parsing new state; %v", err)
			return nil, err
		}

		log.Println("[sandbox] done")
		close(timedout)
		close(errEvent)
		timer.Stop()

		err = s.killContainer(containerID)
		if err != nil {
			log.Printf("[sandbox] error killing container; %v", err)
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

	log.Printf("[sandbox] new state json %s", string(src))

	b, err := stringutil.CompactJSON(src)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &state)
	if err != nil {
		return nil, err
	}

	log.Printf("[sandbox] new state %s", state)
	return b, nil
}

func (s *Sandbox) sendMessage(msg []byte, port string) error {
	// TODO: communicate over IPC
	host := fmt.Sprintf("%s:%s", s.localIP, port)
	log.Printf("[sandbox] sending message to container on host %s; message: %s", host, msg)
	time.Sleep(15 * time.Second)
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Printf("[sandbox] error sending message; %v", err)
		return err
	}
	defer conn.Close()
	log.Printf("[sandbox] writing to %s", host)
	conn.Write(msg)
	conn.Write([]byte("\n"))
	log.Printf("[sandbox] wrote to %s", host)
	return nil
}

func (s *Sandbox) cleanupOnExit() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	log.Printf("[sandbox] caught sig: %+v", sig)
	s.cleanup()
	os.Exit(0)
}

func (s *Sandbox) cleanup() {
	for cid := range s.runningContainers {
		err := s.docker.StopContainer(cid)
		if err != nil {
			log.Printf("[server] error %s", err)
		}
	}
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
