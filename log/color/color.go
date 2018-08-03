package color

import (
	"github.com/fatih/color"
)

var green = color.New(color.FgGreen)
var red = color.New(color.FgRed)
var cyan = color.New(color.FgCyan)
var yellow = color.New(color.FgYellow)
var magenta = color.New(color.FgMagenta)
var white = color.New(color.FgWhite)

// Green ...
func Green(s string, args ...interface{}) string {
	return green.Sprintf(s, args...)
}

// Red ...
func Red(s string, args ...interface{}) string {
	return red.Sprintf(s, args...)
}

// Cyan ...
func Cyan(s string, args ...interface{}) string {
	return cyan.Sprintf(s, args...)
}

// Yellow ...
func Yellow(s string, args ...interface{}) string {
	return yellow.Sprintf(s, args...)
}

// Magenta ...
func Magenta(s string, args ...interface{}) string {
	return magenta.Sprintf(s, args...)
}

// White ...
func White(s string, args ...interface{}) string {
	return white.Sprintf(s, args...)
}

// Info ...
func Info(s string, args ...interface{}) string {
	return cyan.Sprintf(s, args...)
}

// Panic ...
func Panic(s string, args ...interface{}) string {
	return red.Sprintf(s, args...)
}

// Error ...
func Error(s string, args ...interface{}) string {
	return red.Sprintf(s, args...)
}

// Warn ...
func Warn(s string, args ...interface{}) string {
	return yellow.Sprintf(s, args...)
}

// Debug ...
func Debug(s string, args ...interface{}) string {
	return white.Sprintf(s, args...)
}
