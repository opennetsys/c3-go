package c3

import (
	"errors"

	"github.com/c3systems/c3/config"
	"github.com/c3systems/c3/core/server"
)

var (
	// ErrMethodAlreadyRegistred ...
	ErrMethodAlreadyRegistred = errors.New("method already registered")
)

// C3 ...
type C3 struct {
	registeredMethods map[string]func(args ...interface{}) interface{}
}

// NewC3 ...
func NewC3() *C3 {
	return &C3{
		registeredMethods: map[string]func(args ...interface{}) interface{}{},
	}
}

// RegisterMethod ...
func (c3 *C3) RegisterMethod(methodName string, types []string, ifn interface{}) error {
	if _, ok := c3.registeredMethods[methodName]; ok {
		return ErrMethodAlreadyRegistred
	}

	c3.registeredMethods[methodName] = func(args ...interface{}) interface{} {
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
	server.New(&server.Config{
		Host: config.ServerHost,
		Port: config.ServerPort,
	}).Run()
}
