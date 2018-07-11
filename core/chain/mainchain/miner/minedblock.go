package miner

import (
	"bytes"
	"encoding/gob"

	"github.com/c3systems/c3/common/hexutil"
)

// Serialize ...
func (m *MinedBlock) Serialize() ([]byte, error) {
	data := new(bytes.Buffer)
	err := gob.NewEncoder(data).Encode(m)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

// Deserialize ...
func (m *MinedBlock) Deserialize(data []byte) error {
	if m == nil {
		return ErrNilBlock
	}

	byts := bytes.NewBuffer(data)
	return gob.NewDecoder(byts).Decode(m)
}

// SerializeString ...
func (m *MinedBlock) SerializeString() (string, error) {
	bytes, err := m.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(bytes)), nil
}

// DeserializeString ...
func (m *MinedBlock) DeserializeString(hexStr string) error {
	if m == nil {
		return ErrNilBlock
	}

	str, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return m.Deserialize([]byte(str))
}
