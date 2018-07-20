package netutil

import (
	"net"
	"strconv"
	"testing"
)

func TestGetFreePort(t *testing.T) {
	t.Parallel()
	port, err := GetFreePort()
	if err != nil {
		t.Error(err)
	}
	if port == 0 {
		t.Error("port == 0")
	}

	// Try to listen on the port
	l, err := net.Listen("tcp", "localhost"+":"+strconv.Itoa(port))
	if err != nil {
		t.Error(err)
	}
	defer l.Close()
}

func TestLocalIP(t *testing.T) {
	t.Parallel()
	ip, err := LocalIP()
	if err != nil {
		t.Error(err)
	}

	if ip.String() == "" {
		t.Error("expected IP address")
	}

	t.Log(ip)
}
