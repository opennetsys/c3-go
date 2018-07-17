package statechain

import (
	"reflect"
	"testing"
)

var (
	diffHash   = "0xHash"
	diffProps1 = &DiffProps{
		DiffHash: nil,
		Data:     "0x1",
	}
	diffProps2 = &DiffProps{
		DiffHash: &diffHash,
		Data:     "0x1",
	}
)

func TestSerializeDeserializeDiff(t *testing.T) {
	t.Parallel()

	inputs := []*Diff{
		NewDiff(diffProps1),
		NewDiff(diffProps2),
	}

	for idx, input := range inputs {
		bytes, err := input.Serialize()
		if err != nil {
			t.Errorf("test %d failed serialization\n%v", idx+1, err)
		}

		d := new(Diff)
		if err := d.Deserialize(bytes); err != nil {
			t.Errorf("test %d failed deserialization\n%v", idx+1, err)
		}
		if d == nil {
			t.Errorf("test %d failed\nnil diff", idx+1)
		}

		if input.props.Data != d.props.Data {
			t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, *input, *d)
		}

		if !reflect.DeepEqual(input.props.DiffHash, d.props.DiffHash) {
			t.Errorf("test %d failed\nexpected txHash: %v\n received txHash: %v", idx+1, input.props.DiffHash, d.props.DiffHash)
		}
	}
}

func TestSerializeDeserializeStringDiff(t *testing.T) {
	t.Parallel()

	inputs := []*Diff{
		NewDiff(diffProps1),
		NewDiff(diffProps2),
	}

	for idx, input := range inputs {
		str, err := input.SerializeString()
		if err != nil {
			t.Errorf("test %d failed serialization\n%v", idx+1, err)
		}

		d := new(Diff)
		if err := d.DeserializeString(str); err != nil {
			t.Errorf("test %d failed deserialization\n%v", idx+1, err)
		}
		if d == nil {
			t.Errorf("test %d failed\nnil diff", idx+1)
		}

		if input.props.Data != d.props.Data {
			t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, *input, *d)
		}

		if !reflect.DeepEqual(input.props.DiffHash, d.props.DiffHash) {
			t.Errorf("test %d failed\nexpected txHash: %v\n received txHash: %v", idx+1, input.props.DiffHash, d.props.DiffHash)
		}
	}
}
