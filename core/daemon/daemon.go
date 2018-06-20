package daemon

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/docker/docker/daemon"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/libcontainerd"
	"github.com/docker/docker/registry"
)

// Daemon ...
type Daemon struct {
	daemon *daemon.Daemon
}

// New ...
func New() *Daemon {
	cfg := config.New()
	fmt.Println(cfg)
	registryService, err := registry.NewService(registry.ServiceOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(registryService)

	containerdRemote, err := libcontainerd.New(filepath.Join("/tmp/", "containerd"), filepath.Join("/tmp/", "containerd"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(containerdRemote)
	/*

		pluginStore := plugin.NewStore()

		daemon, err := docker.NewDaemon(cfg, registryService, containerdRemote, pluginStore)
		if err != nil {
			log.Fatal(err)
		}
	*/

	return &Daemon{
		//daemon: daemon,
	}
}
