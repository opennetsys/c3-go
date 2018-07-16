package bootstrap

import (
	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3/core/ipfs"
)

// Bootstrap ...
func Bootstrap() {
	err := ipfs.RunDaemon()
	if err != nil {
		log.Fatal("Failed to start IPFS daemon")
	}
}
