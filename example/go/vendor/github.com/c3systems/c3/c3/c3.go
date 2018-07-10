package c3

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/c3systems/c3/common/stringutil"
	c3config "github.com/c3systems/c3/config"
	"github.com/c3systems/c3/core/server"
)

var (
	// ErrMethodAlreadyRegistred ...
	ErrMethodAlreadyRegistred = errors.New("method already registered")
)

// C3 ...
type C3 struct {
	Store             *store
	registeredMethods map[string]func(args ...interface{}) error
	receiver          chan []byte
	state             map[string]string
	statefile         string
}

// store
type store struct {
	*C3
}

// NewC3 ...
func NewC3() *C3 {
	receiver := make(chan []byte)
	c3 := &C3{
		registeredMethods: map[string]func(args ...interface{}) error{},
		receiver:          receiver,
		state:             map[string]string{},
		statefile:         c3config.TempContainerStateFilePath,
	}

	c3.Store = &store{
		c3,
	}

	go func() {
		err := c3.setInitialState()
		if err != nil {
			log.Fatal(err)
		}

		err = c3.listen()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return c3
}

// RegisterMethod ...
func (c3 *C3) RegisterMethod(methodName string, types []string, ifn interface{}) error {
	if _, ok := c3.registeredMethods[methodName]; ok {
		return ErrMethodAlreadyRegistred
	}

	c3.registeredMethods[methodName] = func(args ...interface{}) error {
		switch v := ifn.(type) {
		// TODO: accept arbitrary args
		case func(string, string) error:
			key, ok := args[0].(string)
			if !ok {
				return errors.New("not ok")
			}
			value, ok := args[1].(string)
			if !ok {
				return errors.New("not ok")
			}

			log.Printf("executed method %s with args: %s %s", methodName, key, value)
			v(key, value)
		}
		return nil
	}
	return nil
}

// Serve ...
func (c3 *C3) Serve() {
	server.NewServer(&server.Config{
		Host:     c3config.ServerHost,
		Port:     c3config.ServerPort,
		Receiver: c3.receiver,
	}).Run()
}

// Set ...
// TODO: accept interfaces
func (s *store) Set(key, value string) error {
	s.state[key] = value

	b, err := json.Marshal(s.state)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(s.statefile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	f.Write(b)
	f.Close()
	return nil
}

// Get ...
// TODO: accept interfaces
func (s *store) Get(key string) string {
	v := s.state[key]
	return v
}

func (c3 *C3) setInitialState() error {
	if _, err := os.Stat(c3.statefile); err == nil {
		src, err := ioutil.ReadFile(c3.statefile)
		if err != nil {
			log.Println("fail to read", err)
			return err
		}

		log.Println("json data", string(src))

		b, err := stringutil.CompactJSON(src)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &c3.state)
		if err != nil {
			log.Println("fail to unmarshal", err)
			return err
		}
	}

	return nil
}

// listen ...
func (c3 *C3) listen() error {
	for payload := range c3.receiver {

		var parsed []string
		if err := json.Unmarshal(payload, &parsed); err != nil {
			log.Println(err)
			return err
		}

		method := parsed[0]
		var args []interface{}
		for _, v := range parsed[1:] {
			args = append(args, v)
		}
		if err := c3.registeredMethods[method](args...); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
