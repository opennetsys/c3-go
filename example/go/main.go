package main

import (
	"fmt"

	"github.com/c3systems/c3/c3"
)

// Data ...
type Data struct {
	items map[string]string
}

func (s *Data) setItem(key, value string) error {
	s.items[key] = value
	return nil
}

func (s *Data) getItem(key string) string {
	return s.items[key]
}

func main() {
	fmt.Println("running")
	client := c3.New()
	data := &Data{}
	client.RegisterMethod("setItem", []string{"string", "string"}, data.setItem)
	client.RegisterMethod("getItem", []string{"string"}, data.getItem)
}
