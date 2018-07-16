package main

import (
	"github.com/c3systems/c3/cmd"
	"github.com/c3systems/c3/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.AddHook(logger.ContextHook{})
	//bootstrap.Bootstrap()
	cmd.Execute()
}
