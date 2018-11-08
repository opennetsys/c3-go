package cmd

import (
	"errors"

	colorlog "github.com/c3systems/c3-go/log/color"
)

// error wrap
func errw(err error) error {
	return errors.New(colorlog.Red(err.Error()))
}
