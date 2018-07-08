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
	server := NewServer(&Config{
		Host: Host,
		Port: Port,
	})
	if server == nil {
		t.FailNow()
	}
}

func TestRun(t *testing.T) {
	server := NewServer(&Config{
		Host: Host,
		Port: Port,
	})
	go server.Run()
	time.Sleep(1 * time.Second)
}
