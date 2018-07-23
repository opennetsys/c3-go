package main

import (
	"github.com/c3systems/c3-go/cmd"
	loghooks "github.com/c3systems/c3-go/log/hooks"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.AddHook(loghooks.ContextHook{})
	//bootstrap.Bootstrap()
	cmd.Execute()
}
