package main

import (
	"fmt"

	"github.com/c3systems/c3/c3"
)

var client = c3.NewC3()

// App ...
type App struct {
}

func (s *App) setItem(key, value string) error {
	client.State().Set(key, value)
	return nil
}

func (s *App) getItem(key string) string {
	return client.State().Get(key)
}

func main() {
	fmt.Println("running")
	data := &App{}
	client.RegisterMethod("setItem", []string{"string", "string"}, data.setItem)
	client.RegisterMethod("getItem", []string{"string"}, data.getItem)
	client.Serve()
}
