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
	p    = Props{
		BlockHash:         &hash,
		BlockNumber:       "0x1",
		BlockTime:         "0x5",
		ImageHash:         "imageHash",
		TxsHash:           "txsHash",
		StatePrevDiffHash: "prevStateHash",
		StateCurrentHash:  "currentStateHash",
	}
	p1 = Props{
		BlockHash:         nil,
		BlockNumber:       "0x1",
		BlockTime:         "0x5",
		ImageHash:         "imageHash1",
		TxsHash:           "txsHash",
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
	expecteds := []Props{
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

		if expected.props.TxsHash != actual.props.TxsHash ||
			expected.props.StatePrevDiffHash != actual.props.StatePrevDiffHash ||
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

func TestCID(t *testing.T) {
	blocks := []*Block{
		New(&p),
		New(&p1),
	}

	expecteds := []string{
		"zdpuAmqW24ts9TgR7NhAcamvrMZBjMMW1H3d2ReP9wPvi9aT3",
		"zdpuB1sVfUhZMjDp9E3BD83mGRrM8Ki9YusZ3YQiWoifQyjgy",
	}

	for idx, block := range blocks {
		cid, err := block.CID("block")
		if err != nil {
			t.Fatalf("test %d failed; err getting cid: %v", idx+1, err)
		}
		if cid == nil {
			t.Errorf("test %d failed; expected non-null cid", idx+1)
		}

		if cid.String() != expecteds[idx] {
			t.Errorf("test %d failed; expected %s, received %s", idx+1, expecteds[idx], cid.String())
		}
	}
}

func TestHash(t *testing.T) {
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

		actual, err := block.Hash()
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
