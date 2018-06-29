// +build unit

package stateblock

import (
	"reflect"
	"strings"
	"testing"
	"unicode"
)

var (
	hash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	p    = Props{
		TxsHash:              "txsHash",
		ImageHash:            "imageHash",
		StatePrevDiffHash:    "prevStateHash",
		StateGenesisDiffHash: "genesisStateHash",
		StateCurrentHash:     "currentStateHash",
		BlockNumber:          "0x1",
		TimeStamp:            "0x5",
		Nonce:                "0x1",
		Hash:                 &hash,
	}
	p1 = Props{
		TxsHash:              "txsHash",
		ImageHash:            "imageHash1",
		StatePrevDiffHash:    "prevStateHash",
		StateGenesisDiffHash: "genesisStateHash",
		StateCurrentHash:     "currentStateHash",
		BlockNumber:          "0x1",
		TimeStamp:            "0x5",
		Nonce:                "0x1",
		Hash:                 nil,
	}
	s = `{
		"txsHash": "txsHash",
		"imageHash": "imageHash",
		"statePrevDiffHash": "prevStateHash",
		"stateGenesisDiffHash": "genesisStateHash",
		"stateCurrentHash": "currentStateHash",
		"blockNumber": "0x1",
		"timeStamp": "0x5",
		"nonce": "0x1",
		"hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	}`
	s1 = `{
		"txsHash": "txsHash",
		"imageHash": "imageHash1",
		"statePrevDiffHash": "prevStateHash",
		"stateGenesisDiffHash": "genesisStateHash",
		"stateCurrentHash": "currentStateHash",
		"blockNumber": "0x1",
		"timeStamp": "0x5",
		"nonce": "0x1"
	}`
	//blockHash  = [32]byte{151, 199, 91, 159, 176, 205, 25, 22, 244, 36, 182, 228, 165, 18, 233, 115, 157, 92, 212, 219, 176, 103, 20, 69, 106, 156, 93, 253, 4, 235, 43, 127}
	//blockHash1 = [32]byte{56, 141, 178, 217, 91, 120, 234, 79, 136, 94, 56, 232, 202, 59, 161, 30, 98, 100, 183, 246, 207, 97, 234, 248, 87, 231, 101, 57, 232, 217, 119, 166}
	blockHash  = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	blockHash1 = "91a218ec33e4cf7371f98ce898bc98ed5452e0df7a6a8a6db48d7a5fb6f3c53e"
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

func TestFromBytes(t *testing.T) {
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
		if err := actual.FromBytes([]byte(inputs[idx])); err != nil {
			t.Fatalf("test %d failed; err parsing from bytes: %v", idx+1, err)
		}

		if expected.props.TxsHash != actual.props.TxsHash ||
			expected.props.StatePrevDiffHash != actual.props.StatePrevDiffHash ||
			expected.props.StateGenesisDiffHash != actual.props.StateGenesisDiffHash ||
			expected.props.StateCurrentHash != actual.props.StateCurrentHash ||
			expected.props.BlockNumber != actual.props.BlockNumber ||
			expected.props.TimeStamp != actual.props.TimeStamp ||
			expected.props.Nonce != actual.props.Nonce {
			t.Errorf("test %d failed; expected: %v; received: %v", idx+1, expected, actual)
		}

		if expected.props.Hash == nil && actual.props.Hash != nil ||
			expected.props.Hash != nil && actual.props.Hash == nil {
			t.Errorf("test %d failed; expected hash %v, reecived hash: %v", idx+1, expected.props.Hash, actual.props.Hash)
		}

		if (expected.props.Hash != nil && actual.props.Hash != nil) && *expected.props.Hash != *actual.props.Hash {
			t.Errorf("test %d failed; expected hash %s, reecived hash: %s", idx+1, *expected.props.Hash, *actual.props.Hash)
		}
	}
}

func TestCID(t *testing.T) {
	blocks := []*Block{
		New(&p),
		New(&p1),
	}

	expecteds := []string{
		"zdpuAo94u8Pes5nv3Vi7xoMCuzzwWve93pqYfpAPMYtJP7QaS",
		"zdpuB2BczX1xru7DE3AEFs53df3pe4SfYoFyJ1ZDsXRcMToxz",
	}

	for idx, block := range blocks {
		cid, err := block.CID()
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
