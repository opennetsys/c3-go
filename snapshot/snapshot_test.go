package snapshot

import (
	"log"
	"testing"
)

func TestSnapshot(t *testing.T) {
	svc := New()
	imageHash := ""
	stateBlockNumber := 1
	err := svc.Snapshot(imageHash, stateBlockNumber)
	if err != nil {
		log.Fatal(err)
	}
}
