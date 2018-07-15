package sdk

import (
	"fmt"
	"testing"
)

func TestRegisterMethod(t *testing.T) {
	t.Parallel()
	c3 := NewC3()

	err := c3.RegisterMethod("setItem", []string{"string", "string"}, func(key, value string) error {
		fmt.Println("test setItem called with:", key, value)
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	t.Parallel()
	c3 := NewC3()

	c3.State().Set("foo", "bar")
	value := c3.State().Get("foo")
	if value != "bar" {
		t.Error("expected match")
	}
}

func TestState(t *testing.T) {
	t.Parallel()
	c3 := NewC3()

	err := c3.RegisterMethod("setItem", []string{"string", "string"}, func(key, value string) error {
		fmt.Println("test setItem called with:", key, value)
		c3.State().Set(key, value)
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	err = c3.Process([]byte(`[["setItem", "foo", "bar"],["setItem", "hello", "world"]]`))
	if err != nil {
		t.Error(err)
	}

	value := c3.State().Get("foo")
	if value != "bar" {
		t.Errorf("expected match; got %s", value)
	}

	value = c3.State().Get("hello")
	if value != "world" {
		t.Errorf("expected match; got %s", value)
	}
}
