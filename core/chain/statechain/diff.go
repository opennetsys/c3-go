package statechain

import (
	"encoding/json"

	"github.com/c3systems/c3/common/hashing"
	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/merkletree"
)

// NewDiff ...
func NewDiff(props *DiffProps) *Diff {
	if props == nil {
		return &Diff{}
	}

	return &Diff{
		props: *props,
	}
}

// Props ...
func (d Diff) Props() DiffProps {
	return d.props
}

// Serialize ...
func (d *Diff) Serialize() ([]byte, error) {
	return d.MarshalJSON()
}

// Deserialize ...
func (d *Diff) Deserialize(data []byte) error {
	return d.UnmarshalJSON(data)
}

// SerializeString ...
func (d Diff) SerializeString() (string, error) {
	bytes, err := d.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(bytes)), nil
}

// DeserializeString ...
func (d *Diff) DeserializeString(hexStr string) error {
	if d == nil {
		return ErrNilDiff
	}

	str, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return d.Deserialize([]byte(str))
}

// CalculateHash ...
func (d Diff) CalculateHash() (string, error) {
	bytes, err := d.CalculateHashBytes()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(bytes)), nil
}

// CalculateHashBytes ...
func (d Diff) CalculateHashBytes() ([]byte, error) {
	tmpDiff := Diff{
		props: DiffProps{
			Data: d.props.Data,
		},
	}

	bytes, err := tmpDiff.Serialize()
	if err != nil {
		return nil, err
	}

	hashedBytes := hashing.Hash(bytes)
	return hashedBytes[:], nil
}

// Equals ...
func (d Diff) Equals(other merkletree.Content) (bool, error) {
	dHash, err := d.CalculateHashBytes()
	if err != nil {
		return false, err
	}

	oHash, err := other.CalculateHashBytes()
	if err != nil {
		return false, err
	}

	return string(dHash) == string(oHash), nil
}

// SetHash ...
func (d *Diff) SetHash() error {
	if d == nil {
		return ErrNilDiff
	}

	hash, err := d.CalculateHash()
	if err != nil {
		return err
	}

	d.props.DiffHash = &hash

	return nil
}

// MarshalJSON ...
func (d *Diff) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.props)
}

// UnmarshalJSON ...
func (d *Diff) UnmarshalJSON(data []byte) error {
	var props DiffProps
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	d.props = props

	return nil
}
