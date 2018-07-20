package server

import (
	"bufio"
	"fmt"
	"net"

	loghooks "github.com/c3systems/c3/log/hooks"
	log "github.com/sirupsen/logrus"
)

// Server ...
type Server struct {
	host     string
	port     int
	receiver chan []byte
}

// Client ...
type Client struct {
	conn    net.Conn
	channel chan []byte
}

// Config ...
type Config struct {
	Host     string
	Port     int
	Receiver chan []byte
}

// NewServer ...
func NewServer(config *Config) *Server {
	return &Server{
		host:     config.Host,
		port:     config.Port,
		receiver: config.Receiver,
	}
}

// Run ...
func (server *Server) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%v", server.host, server.port))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		client := &Client{
			conn:    conn,
			channel: server.receiver,
		}
		go client.handleRequest()
	}
}

func (client *Client) handleRequest() {
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}
		fmt.Printf("Message incoming: %s", message)
		client.channel <- []byte(message)
		client.conn.Write([]byte("Message received.\n"))
	}
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
