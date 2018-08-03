package p2p

import (
	"reflect"
	"testing"

	"github.com/c3systems/c3-go/core/chain/mainchain"

	bfmt "github.com/ipfs/go-block-format"
	datastore "github.com/ipfs/go-datastore"
	bstore "github.com/ipfs/go-ipfs-blockstore"
)

func TestBasicBlock(t *testing.T) {
	hash := "fakeHash"
	block := mainchain.New(&mainchain.Props{
		BlockHash:             &hash,
		BlockNumber:           "1",
		BlockTime:             "0",
		ImageHash:             "fakeImageHash",
		StateBlocksMerkleHash: "fakeMerkle",
		PrevBlockHash:         "fakePrevHash",
		Nonce:                 "fakeNonce",
		Difficulty:            "fakeDifficulty",
	})

	c, err := GetMainchainBlockCID(block)
	if err != nil {
		t.Error(err)
	}
	t.Logf("cid: %s", c.String())

	bytes, err := block.Serialize()
	if err != nil {
		t.Error(err)
	}

	basicIPFSBlock, err := bfmt.NewBlockWithCid(bytes, c)
	if err != nil {
		t.Error(err)
	}
	t.Logf("basic cid: %s", basicIPFSBlock.Cid().String())
}

func TestLocalPutAndFetchMainchainBlock(t *testing.T) {
	memStore := datastore.NewMapDatastore()
	blocks := bstore.NewBlockstore(memStore)

	hash := "fakeHash"
	block := mainchain.New(&mainchain.Props{
		BlockHash:             &hash,
		BlockNumber:           "1",
		BlockTime:             "0",
		ImageHash:             "fakeImageHash",
		StateBlocksMerkleHash: "fakeMerkle",
		PrevBlockHash:         "fakePrevHash",
		Nonce:                 "fakeNonce",
		Difficulty:            "fakeDifficulty",
	})

	c, err := GetMainchainBlockCID(block)
	if err != nil {
		t.Error(err)
	}
	t.Logf("cid: %s", c.String())

	bytes, err := block.Serialize()
	if err != nil {
		t.Error(err)
	}

	basicIPFSBlock, err := bfmt.NewBlockWithCid(bytes, c)
	if err != nil {
		t.Error(err)
	}

	if err := blocks.Put(basicIPFSBlock); err != nil {
		t.Error(err)
	}

	has, err := blocks.Has(c)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("block store does not have our key!")
	}

	data, err := blocks.Get(c)
	if err != nil {
		t.Error(err)
	}

	received := new(mainchain.Block)
	if err := received.Deserialize(data.RawData()); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(*block, received) {
		t.Errorf("expected: %v\nreceived: %v", *block, received)
	}
}
