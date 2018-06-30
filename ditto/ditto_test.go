package ditto

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	srv := New(&Config{})
	_ = srv
}

func TestUploadImage(t *testing.T) {
	srv := New(&Config{})
	filepath := "./test_data/hello-world.tar"
	reader, err := os.Open(filepath)
	if err != nil {
		t.Error(err)
	}
	err = srv.UploadImage(reader)
	if err != nil {
		t.Error(err)
	}
}
