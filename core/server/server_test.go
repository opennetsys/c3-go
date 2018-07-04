package server

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	server := New(&Config{
		Host: "localhost",
		Port: 3333,
	})
	_ = server
}

func TestRun(t *testing.T) {
	server := New(&Config{
		Host: "localhost",
		Port: 3333,
	})
	go server.Run()
	time.Sleep(1 * time.Second)
}
