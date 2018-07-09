package c3

import (
	"errors"

	"github.com/c3systems/c3/config"
	"github.com/c3systems/c3/core/server"
)

// Service ...
type Service struct {
	registeredMethods map[string]func(args ...interface{}) interface{}
}

// New ...
func New() *Service {
	go server.New(&server.Config{
		Host: config.ServerHost,
		Port: config.ServerPort,
	}).Run()

	return &Service{
		registeredMethods: map[string]func(args ...interface{}) interface{}{},
	}
}

// RegisterMethod ...
func (s *Service) RegisterMethod(methodName string, types []string, ifn interface{}) error {
	if _, ok := s.registeredMethods[methodName]; ok {
		return errors.New("method already registered")
	}

	s.registeredMethods[methodName] = func(args ...interface{}) interface{} {
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
