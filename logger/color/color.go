package color

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var magenta = color.New(color.FgMagenta).SprintFunc()
var white = color.New(color.FgWhite).SprintFunc()

// Green ...
func Green(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Println(green(s))
	} else {
		log.Printf(green(s), args...)
	}
}

// Red ...
func Red(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Println(red(s))
	} else {
		log.Printf(red(s), args...)
	}
}

// Cyan ...
func Cyan(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Println(cyan(s))
	} else {
		log.Printf(cyan(s), args...)
	}
}

// Yellow ...
func Yellow(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Println(yellow(s))
	} else {
		log.Printf(yellow(s), args...)
	}
}

// Magenta ...
func Magenta(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Println(magenta(s))
	} else {
		log.Printf(magenta(s), args...)
	}
}

// White ...
func White(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Println(white(s))
	} else {
		log.Printf(white(s), args...)
	}
}

// Info ...
func Info(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Errorln(cyan(s))
	} else {
		log.Errorf(cyan(s), args...)
	}
}

// Panic ...
func Panic(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Panicln(red(s))
	} else {
		log.Panicf(red(s), args...)
	}
}

// Error ...
func Error(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Errorln(red(s))
	} else {
		log.Errorf(red(s), args...)
	}
}

// Warn ...
func Warn(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Warnln(yellow(s))
	} else {
		log.Warnf(yellow(s), args...)
	}
}

// Debug ...
func Debug(s string, args ...interface{}) {
	if len(args) == 0 {
		log.Debugln(white(s))
	} else {
		log.Debugf(white(s), args...)
	}
}
