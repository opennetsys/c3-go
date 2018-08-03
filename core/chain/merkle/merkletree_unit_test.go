// +build unit

package merkle

import (
	"reflect"
	"testing"
)

var (
	hash   = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	props1 = &TreeProps{
		MerkleTreeRootHash: nil,
		Kind:               "kind1",
		Hashes: []string{
			"1",
			"2",
			"3",
		},
	}
	props2 = &TreeProps{
		MerkleTreeRootHash: &hash,
		Kind:               "kind2",
		Hashes: []string{
			"4",
			"5",
			"6",
		},
	}
)

func TestSerializeDeserialize(t *testing.T) {
	t.Parallel()

	t1, err := New(props1)
	if err != nil {
		t.Error(err)
	}
	t2, err := New(props2)
	if err != nil {
		t.Error(err)
	}

	inputs := []*Tree{
		t1,
		t2,
	}

	for idx, input := range inputs {
		bytes, err := input.Serialize()
		if err != nil {
			t.Errorf("test %d failed serialization\n%v", idx+1, err)
		}

		tr := new(Tree)
		if err := tr.Deserialize(bytes); err != nil {
			t.Errorf("test %d failed deserialization\n%v", idx+1, err)
		}
		if tr == nil {
			t.Errorf("test %d failed\nnil tree", idx+1)
		}

		if input.props.Kind != tr.props.Kind {
			t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, *input, *tr)
		}

		if !reflect.DeepEqual(input.props.MerkleTreeRootHash, tr.props.MerkleTreeRootHash) {
			t.Errorf("test %d failed\nexpected: %v\n received: %v", idx+1, input.props.MerkleTreeRootHash, tr.props.MerkleTreeRootHash)
		}

		if !reflect.DeepEqual(input.props.Hashes, tr.props.Hashes) {
			t.Errorf("test %d failed\nexpected hashes: %v\n received hashes: %v", idx+1, input.props.Hashes, tr.props.Hashes)
		}
	}
}

func TestSerializeDeserializeString(t *testing.T) {
	t.Parallel()

	t1, err := New(props1)
	if err != nil {
		t.Error(err)
	}
	t2, err := New(props2)
	if err != nil {
		t.Error(err)
	}

	inputs := []*Tree{
		t1,
		t2,
	}

	for idx, input := range inputs {
		str, err := input.SerializeString()
		if err != nil {
			t.Errorf("test %d failed serialization\n%v", idx+1, err)
		}

		tr := new(Tree)
		if err := tr.DeserializeString(str); err != nil {
			t.Errorf("test %d failed deserialization\n%v", idx+1, err)
		}
		if tr == nil {
			t.Errorf("test %d failed\nnil t", idx+1)
		}

		if input.props.Kind != tr.props.Kind {
			t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, *input, *tr)
		}

		if !reflect.DeepEqual(input.props.MerkleTreeRootHash, tr.props.MerkleTreeRootHash) {
			t.Errorf("test %d failed\nexpected: %v\n received: %v", idx+1, input.props.MerkleTreeRootHash, tr.props.MerkleTreeRootHash)
		}

		if !reflect.DeepEqual(input.props.Hashes, tr.props.Hashes) {
			t.Errorf("test %d failed\nexpected hashes: %v\n received hashes: %v", idx+1, input.props.Hashes, tr.props.Hashes)
		}
	}
}
