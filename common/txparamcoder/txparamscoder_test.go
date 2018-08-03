package txparamcoder

import (
	"testing"
)

// TODO: table test

func TestEncodeMethodName(t *testing.T) {
	encoded := EncodeMethodName("myMethod")
	if encoded != "0xcbbf8a92fd90f20d804bf89baafba8eb0faef816def7231faa7d21eec7b65a6a" {
		t.Error("expected match")
	}
}

func TestEncodeParam(t *testing.T) {
	encoded := EncodeParam("hello")
	if encoded != "0x68656c6c6f" {
		t.Error("expected match")
	}
}

func TestEncodeParams(t *testing.T) {
	encoded := EncodeParams("hello", "world")
	if encoded[0] != "0x68656c6c6f" {
		t.Error("expected match")
	}

	if encoded[1] != "0x776f726c64" {
		t.Error("expected match")
	}
}

func TestToJSONArray(t *testing.T) {
	js := ToJSONArray("hello", "world")
	if string(js) != `["hello","world"]` {
		t.Error("expected match")
	}
}

func TestAppendJSONArrays(t *testing.T) {
	js := AppendJSONArrays(
		ToJSONArray("hello", "world"),
		ToJSONArray("foo", "bar"),
	)

	if string(js) != `[["hello","world"],["foo","bar"]]` {
		t.Error("expected match")
	}
}
