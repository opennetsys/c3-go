package txparamcoder

import (
	"encoding/json"
	"log"

	"github.com/c3systems/c3-go/common/hashing"
	"github.com/c3systems/c3-go/common/hexutil"
)

// EncodeMethodName ...
func EncodeMethodName(name string) string {
	return hashing.HashToHexString([]byte(name))
}

// EncodeParam ...
func EncodeParam(arg string) string {
	return hexutil.EncodeToString([]byte(arg))
}

// EncodeParams ...
func EncodeParams(args ...string) []string {
	var encoded []string
	for _, arg := range args {
		encoded = append(encoded, hexutil.EncodeToString([]byte(arg)))
	}

	return encoded
}

// ToJSONArray ...
func ToJSONArray(args ...string) []byte {
	js, err := json.Marshal(args)
	if err != nil {
		log.Fatal(err)
	}
	return js
}

// AppendJSONArrays ...
func AppendJSONArrays(args ...[]byte) []byte {
	var combined [][]string
	for _, arg := range args {
		var js []string
		err := json.Unmarshal(arg, &js)
		if err != nil {
			log.Fatal(err)
		}
		combined = append(combined, js)
	}
	js, err := json.Marshal(combined)
	if err != nil {
		log.Fatal(err)
	}
	return js
}
