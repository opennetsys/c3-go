package coder

import "github.com/ethereum/go-ethereum/rlp"

// Serialize ...
func Serialize(v interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(v)
}

// Deserialize ...
// note: v must be a pointer
func Deserialize(data []byte, v interface{}) error {
	if err := rlp.DecodeBytes(data, v); err != nil {
		return err
	}

	return nil
}
