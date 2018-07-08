package c3

import "testing"

func setItem(key, value string) error {
	return nil
}

func TestRegisterMethod(t *testing.T) {
	c3 := NewC3()

	err := c3.RegisterMethod("setItem", []string{"string", "string"}, setItem)
	if err != nil {
		t.Error(err)
	}
}
