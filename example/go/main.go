package main

import (
	"fmt"

	c3 "github.com/c3systems/c3-sdk-go"
)

var client = c3.NewC3()

// App ...
type App struct {
}

func (s *App) setItem(key, value string) error {
	client.State().Set([]byte(key), []byte(value))
	return nil
}

func (s *App) getItem(key string) string {
	v, found := client.State().Get([]byte(key))
	if !found {
		return ""
	}

	return string(v)
}

func main() {
	fmt.Println("running")
	data := &App{}
	client.RegisterMethod("setItem", []string{"string", "string"}, data.setItem)
	client.RegisterMethod("getItem", []string{"string"}, data.getItem)
	client.Serve()
}
