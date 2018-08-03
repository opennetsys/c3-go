package types

import "github.com/c3systems/merkletree"

// ChainObject ...
type ChainObject interface {
	Serialize() ([]byte, error)
	Deserialize(bytes []byte) error
	SerializeString() (string, error)
	DeserializeString(hexStr string) error
	CalculateHash() (string, error)
	CalculateHashBytes() ([]byte, error)
	Equals(other merkletree.Content) (bool, error)
	SetHash() error
}
