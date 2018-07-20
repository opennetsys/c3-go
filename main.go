package main

import (
	"github.com/c3systems/c3/cmd"
	loghooks "github.com/c3systems/c3/log/hooks"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.AddHook(loghooks.ContextHook{})
	//bootstrap.Bootstrap()
	cmd.Execute()
}
