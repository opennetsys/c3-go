package daemon

import (
	"fmt"

	"github.com/docker/docker/daemon/config"
)

// Daemon ...
type Daemon struct {
	daemon *docker.Daemon
}

// New ...
func New() *Daemon {
	cfg := config.New()
	fmt.Println(cfg)
	/*
		registryService, err := registry.NewService()
		if err != nil {
			return err
		}

		containerdRemote, err := libcontainerd.New(filepath.Join("/", "containerd"), filepath.Join("/", "containerd"))
		if err != nil {
			return err
		}

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
