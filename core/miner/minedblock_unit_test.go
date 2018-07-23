// +build unit

package miner

import (
	"reflect"
	"testing"

	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/merkle"
	"github.com/c3systems/c3-go/core/chain/statechain"
)

var (
	b1                 = &MinedBlock{}
	mainchainBlockHash = "0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	mainchainBlockSig  = &mainchain.MinerSig{
		R: "0x1",
		S: "0x1",
	}
	mainchainBlockProps1 = &mainchain.Props{
		BlockHash:             &mainchainBlockHash,
		BlockNumber:           "0x1",
		BlockTime:             "0x5",
		ImageHash:             "0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		StateBlocksMerkleHash: "0xstateBlocksHash",
		PrevBlockHash:         "0xprevBlockHash",
		Nonce:                 "0x1",
		Difficulty:            "0x1",
		MinerAddress:          "0x123",
		MinerSig:              mainchainBlockSig,
	}
	mainchainBlockProps2 = &mainchain.Props{
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
	merkleTreeHash   = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	merkleTreeProps1 = &merkle.TreeProps{
		MerkleTreeRootHash: nil,
		Kind:               "kind1",
		Hashes: []string{
			"1",
			"2",
			"3",
		},
	}
	merkleTreeProps2 = &merkle.TreeProps{
		MerkleTreeRootHash: &merkleTreeHash,
		Kind:               "kind2",
		Hashes: []string{
			"4",
			"5",
			"6",
		},
	}
	diffHash   = "0xHash"
	diffProps1 = &statechain.DiffProps{
		DiffHash: nil,
		Data:     "0x1",
	}
	diffProps2 = &statechain.DiffProps{
		DiffHash: &diffHash,
		Data:     "0x1",
	}
	statechainBlockHash   = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	statechainBlockProps1 = &statechain.BlockProps{
		BlockHash:         &merkleTreeHash,
		BlockNumber:       "0x1",
		BlockTime:         "0x5",
		ImageHash:         "imageHash",
		StatePrevDiffHash: "prevStateHash",
		StateCurrentHash:  "currentStateHash",
	}
	statechainBlockProps2 = &statechain.BlockProps{
		BlockHash:         nil,
		BlockNumber:       "0x1",
		BlockTime:         "0x5",
		ImageHash:         "imageHash1",
		StatePrevDiffHash: "prevStateHash",
		StateCurrentHash:  "currentStateHash",
	}
	txPayload = []byte(`{"foo":"bar"}`)
	txHash    = "0xHash"
	txSig     = &statechain.TxSig{
		R: " 0x0",
		S: "0x1",
	}
	txProps1 = &statechain.TransactionProps{
		TxHash:    nil,
		ImageHash: "0x1",
		Method:    "0x2",
		Payload:   nil,
		From:      "0x3",
		Sig:       nil,
	}
	txProps2 = &statechain.TransactionProps{
		TxHash:    &txHash,
		ImageHash: "0x1",
		Method:    "0x2",
		Payload:   txPayload,
		From:      "0x3",
		Sig:       txSig,
	}
)

func TestSerializeDeserialize(t *testing.T) {
	t.Parallel()

	inputs, err := buildMinedBlockInputs()
	if err != nil {
		t.Fatal(err)
	}

	for idx, input := range inputs {
		bytes, err := input.Serialize()
		if err != nil {
			t.Fatalf("test %d faild serialiation\n%v", idx+1, err)
		}

		mined := new(MinedBlock)
		if err := mined.Deserialize(bytes); err != nil {
			t.Fatalf("test %d failed deserialization\n%v", idx+1, err)
		}

		isMinedBlockEqual(t, idx, input, mined)
	}
}

func TestSerializeDeserializeString(t *testing.T) {
	t.Parallel()

	inputs, err := buildMinedBlockInputs()
	if err != nil {
		t.Fatal(err)
	}

	for idx, input := range inputs {
		str, err := input.SerializeString()
		if err != nil {
			t.Fatalf("test %d faild serialiation\n%v", idx+1, err)
		}

		mined := new(MinedBlock)
		if err := mined.DeserializeString(str); err != nil {
			t.Fatalf("test %d failed deserialization\n%v", idx+1, err)
		}

		isMinedBlockEqual(t, idx, input, mined)
	}
}

func buildMinedBlockInputs() ([]*MinedBlock, error) {
	t1 := statechain.NewTransaction(txProps1)
	t2 := statechain.NewTransaction(txProps2)

	sb1 := statechain.New(statechainBlockProps1)
	sb2 := statechain.New(statechainBlockProps2)

	d1 := statechain.NewDiff(diffProps1)
	d2 := statechain.NewDiff(diffProps2)

	tr1, err := merkle.New(merkleTreeProps1)
	if err != nil {
		return nil, err
	}
	tr2, err := merkle.New(merkleTreeProps2)
	if err != nil {
		return nil, err
	}

	mb1 := mainchain.New(mainchainBlockProps1)
	mb2 := mainchain.New(mainchainBlockProps2)

	return []*MinedBlock{
		b1,
		&MinedBlock{
			NextBlock:     mb1,
			PreviousBlock: mb2,
			StatechainBlocksMap: map[string]*statechain.Block{
				"foo": sb1,
				"bar": sb2,
			},
			TransactionsMap: map[string]*statechain.Transaction{
				"foo": t1,
				"bar": t2,
			},
			DiffsMap: map[string]*statechain.Diff{
				"foo": d1,
				"bar": d2,
			},
			MerkleTreesMap: map[string]*merkle.Tree{
				"foo": tr1,
				"bar": tr2,
			},
		},
	}, nil
}

func isMinedBlockEqual(t *testing.T, idx int, input, mined *MinedBlock) {
	if !reflect.DeepEqual(input.NextBlock, mined.NextBlock) {
		t.Errorf("test %d failed\nexpected next block: %v\nreceived next block: %v", idx+1, input.NextBlock, mined.NextBlock)
	}

	if !reflect.DeepEqual(input.PreviousBlock, mined.PreviousBlock) {
		t.Errorf("test %d failed\nexpected previous block: %v\nreceived previous block: %v", idx+1, input.PreviousBlock, mined.PreviousBlock)
	}

	if len(input.MerkleTreesMap) == len(mined.MerkleTreesMap) {
		for k, v := range input.MerkleTreesMap {
			v1, ok := mined.MerkleTreesMap[k]
			if !ok {
				t.Errorf("test %d failed\n merkle tree maps key %s not present", idx+1, k)
			}

			if v.Props().Kind != v1.Props().Kind {
				t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, v.Props(), v1.Props())
			}

			if !reflect.DeepEqual(v.Props().MerkleTreeRootHash, v1.Props().MerkleTreeRootHash) {
				t.Errorf("test %d failed\nexpected: %v\n received: %v", idx+1, v.Props().MerkleTreeRootHash, v1.Props().MerkleTreeRootHash)
			}

			if !reflect.DeepEqual(v.Props().Hashes, v1.Props().Hashes) {
				t.Errorf("test %d failed\nexpected hashes: %v\n received hashes: %v", idx+1, v.Props().Hashes, v1.Props().Hashes)
			}
		}
	}

	if len(input.DiffsMap) == len(mined.DiffsMap) {
		for k, v := range input.DiffsMap {
			v1, ok := mined.DiffsMap[k]
			if !ok {
				t.Errorf("test %d failed\n diff maps key %s not present", idx+1, k)
			}

			if v.Props().Data != v1.Props().Data {
				t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, v.Props(), v1.Props())
			}

			if !reflect.DeepEqual(v.Props().DiffHash, v1.Props().DiffHash) {
				t.Errorf("test %d failed\nexpected txHash: %v\n received txHash: %v", idx+1, v.Props().DiffHash, v1.Props().DiffHash)
			}
		}
	}

	if len(input.StatechainBlocksMap) == len(mined.StatechainBlocksMap) {
		for k, v := range input.StatechainBlocksMap {
			v1, ok := mined.StatechainBlocksMap[k]
			if !ok {
				t.Errorf("test %d failed\n state blocks maps key %s not present", idx+1, k)
			}

			if !reflect.DeepEqual(v.Props(), v1.Props()) {
				t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, v.Props(), v1.Props())
			}
		}
	}

	if len(input.TransactionsMap) == len(mined.TransactionsMap) {
		for k, v := range input.TransactionsMap {
			v1, ok := mined.TransactionsMap[k]
			if !ok {
				t.Errorf("test %d failed\n diff maps key %s not present", idx+1, k)
			}

			if v.Props().ImageHash != v1.Props().ImageHash || v.Props().Method != v1.Props().Method || v.Props().From != v1.Props().From {
				t.Errorf("test %d failed\nexpected: %v\nreceived: %v", idx+1, *v, *v1)
			}

			if !reflect.DeepEqual(v.Props().Payload, v1.Props().Payload) {
				t.Errorf("test %d failed\nexpected payload type %T: %v\n received payload %T: %v", idx+1, v.Props().Payload, v.Props().Payload, v1.Props().Payload, v1.Props().Payload)
			}

			if !reflect.DeepEqual(v.Props().TxHash, v1.Props().TxHash) {
				t.Errorf("test %d failed\nexpected txHash: %v\n received txHash: %v", idx+1, v.Props().TxHash, v1.Props().TxHash)
			}
		}
	}
}
