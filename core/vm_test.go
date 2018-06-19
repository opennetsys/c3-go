package vm

import "testing"

func TestNew(t *testing.T) {
	vm := New()
	_ = vm
}

func TestListImages(t *testing.T) {
	vm := New()
	vm.ListImages()
}
