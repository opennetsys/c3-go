package sandbox

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/common/txparamcoder"
)

// Docker file found in /example
var imageID = "QmUJtVGaAwLihH7RxPjRqTR3VxtjA4JAqxGTtBsuuRbTuR"

//var imageID = "0192547c385a"

func init() {
	// Makefile will set this env var
	if os.Getenv("IMAGEID") != "" {
		imageID = os.Getenv("IMAGEID")
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	sb := New(nil)
	if sb == nil {
		t.Error("expected instance")
	}
}

func TestPayload(t *testing.T) {
	t.Parallel()
	sb := New(nil)

	payload := txparamcoder.ToJSONArray(
		txparamcoder.EncodeMethodName("setItem"),
		txparamcoder.EncodeParam("foo"),
		txparamcoder.EncodeParam("bar"),
	)

	newState, err := sb.Play(&PlayConfig{
		ImageID: imageID,
		Payload: payload,
	})

	if err != nil {
		t.Error(err)
	}

	expectedStateMap := map[string]string{
		hexutil.EncodeToString([]byte("foo")): hexutil.EncodeToString([]byte("bar")),
	}

	expectedState, err := json.Marshal(expectedStateMap)

	if !reflect.DeepEqual(newState, expectedState) {
		t.Error("expected new state")
	}
}

func TestInitialState(t *testing.T) {
	t.Parallel()
	sb := New(nil)
	payload := txparamcoder.ToJSONArray(
		txparamcoder.EncodeMethodName("setItem"),
		txparamcoder.EncodeParam("foo"),
		txparamcoder.EncodeParam("bar"),
	)
	initialStateMap := map[string]string{
		hexutil.EncodeToString([]byte("hello")): hexutil.EncodeToString([]byte("world")),
	}

	initialState, err := json.Marshal(initialStateMap)
	if err != nil {
		t.Error(err)
	}

	newState, err := sb.Play(&PlayConfig{
		ImageID:      imageID,
		Payload:      payload,
		InitialState: initialState,
	})

	if err != nil {
		t.Error(err)
	}

	expectedStateMap := map[string]string{
		hexutil.EncodeToString([]byte("foo")):   hexutil.EncodeToString([]byte("bar")),
		hexutil.EncodeToString([]byte("hello")): hexutil.EncodeToString([]byte("world")),
	}

	expectedState, err := json.Marshal(expectedStateMap)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("expected new state; got %s", string(newState))
	}
}

func TestMultipleInputs(t *testing.T) {
	t.Parallel()
	sb := New(nil)
	payload := txparamcoder.AppendJSONArrays(
		txparamcoder.ToJSONArray(
			txparamcoder.EncodeMethodName("setItem"),
			txparamcoder.EncodeParam("foo"),
			txparamcoder.EncodeParam("bar"),
		),
		txparamcoder.ToJSONArray(
			txparamcoder.EncodeMethodName("setItem"),
			txparamcoder.EncodeParam("hello"),
			txparamcoder.EncodeParam("mars"),
		),
	)
	initialState := []byte(fmt.Sprintf(`{%q:%q}`, hexutil.EncodeToString([]byte("hello")), hex.EncodeToString([]byte("world"))))
	newState, err := sb.Play(&PlayConfig{
		ImageID:      imageID,
		Payload:      payload,
		InitialState: initialState,
	})

	if err != nil {
		t.Error(err)
	}

	expectedStateMap := map[string]string{
		hexutil.EncodeToString([]byte("foo")):   hexutil.EncodeToString([]byte("bar")),
		hexutil.EncodeToString([]byte("hello")): hexutil.EncodeToString([]byte("mars")),
	}

	expectedState, err := json.Marshal(expectedStateMap)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("expected new state; got %s", string(newState))
	}
}
