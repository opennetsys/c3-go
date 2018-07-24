package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3-go/common/hashing"
	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/common/stringutil"
	c3config "github.com/c3systems/c3-go/config"
	"github.com/c3systems/c3-go/core/server"
	loghooks "github.com/c3systems/c3-go/log/hooks"
)

var (
	// ErrMethodAlreadyRegistered ...
	ErrMethodAlreadyRegistered = errors.New("method already registered")
	// ErrMethodNotExists ...
	ErrMethodNotExists = errors.New("method does not exist")
)

// State ...
type State struct {
	state map[string]string
}

// C3 ...
type C3 struct {
	registeredMethods map[string]func(args ...interface{}) error
	receiver          chan []byte
	state             *State
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
		state: &State{
			state: map[string]string{},
		},
		statefile: c3config.TempContainerStateFilePath,
	}

	err := c3.setInitialState()
	if err != nil {
		log.Fatalf("[c3] %s", err)
	}

	go func() {
		err = c3.listen()
		if err != nil {
			log.Fatalf("[c3] %s", err)
		}
	}()

	return c3
}

// RegisterMethod ...
func (c3 *C3) RegisterMethod(methodName string, types []string, ifn interface{}) error {
	methodNameHash := hashing.HashToHexString([]byte(methodName))
	if _, ok := c3.registeredMethods[methodNameHash]; ok {
		return ErrMethodAlreadyRegistered
	}

	c3.registeredMethods[methodNameHash] = func(args ...interface{}) error {
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

			keyBytes, _ := hexutil.DecodeString(key)
			key = string(keyBytes)

			valueBytes, _ := hexutil.DecodeString(value)
			value = string(valueBytes)

			log.Printf("[c3] executed method %s with args: %s %s", methodName, key, value)
			err := v(key, value)
			if err != nil {
				log.Errorf("[c3] method failed %s", err)
				log.Fatalf("[c3] %s", err)
			}
		case func(string, string) (string, error):
			key, ok := args[0].(string)
			if !ok {
				return errors.New("not ok")
			}
			value, ok := args[1].(string)
			if !ok {
				return errors.New("not ok")
			}

			keyBytes, _ := hexutil.DecodeString(key)
			key = string(keyBytes)

			valueBytes, _ := hexutil.DecodeString(value)
			value = string(valueBytes)

			log.Printf("[c3] executed method %s with args: %s %s", methodName, key, value)
			str, err := v(key, value)
			if err != nil {
				log.Errorf("[c3] method failed %s", err)
				log.Fatalf("[c3] %s", err)
			}

			log.Printf("[c3] result %s", str)
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

// State ...
func (c3 *C3) State() *State {
	return c3.state
}

// Set ...
func (s *State) Set(key, value []byte) error {
	s.state[hexutil.EncodeToString(key)] = hexutil.EncodeToString(value)
	fmt.Println("setting state k/v", key, value)
	fmt.Println("latest state:", s.state)

	b, err := json.Marshal(s.state)
	if err != nil {
		return err
	}

	log.Printf("[c3] marshed state %s", string(b))

	f, err := os.OpenFile(c3config.TempContainerStateFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Error("[c3] failed to store file")
		return err
	}

	defer f.Close()
	f.Write(b)

	return nil
}

// Get ...
func (s *State) Get(key []byte) ([]byte, bool) {
	value, ok := s.state[hexutil.EncodeToString(key)]
	v, err := hexutil.DecodeString(value)
	if err != nil {
		return nil, false
	}
	return v, ok
}

func (c3 *C3) setInitialState() error {
	if _, err := os.Stat(c3.statefile); err == nil {
		src, err := ioutil.ReadFile(c3.statefile)
		if err != nil {
			log.Errorf("[c3] fail to read %s", err)
			return err
		}

		log.Printf("[c3] json data %s", string(src))

		if len(src) == 0 {
			return nil
		}

		log.Println("attempting to load initial state", string(src))

		b, err := stringutil.CompactJSON(src)
		if err != nil {
			log.Errorf("[c3] failed to compact %s", err)
			return err
		}

		err = json.Unmarshal(b, &c3.state.state)
		if err != nil {
			log.Errorf("[c3] fail to unmarshal %s", err)
			return err
		}
	} else {
		log.Error("[c3] state file not found")
	}

	log.Println("initial state loaded")

	return nil
}

// process ...
func (c3 *C3) process(payload []byte) error {
	var ifcs []interface{}
	if err := json.Unmarshal(payload, &ifcs); err != nil {
		log.Errorf("[c3] %s", err)
		return err
	}

	// if format is [a, b, c]
	_, ok := ifcs[0].(string)
	if ok {
		v := make([]string, len(ifcs))
		for i, k := range ifcs {
			v[i] = k.(string)
		}

		err := c3.invoke(v[0], v[1:])
		if err != nil {
			return err
		}

		return nil
	}

	// if format is [[a, b, c], [a, b, c]]
	for i := range ifcs {
		ifc := ifcs[i].([]interface{})
		v := make([]string, len(ifc))
		for j, k := range ifc {
			v[j] = k.(string)
		}

		err := c3.invoke(v[0], v[1:])
		if err != nil {
			return err
		}
	}

	return nil
}

// invoke ...
func (c3 *C3) invoke(method string, params []string) error {
	var args []interface{}
	for _, v := range params {
		args = append(args, v)
	}

	fn, ok := c3.registeredMethods[method]
	if !ok {
		return ErrMethodNotExists
	}
	return fn(args...)
}

// listen ...
func (c3 *C3) listen() error {
	for payload := range c3.receiver {
		err := c3.process(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
