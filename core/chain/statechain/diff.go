package statechain

import (
	"encoding/json"
	"errors"

	"github.com/c3systems/c3-go/common/coder"
	"github.com/c3systems/c3-go/common/hashutil"
	"github.com/c3systems/c3-go/common/hexutil"
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
	tmp := BuildCoderFromDiff(d)
	bytes, err := tmp.Marshal()
	if err != nil {
		return nil, err
	}

	return coder.AppendCode(bytes), nil
}

// Deserialize ...
func (d *Diff) Deserialize(data []byte) error {
	if data == nil {
		return errors.New("nil bytes")
	}
	if d == nil {
		return errors.New("nil diff")
	}

	_, bytes, err := coder.StripCode(data)
	if err != nil {
		return err
	}

	props, err := BuildDiffPropsFromBytes(bytes)
	if err != nil {
		return err
	}

	d.props = *props

	return nil
}

// SerializeString ...
func (d Diff) SerializeString() (string, error) {
	bts, err := d.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeToString(bts), nil
}

// DeserializeString ...
func (d *Diff) DeserializeString(hexStr string) error {
	if d == nil {
		return ErrNilDiff
	}

	bts, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return d.Deserialize(bts)
}

// CalculateHash ...
func (d Diff) CalculateHash() (string, error) {
	bts, err := d.CalculateHashBytes()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeToString(bts), nil
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

	hashedBytes := hashutil.Hash(bytes)
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

// BuildCoderFromDiff ...
func BuildCoderFromDiff(d *Diff) *coder.Diff {
	tmp := &coder.Diff{
		Data: d.props.Data,
	}

	// note: is there a better way to handle nil with protobuff?
	if d.props.DiffHash != nil {
		tmp.DiffHash = *d.props.DiffHash
	}

	return tmp
}

// BuildDiffPropsFromBytes ...
func BuildDiffPropsFromBytes(data []byte) (*DiffProps, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	c, err := BuildBytesCoderFromBytes(data)
	if err != nil {
		return nil, err
	}

	return BuildDiffPropsFromCoder(c)
}

// BuildBytesCoderFromBytes ...
func BuildBytesCoderFromBytes(data []byte) (*coder.Diff, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	tmp := new(coder.Diff)
	if err := tmp.Unmarshal(data); err != nil {
		return nil, err
	}
	if tmp == nil {
		return nil, errors.New("nil output")
	}

	return tmp, nil
}

// BuildDiffPropsFromCoder ...
func BuildDiffPropsFromCoder(tmp *coder.Diff) (*DiffProps, error) {
	if tmp == nil {
		return nil, errors.New("nil coder")
	}

	props := &DiffProps{
		Data: tmp.Data,
	}
	if tmp.DiffHash != "" {
		s := tmp.DiffHash
		props.DiffHash = &s
	}

	return props, nil
}
