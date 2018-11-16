package main

import (
	"fmt"
	"os"

	"github.com/c3systems/c3-go/registry/util"
)

func main() {
	input := os.Args[1]
	output := util.DockerizeHash(input)
	fmt.Println(output)
}
