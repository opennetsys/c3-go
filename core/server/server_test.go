package server

import (
	"testing"
)

func TestNew(t *testing.T) {
	server := New(&Config{
		Host: "localhost",
		Port: 3333,
	})
	server.Run()
}
