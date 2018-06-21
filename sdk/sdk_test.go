package sdk

import "testing"

func setItem(key, value string) error {
	return nil
}

func TestRegisterMethod(t *testing.T) {
	srv := New()

	err := srv.RegisterMethod("setItem", []string{"string", "string"}, setItem)
	if err != nil {
		t.Error(err)
	}
}
