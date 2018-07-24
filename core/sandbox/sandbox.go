package sandbox

import (
	"archive/tar"
	"bytes"
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

	"github.com/c3systems/c3-go/common/netutil"
	"github.com/c3systems/c3-go/common/stringutil"
	c3config "github.com/c3systems/c3-go/config"
	"github.com/c3systems/c3-go/core/docker"
	loghooks "github.com/c3systems/c3-go/log/hooks"
	"github.com/c3systems/c3-go/registry"
	regutil "github.com/c3systems/c3-go/registry/util"
)

// Ensure the service implements the interface
var _ Interface = (*Service)(nil)

// Service ...
type Service struct {
	docker            docker.Interface
	registry          registry.Interface
	sock              string
	runningContainers map[string]bool
	localIP           string
}

// Config ...
type Config struct {
	docker   docker.Interface
	registry registry.Interface
}

// New ...
func New(config *Config) *Service {
	localIP, err := netutil.LocalIP()
	if err != nil {
		log.Fatalf("[sandbox] %s", err)
	}

	if config == nil {
		dockerLocalRegistryHost := os.Getenv("DOCKER_LOCAL_REGISTRY_HOST")
		if dockerLocalRegistryHost == "" {
			dockerLocalRegistryHost = localIP.String()
		}

		docker := docker.NewClient()
		reg := registry.NewRegistry(&registry.Config{
			DockerLocalRegistryHost: dockerLocalRegistryHost,
		})

		config = &Config{
			docker:   docker,
			registry: reg,
		}
	}

	sb := &Service{
		docker:            config.docker,
		registry:          config.registry,
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
func (s *Service) Play(config *PlayConfig) ([]byte, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	var dockerImageID = config.ImageID
	var err error

	// If it's an IPFS hash then pull it from IPFS
	if strings.HasPrefix(config.ImageID, "Qm") {
		dockerizedHash := regutil.DockerizeHash(config.ImageID)
		hasImage, err := s.docker.HasImage(dockerizedHash)
		if err != nil {
			return nil, err
		}
		if hasImage {
			log.Printf("[sandbox] using cached image %s", dockerizedHash)
			dockerImageID = dockerizedHash
		} else {
			log.Printf("[sandbox] image not cached, pulling %s", config.ImageID)
			dockerImageID, err = s.registry.PullImage(config.ImageID)
			if err != nil {
				return nil, err
			}
		}
	}

	log.Printf("[sandbox] running docker image %s", dockerImageID)

	hp, err := netutil.GetFreePort()
	if err != nil {
		return nil, err
	}

	hostPort := strconv.Itoa(hp)
	containerID, err := s.docker.CreateContainer(dockerImageID, nil, &docker.CreateContainerConfig{
		Volumes: map[string]string{
		// sock binding will be required for spawning sibling containers
		// container:host
		//"/var/run/docker.sock": "/var/run/docker.sock",
		//"/tmp": tmpdir,
		},
		Ports: map[string]string{
			"3333": hostPort,
		},
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	body := config.InitialState
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{
		Name: c3config.TempContainerStateFileName,
		Mode: 0600,
		Size: int64(len(body)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(body)); err != nil {
		return nil, err
	}
	defer tw.Close()

	s.runningContainers[containerID] = true

	r := bytes.NewReader(buf.Bytes())
	err = s.docker.CopyToContainer(containerID, "/tmp", r)
	if err != nil {
		return nil, err
	}

	err = s.docker.StartContainer(containerID)
	if err != nil {
		return nil, err
	}

	done := make(chan bool)
	timedout := make(chan bool)
	errEvent := make(chan error)

	go func() {
		// Wait for application to start up
		// TODO: optimize
		time.Sleep(3 * time.Second)
		err := s.sendMessage(config.Payload, hostPort)
		if err != nil {
			log.Errorf("[sandbox] error sending message; %v", err)
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
			log.Error("[sandbox] writing to timeout channel")
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
			log.Errorf("[sandbox] error killing container on error event; %v", err)
			return nil, err
		}

		return nil, e
	case <-timedout:
		log.Error("[sandbox] timedout")
		close(timedout)
		close(errEvent)

		err := s.killContainer(containerID)
		if err != nil {
			log.Errorf("[sandbox] error killing container after timeout; %v", err)
			return nil, err
		}

		return nil, errors.New("timedout")
	case <-done:
		log.Println("[sandbox] reading new state...")
		cmd := []string{"bash", "-c", "cat " + c3config.TempContainerStateFilePath}
		resp, err := s.docker.ContainerExec(containerID, cmd)
		if err != nil {
			log.Errorf("[sandbox] error calling exec on contaienr; %v", err)
			return nil, err
		}

		result, err := parseNewState(resp)
		if err != nil {
			log.Errorf("[sandbox] error parsing new state; %v", err)
			return nil, err
		}

		log.Println("[sandbox] done")
		close(timedout)
		close(errEvent)
		timer.Stop()

		err = s.killContainer(containerID)
		if err != nil {
			log.Errorf("[sandbox] error killing container; %v", err)
			return nil, err
		}

		return result, nil
	}
}

func (s *Service) killContainer(containerID string) error {
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

func (s *Service) sendMessage(msg []byte, port string) error {
	// TODO: communicate over IPC
	host := fmt.Sprintf("%s:%s", s.localIP, port)
	log.Printf("[sandbox] sending message to container on host %s; message: %s", host, msg)
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Errorf("[sandbox] error sending message; %v", err)
		return err
	}
	defer conn.Close()
	log.Printf("[sandbox] writing to %s", host)
	conn.Write(msg)
	conn.Write([]byte("\n"))
	log.Printf("[sandbox] wrote to %s", host)
	// TODO: inspect container to see if it's completed the task
	time.Sleep(5 * time.Second)
	return nil
}

func (s *Service) cleanupOnExit() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	log.Printf("[sandbox] caught sig: %+v", sig)
	s.cleanup()
	os.Exit(0)
}

func (s *Service) cleanup() {
	for cid := range s.runningContainers {
		err := s.docker.StopContainer(cid)
		if err != nil {
			log.Errorf("[server] error %s", err)
		}
	}
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
