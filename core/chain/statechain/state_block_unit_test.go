// +build unit

package statechain

import (
	"reflect"
	"strings"
	"testing"
	"unicode"
)

var (
	hash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	p    = BlockProps{
		BlockHash:         &hash,
		BlockNumber:       "0x1",
		BlockTime:         "0x5",
		ImageHash:         "imageHash",
		StatePrevDiffHash: "prevStateHash",
		StateCurrentHash:  "currentStateHash",
	}
	p1 = BlockProps{
		BlockHash:         nil,
		BlockNumber:       "0x1",
		BlockTime:         "0x5",
		ImageHash:         "imageHash1",
		StatePrevDiffHash: "prevStateHash",
		StateCurrentHash:  "currentStateHash",
	}
	s = `{
		"blockHash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"blockNumber": "0x1",
		"blockTime": "0x5",
		"imageHash": "imageHash",
		"txsHash": "txsHash",
		"statePrevDiffHash": "prevStateHash",
		"stateCurrentHash": "currentStateHash"
	}`
	s1 = `{
		"blockNumber": "0x1",
		"blockTime": "0x5",
		"imageHash": "imageHash1",
		"txsHash": "txsHash",
		"statePrevDiffHash": "prevStateHash",
		"stateCurrentHash": "currentStateHash"
	}`
	//blockHash  = [32]byte{151, 199, 91, 159, 176, 205, 25, 22, 244, 36, 182, 228, 165, 18, 233, 115, 157, 92, 212, 219, 176, 103, 20, 69, 106, 156, 93, 253, 4, 235, 43, 127}
	//blockHash1 = [32]byte{56, 141, 178, 217, 91, 120, 234, 79, 136, 94, 56, 232, 202, 59, 161, 30, 98, 100, 183, 246, 207, 97, 234, 248, 87, 231, 101, 57, 232, 217, 119, 166}
	blockHash  = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	blockHash1 = "647d8ddefae7c4374de27c19fedeaef9d2ec6b72fa40f6b4014657a71e2bfc95"
)

func TestNew(t *testing.T) {
	t.Parallel()

	expecteds := []Block{
		Block{
			props: p,
		},
		Block{
			props: p1,
		},
	}

	actuals := []*Block{
		New(&p),
		New(&p1),
	}

	for idx, expected := range expecteds {
		if actuals[idx] == nil || !reflect.DeepEqual(expected, *actuals[idx]) {
			t.Errorf("test #%d failed; expected: %v; received: %v", idx+1, expected, actuals[idx])
		}
	}

}

func TestProps(t *testing.T) {
	t.Parallel()

	expecteds := []BlockProps{
		p,
		p1,
	}

	blocks := []*Block{
		New(&expecteds[0]),
		New(&expecteds[1]),
	}

	for idx, block := range blocks {
		if block == nil {
			t.Fatalf("test %d failed; expected non-nil block", idx+1)
		}

		if !reflect.DeepEqual(block.Props(), expecteds[idx]) {
			t.Errorf("test #%d failed; expected: %v; received: %v", idx+1, expecteds[idx], block.Props())
		}
	}
}

func TestSerialize(t *testing.T) {
	// TODO: fix!
	t.Skip()
	expecteds := []string{
		s,
		s1,
	}

	blocks := []*Block{
		New(&p),
		New(&p1),
	}

	for idx, block := range blocks {
		if block == nil {
			t.Fatalf("test %d failed; expected non-nil block", idx+1)
		}

		actual, err := block.Serialize()
		if err != nil {
			t.Fatalf("test %d failed; err serializing block: %v", idx+1, err)
		}

		expected := []byte(removeWhiteSpace(expecteds[idx]))
		actual = []byte(removeWhiteSpace(string(actual)))
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("test %d failed; expected: %s; received: %s", idx+1, string(expected), string(actual))
		}
	}
}

func TestDeserialize(t *testing.T) {
	// TODO: fix!
	t.Skip()
	expecteds := []Block{
		Block{
			props: p,
		},
		Block{
			props: p1,
		},
	}

	inputs := []string{
		s,
		s1,
	}

	for idx, expected := range expecteds {
		var actual Block
		if err := actual.Deserialize([]byte(inputs[idx])); err != nil {
			t.Fatalf("test %d failed; err parsing from bytes: %v", idx+1, err)
		}

		if expected.props.StatePrevDiffHash != actual.props.StatePrevDiffHash ||
			expected.props.StateCurrentHash != actual.props.StateCurrentHash ||
			expected.props.ImageHash != actual.props.ImageHash ||
			expected.props.BlockNumber != actual.props.BlockNumber ||
			expected.props.BlockTime != actual.props.BlockTime {
			t.Errorf("test %d failed; expected: %v; received: %v", idx+1, expected, actual)
		}

		if expected.props.BlockHash == nil && actual.props.BlockHash != nil ||
			expected.props.BlockHash != nil && actual.props.BlockHash == nil {
			t.Errorf("test %d failed; expected hash %v, reecived hash: %v", idx+1, expected.props.BlockHash, actual.props.BlockHash)
		}

		if (expected.props.BlockHash != nil && actual.props.BlockHash != nil) && *expected.props.BlockHash != *actual.props.BlockHash {
			t.Errorf("test %d failed; expected hash %s, reecived hash: %s", idx+1, *expected.props.BlockHash, *actual.props.BlockHash)
		}
	}
}

func TestSerializeDeserializeStatechainBlock(t *testing.T) {
	t.Parallel()

	blocks := []*Block{
		New(&p),
		New(&p1),
	}

	for idx, input := range blocks {
		bytes, err := input.Serialize()
		if err != nil {
			t.Errorf("test %d failed to serialize\n%v", idx+1, err)
		}

		b := new(Block)
		if err := b.Deserialize(bytes); err != nil {
			t.Errorf("test %d failed to deserialize\n%v", idx+1, err)
		}
		if b == nil {
			t.Fatalf("test %d failed; block is nil", idx+1)
		}

		if !reflect.DeepEqual(input, b) {
			t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, *input, *b)
		}
	}
}

func TestSerializeDeserializeStringStatechainBlock(t *testing.T) {
	t.Parallel()

	blocks := []*Block{
		New(&p),
		New(&p1),
	}

	for idx, input := range blocks {
		str, err := input.SerializeString()
		if err != nil {
			t.Errorf("test %d failed to serialize\n%v", idx+1, err)
		}

		b := new(Block)
		if err := b.DeserializeString(str); err != nil {
			t.Errorf("test %d failed to deserialize\n%v", idx+1, err)
		}
		if b == nil {
			t.Fatalf("test %d failed; block is nil", idx+1)
		}

		if !reflect.DeepEqual(input, b) {
			t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, *input, *b)
		}
	}
}

func TestHash(t *testing.T) {
	// TODO: fix!
	t.Skip()

	expecteds := []string{
		blockHash,
		blockHash1,
	}

	blocks := []*Block{
		New(&p),
		New(&p1),
	}

	for idx, block := range blocks {
		if block == nil {
			t.Fatalf("test %d failed; expected non-nil block", idx+1)
		}

		actual, err := block.CalculateHash()
		if err != nil {
			t.Fatalf("test %d failed; err serializing block: %v", idx+1, err)
		}

		if !reflect.DeepEqual(expecteds[idx], actual) {
			t.Errorf("test %d failed; expected: %v; received: %v", idx+1, expecteds[idx], actual)
		}
	}

}

func removeWhiteSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
