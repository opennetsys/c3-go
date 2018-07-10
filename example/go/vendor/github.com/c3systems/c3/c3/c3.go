package c3

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/c3systems/c3/config"
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
		statefile:         config.TempContainerStatePath,
	}

	c3.Store = &store{
		c3,
	}

	go c3.setInitialState()
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
		// TODO: accept arbitrary args
		case func(string, string) error:
			key, ok := args[0].(string)
			if !ok {
				log.Fatal("not ok")
			}
			value, ok := args[1].(string)
			if !ok {
				log.Fatal("not ok")
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
		Host:     config.ServerHost,
		Port:     config.ServerPort,
		Receiver: c3.receiver,
	}).Run()
}

// Set ...
// TODO: accept interfaces
func (s *store) Set(key, value string) {
	s.state[key] = value

	b, err := json.Marshal(s.state)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(s.statefile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	f.Write(b)
	f.Close()
}

// Get ...
// TODO: accept interfaces
func (s *store) Get(key string) string {
	v := s.state[key]
	return v
}

func (c3 *C3) setInitialState() {
	if _, err := os.Stat(c3.statefile); err == nil {
		src, err := ioutil.ReadFile(c3.statefile)
		if err != nil {
			log.Fatalln("fail to read", err)
		}

		b := new(bytes.Buffer)

		re := regexp.MustCompile(`\\n`)
		s := re.ReplaceAllString(string(src), "")
		fmt.Println("raw data", s)

		if err := json.Compact(b, []byte(s)); err != nil {
			log.Fatalln("fail to compact", err)
		}

		log.Println("json data", string(b.Bytes()))

		err = json.Unmarshal(b.Bytes(), &c3.state)
		if err != nil {
			log.Fatalln("fail to unmarshal", err)
		}
	}
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
