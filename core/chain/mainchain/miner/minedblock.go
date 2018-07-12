package miner

import (
	"encoding/json"

	"github.com/c3systems/c3/common/hexutil"
)

// Serialize ...
func (m *MinedBlock) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// Deserialize ...
func (m *MinedBlock) Deserialize(data []byte) error {
	return json.Unmarshal(data, m)
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
