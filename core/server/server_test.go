package server

import (
	"testing"
	"time"
)

var (
	Host = "localhost"
	Port = 3333
)

func TestNew(t *testing.T) {
	t.Parallel()
	server := NewServer(&Config{
		Host: Host,
		Port: Port,
	})
	if server == nil {
		t.Error("expected instance")
	}
}

func TestRun(t *testing.T) {
	t.Parallel()
	server := NewServer(&Config{
		Host: Host,
		Port: Port,
	})
	go server.Run()
	time.Sleep(1 * time.Second)
}
