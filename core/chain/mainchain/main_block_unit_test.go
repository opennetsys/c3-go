// +build unit

package mainchain

import (
	"reflect"
	"strings"
	"testing"
	"unicode"

	log "github.com/sirupsen/logrus"
)

var (
	hash = "0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	bSig = MinerSig{
		R: "0x1",
		S: "0x1",
	}
	p = Props{
		BlockHash:             &hash,
		BlockNumber:           "0x1",
		BlockTime:             "0x5",
		ImageHash:             "0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		StateBlocksMerkleHash: "0xstateBlocksHash",
		PrevBlockHash:         "0xprevBlockHash",
		Nonce:                 "0x1",
		Difficulty:            "0x1",
		MinerAddress:          "0x123",
		MinerSig:              &bSig,
	}
	p1 = Props{
		BlockHash:             nil,
		BlockNumber:           "0x2",
		BlockTime:             "0x5",
		ImageHash:             "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		StateBlocksMerkleHash: "stateBlocksHash",
		PrevBlockHash:         "prevBlockHash",
		Nonce:                 "0x1",
		Difficulty:            "0x1",
		MinerAddress:          "",
		MinerSig:              nil,
	}
	s = `{
		"blockHash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"blockNumber": "0x1",
		"blockTime": "0x5",
		"imageHash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"stateBlocksHash": "stateBlocksHash",
		"prevBlockHash": "prevBlockHash",
		"nonce": "0x1",
		"difficulty": "0x1"
	}`
	s1 = `{
		"blockNumber": "0x2",
		"blockTime": "0x5",
		"imageHash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"stateBlocksHash": "stateBlocksHash",
		"prevBlockHash": "prevBlockHash",
		"nonce": "0x1",
		"difficulty": "0x1"
	}`
	//blockHash  = [32]byte{151, 199, 91, 159, 176, 205, 25, 22, 244, 36, 182, 228, 165, 18, 233, 115, 157, 92, 212, 219, 176, 103, 20, 69, 106, 156, 93, 253, 4, 235, 43, 127}
	//blockHash1 = [32]byte{56, 141, 178, 217, 91, 120, 234, 79, 136, 94, 56, 232, 202, 59, 161, 30, 98, 100, 183, 246, 207, 97, 234, 248, 87, 231, 101, 57, 232, 217, 119, 166}
	blockHash  = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	blockHash1 = "40a7aa8078f0487485ee7ab3fa497e5b8302eb96c87c1474a8a8eda591503599"
)

func TestNew(t *testing.T) {
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

	actuals := []*Block{
		New(&p),
		New(&p1),
	}

	for idx, expected := range expecteds {
		if actuals[idx] == nil || !reflect.DeepEqual(expected, *actuals[idx]) {
			t.Errorf("test #%d failed; expected: %v; received: %v", idx+1, expected, actuals[idx])
		}
	}

	// A nil props and a props with the incorrect imageHash should return the correct image hash
	inputs := []*Props{
		nil,
		&Props{
			ImageHash: "foo",
		},
	}

	for idx, input := range inputs {
		actual := New(input)
		if actual == nil || actual.Props().ImageHash != ImageHash {
			t.Errorf("test #%d failed; expected: %v; received: %v", idx+1, ImageHash, actual.Props().ImageHash)
		}
	}
}

func TestProps(t *testing.T) {
	t.Parallel()

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
			t.Errorf("test %d failed; expected non-nil block", idx+1)
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
			t.Errorf("test %d failed; expected non-nil block", idx+1)
		}

		actual, err := block.Serialize()
		if err != nil {
			t.Errorf("test %d failed; err serializing block: %v", idx+1, err)
		}
		t.Log(actual)

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
			t.Errorf("test %d failed; err parsing from bytes: %v", idx+1, err)
		}

		if expected.props.StateBlocksMerkleHash != actual.props.StateBlocksMerkleHash ||
			expected.props.PrevBlockHash != actual.props.PrevBlockHash ||
			expected.props.BlockNumber != actual.props.BlockNumber ||
			expected.props.ImageHash != actual.props.ImageHash ||
			expected.props.Nonce != actual.props.Nonce ||
			expected.props.Difficulty != actual.props.Difficulty ||
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
			t.Errorf("test %d failed; expected non-nil block", idx+1)
		}

		actual, err := block.CalculateHash()
		if err != nil {
			t.Errorf("test %d failed; err serializing block: %v", idx+1, err)
		}

		if !reflect.DeepEqual(expecteds[idx], actual) {
			t.Errorf("test %d failed; expected: %v; received: %v", idx+1, expecteds[idx], actual)
		}
	}

}

func TestGenesisBlockHash(t *testing.T) {
	t.Parallel()

	hash, err := GenesisBlock.CalculateHash()
	if err != nil {
		t.Error(err)
	}

	if hash != GenesisBlockHash {
		log.Errorf("received %s\nexpected %s", hash, GenesisBlockHash)
	}
}

func TestSerializeDeserialize(t *testing.T) {
	t.Parallel()

	b := &Block{
		props: p,
	}
	b1 := &Block{
		props: p1,
	}

	inputs := []*Block{b, b1}
	for idx, in := range inputs {
		bytes, err := in.Serialize()
		if err != nil {
			t.Errorf("test %d failed\nerr serializing: %v", idx+1, err)
		}

		tmpBlock := new(Block)
		if err := tmpBlock.Deserialize(bytes); err != nil {
			t.Errorf("test %d failed\nerr deserializing: %v", idx+1, err)
		}

		if !reflect.DeepEqual(in, tmpBlock) {
			if tmpBlock.props.BlockHash != nil {
				t.Logf("block hash %s", *tmpBlock.props.BlockHash)
			}
			t.Errorf("test #%d faild\nexpected: %v\nreceived: %v\n", idx+1, *in, *tmpBlock)
		}
	}
}

func TestSerializeDeserializeString(t *testing.T) {
	t.Parallel()

	b := &Block{
		props: p,
	}
	b1 := &Block{
		props: p1,
	}

	inputs := []*Block{b, b1}
	for idx, in := range inputs {
		str, err := in.SerializeString()
		if err != nil {
			t.Errorf("test %d failed\nerr serializing: %v", idx+1, err)
		}

		tmpBlock := new(Block)
		if err := tmpBlock.DeserializeString(str); err != nil {
			t.Errorf("test %d failed\nerr deserializing: %v", idx+1, err)
		}

		if !reflect.DeepEqual(in, tmpBlock) {
			if tmpBlock.props.BlockHash != nil {
				t.Logf("block hash %s", *tmpBlock.props.BlockHash)
			}
			t.Errorf("test #%d faild\nexpected: %v\nreceived: %v\n", idx+1, *in, *tmpBlock)
		}
	}
}

func TestDeepEqual(t *testing.T) {
	t.Parallel()

	b := &Block{
		props: p,
	}
	b1 := &Block{
		props: p,
	}

	if !reflect.DeepEqual(b, b1) {
		t.Error("not equal")
	}
}

func TestSerializeDeserializeSig(t *testing.T) {
	// TODO: fix!
	t.Skip()

	//bytes, err := coder.Serialize(bSig)
	//if err != nil {
	//t.Error(err)
	//}

	//sig := new(BlockSig)
	//if err := coder.Deserialize(bytes, &sig); err != nil {
	//t.Error(err)
	//}
	//if sig == nil {
	//t.Error("nil sig")
	//}

	//if !reflect.DeepEqual(*sig, bSig) {
	//t.Errorf("expected %v\nreceived %v", bSig, *sig)
	//}
}

func removeWhiteSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
