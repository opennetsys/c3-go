package p2p

import (
	"log"
	"reflect"
	"testing"

	"github.com/c3systems/c3/core/chain/mainchain"

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
		StateBlockHashes: []*string{
			&hash,
			&hash,
		},
		PrevBlockHash: "fakePrevHash",
		Nonce:         "fakeNonce",
		Difficulty:    "fakeDifficulty",
	})

	c, err := GetMainchainBlockCID(block)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("cid: %s", c.String())

	bytes, err := block.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	basicIPFSBlock, err := bfmt.NewBlockWithCid(bytes, c)
	if err != nil {
		t.Fatal(err)
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
		StateBlockHashes: []*string{
			&hash,
			&hash,
		},
		PrevBlockHash: "fakePrevHash",
		Nonce:         "fakeNonce",
		Difficulty:    "fakeDifficulty",
	})

	c, err := GetMainchainBlockCID(block)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("cid: %s", c.String())

	bytes, err := block.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	basicIPFSBlock, err := bfmt.NewBlockWithCid(bytes, c)
	if err != nil {
		t.Fatal(err)
	}

	if err := blocks.Put(basicIPFSBlock); err != nil {
		t.Fatal(err)
	}

	has, err := blocks.Has(c)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal("block store does not have our key!")
	}

	data, err := blocks.Get(c)
	if err != nil {
		t.Fatal(err)
	}

	received := new(mainchain.Block)
	if err := received.Deserialize(data.RawData()); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*block, received) {
		t.Errorf("expected: %v\nreceived: %v", *block, received)
	}
}
