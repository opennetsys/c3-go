package main

import (
	"github.com/miguelmota/c3/core/client"
)

func main() {
	cl := client.New()
	cl.ListImages()
}
