package c3

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/c3systems/c3/common/fscache"
	"github.com/c3systems/c3/config"
	"github.com/c3systems/c3/core/server"
)

var (
	// ErrMethodAlreadyRegistred ...
	ErrMethodAlreadyRegistred = errors.New("method already registered")
)

// C3 ...
type C3 struct {
	Store             store
	registeredMethods map[string]func(args ...interface{}) error
	receiver          chan []byte
}

// store
type store struct{}

// NewC3 ...
func NewC3() *C3 {
	receiver := make(chan []byte)
	c3 := &C3{
		registeredMethods: map[string]func(args ...interface{}) error{},
		receiver:          receiver,
	}

	go c3.listen()

	return c3
}

// RegisterMethod ...
func (c3 *C3) RegisterMethod(methodName string, types []string, ifn interface{}) error {
	if _, ok := c3.registeredMethods[methodName]; ok {
		return ErrMethodAlreadyRegistred
	}

	c3.registeredMethods[methodName] = func(args ...interface{}) error {
		switch v := ifn.(type) {
		case func(string, string) error:
			key, ok := args[0].(string)
			if !ok {
			}
			value, ok := args[0].(string)
			v(key, value)
		}
		return nil
	}
	return nil
}

// Serve ...
func (c3 *C3) Serve() {
	server.NewServer(&server.Config{
		Host:     config.ServerHost,
		Port:     config.ServerPort,
		Receiver: c3.receiver,
	}).Run()
}

// Set ...
// TODO: accept interfaces
func (s *store) Set(key, value string) {
	err := fscache.Set(key, value, 1*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
}

// Get ...
// TODO: accept interfaces
func (s *store) Get(key string) string {
	var value string
	found, err := fscache.Get(key, &value)
	if err != nil {
		log.Fatal(err)
	}
	if found {
		return value
	}

	return ""
}

// listen ...
func (c3 *C3) listen() {
	for payload := range c3.receiver {

		var parsed []string
		if err := json.Unmarshal(payload, &parsed); err != nil {
			log.Fatal(err)
		}

		method := parsed[0]
		var args []interface{}
		for _, v := range parsed[1:] {
			args = append(args, v)
		}
		if err := c3.registeredMethods[method](args...); err != nil {
			log.Fatal(err)
		}
	}
}
