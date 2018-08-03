package color

import (
	"log"
	"testing"
)

func TestGreen(t *testing.T) {
	log.Println(Green("hello %s", "bob"))
}
