// +build unit

package miner

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/common/fileutil"

	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/statechain"
	"github.com/c3systems/c3-go/core/p2p/mock"

	"github.com/golang/mock/gomock"
)

func TestCheckBlockHashAgainstDifficulty(t *testing.T) {
	t.Parallel()

	type input struct {
		block *mainchain.Block
	}
	type test struct {
		input    input
		expected bool
		err      error
	}

	var (
		hash1 string = hexutil.AddPrefix("01")
		hash2 string = hexutil.AddPrefix("1")
		hash3 string = "foo"
	)

	tests := []test{
		test{
			input: input{
				block: mainchain.New(&mainchain.Props{
					BlockHash:  &hash1,
					Difficulty: hexutil.EncodeUint64(1),
				}),
			},
			expected: true,
			err:      nil,
		},
		test{
			input: input{
				block: mainchain.New(&mainchain.Props{
					BlockHash:  &hash1,
					Difficulty: hexutil.EncodeUint64(2),
				}),
			},
			expected: false,
			err:      nil,
		},
		test{
			input: input{
				block: mainchain.New(&mainchain.Props{
					BlockHash:  &hash2,
					Difficulty: hexutil.EncodeUint64(0),
				}),
			},
			expected: true,
			err:      nil,
		},
		test{
			input: input{
				block: mainchain.New(&mainchain.Props{
					BlockHash:  &hash3,
					Difficulty: hexutil.EncodeUint64(2),
				}),
			},
			expected: false,
			err:      nil,
		},
	}

	for idx, tt := range tests {
		ok, err := CheckBlockHashAgainstDifficulty(tt.input.block)

		if tt.err != err {
			t.Errorf("test %d failed\nexpected err %v\nreceived err %v", idx+1, tt.err, err)
		}

		if tt.expected != ok {
			t.Errorf("test %d failed\nexpected %v\nreceived %v", idx+1, tt.expected, ok)
		}
	}
}

func TestCheckHashAgainstDifficulty(t *testing.T) {
	t.Parallel()

	type input struct {
		hashHex    string
		difficulty uint64
	}
	type test struct {
		input    input
		expected bool
		err      error
	}

	tests := []test{
		test{
			input: input{
				hashHex:    hexutil.AddPrefix("01"),
				difficulty: 1,
			},
			expected: true,
			err:      nil,
		},
		test{
			input: input{
				hashHex:    hexutil.AddPrefix("01"),
				difficulty: 2,
			},
			expected: false,
			err:      nil,
		},
		test{
			input: input{
				hashHex:    hexutil.AddPrefix("1"),
				difficulty: 0,
			},
			expected: true,
			err:      nil,
		},
		test{
			input: input{
				hashHex:    "foo",
				difficulty: 2,
			},
			expected: false,
			err:      nil,
		},
	}

	for idx, tt := range tests {
		ok, err := CheckHashAgainstDifficulty(tt.input.hashHex, tt.input.difficulty)

		if tt.err != err {
			t.Errorf("test %d failed\nexpected err %v\nreceived err %v", idx+1, tt.err, err)
		}

		if tt.expected != ok {
			t.Errorf("test %d failed\nexpected %v\nreceived %v", idx+1, tt.expected, ok)
		}
	}
}

func TestBuildTxsMap(t *testing.T) {
	t.Parallel()

	tx1 := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: "foo",
	})
	tx2 := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: "bar",
	})
	tx3 := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: "foobar",
	})
	transactions := []*statechain.Transaction{
		tx1,
		tx2,
		tx3,
		tx1,
		tx2,
		tx1,
	}

	m := BuildTxsMap(transactions)

	expected := make(statechain.TransactionsMap)
	expected[tx1.Props().ImageHash] = []*statechain.Transaction{
		tx1,
		tx1,
		tx1,
	}
	expected[tx2.Props().ImageHash] = []*statechain.Transaction{
		tx2,
		tx2,
	}
	expected[tx3.Props().ImageHash] = []*statechain.Transaction{
		tx3,
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected %v\nreceived %v", expected, m)
	}
}

func TestVerifyTransaction(t *testing.T) {
	t.Parallel()

	type test struct {
		input    *statechain.Transaction
		expected bool
		err      error
	}

	priv, pub, err := c3crypto.NewKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	addr, err := c3crypto.EncodeAddress(pub)
	if err != nil {
		t.Fatal(err)
	}

	fakeHash := "foo"
	nilSigTx := statechain.NewTransaction(&statechain.TransactionProps{
		TxHash: &fakeHash,
	})
	wrongHash := statechain.NewTransaction(&statechain.TransactionProps{
		TxHash: &fakeHash,
		From:   addr,
		Sig:    new(statechain.TxSig),
	})
	goodHash, err := wrongHash.CalculateHash()
	if err != nil {
		t.Fatal(err)
	}
	r, s, err := c3crypto.Sign(priv, []byte(goodHash))
	if err != nil {
		t.Fatal(err)
	}
	wrongSig := statechain.NewTransaction(&statechain.TransactionProps{
		TxHash: &goodHash,
		From:   addr,
		Sig: &statechain.TxSig{
			R: hexutil.EncodeBigInt(r),
			S: hexutil.EncodeBigInt(r),
		},
	})
	goodTx := statechain.NewTransaction(&statechain.TransactionProps{
		TxHash: &goodHash,
		From:   addr,
		Sig: &statechain.TxSig{
			R: hexutil.EncodeBigInt(r),
			S: hexutil.EncodeBigInt(s),
		},
	})

	tests := []test{
		test{
			input:    nil,
			expected: false,
			err:      ErrNilTx,
		},
		test{
			input:    new(statechain.Transaction),
			expected: false,
			err:      nil,
		},
		test{
			input:    nilSigTx,
			expected: false,
			err:      nil,
		},
		test{
			input:    wrongHash,
			expected: false,
			err:      nil,
		},
		test{
			input:    wrongSig,
			expected: false,
			err:      nil,
		},
		test{
			input:    goodTx,
			expected: true,
			err:      nil,
		},
	}

	for idx, tt := range tests {
		ok, err := VerifyTransaction(tt.input)

		if tt.err != err {
			t.Errorf("test %d failed\nexpected %v\nreceived %v", idx+1, tt.err, err)
		}

		if tt.expected != ok {
			t.Errorf("test %d failed\nexpected %v\nreceived %v", idx+1, tt.expected, ok)
		}
	}
}

func TestGatherDiffs(t *testing.T) {
	t.Parallel()

	// 1. mock the p2p service
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockP2P := mock_p2p.NewMockInterface(mockCtrl)

	// 2. build a fake statechain blocks and diff
	block := statechain.New(&statechain.BlockProps{
		BlockNumber:       "fakeNumber",
		StatePrevDiffHash: "foo",
		PrevBlockHash:     "bar",
	})
	genesisBlock := statechain.New(&statechain.BlockProps{
		BlockNumber:       mainchain.GenesisBlock.Props().BlockNumber,
		StatePrevDiffHash: "foo",
		PrevBlockHash:     "bar",
	})
	diff := statechain.NewDiff(&statechain.DiffProps{
		Data: "foobar",
	})

	// 3. add the expected mockp2p calls and returns
	mockP2P.
		EXPECT().
		GetStatechainDiff(gomock.Any()).
		Return(diff, nil)

	mockP2P.
		EXPECT().
		GetStatechainBlock(gomock.Any()).
		Return(block, nil)

	mockP2P.
		EXPECT().
		GetStatechainDiff(gomock.Any()).
		Return(diff, nil)

	mockP2P.
		EXPECT().
		GetStatechainBlock(gomock.Any()).
		Return(genesisBlock, nil)

	mockP2P.
		EXPECT().
		GetStatechainDiff(gomock.Any()).
		Return(diff, nil)

	// 4. run the function
	diffs, err := GatherDiffs(context.Background(), mockP2P, block)
	if err != nil {
		t.Fatal(err)
	}

	// 5. compare to the expected
	expected := []*statechain.Diff{
		diff,
		diff,
		diff,
	}

	if !reflect.DeepEqual(expected, diffs) {
		t.Errorf("expected %v\nreceived %v", expected, diffs)
	}

	// 6. do it again but this time with just a genesis block
	mockP2P.
		EXPECT().
		GetStatechainDiff(gomock.Any()).
		Return(diff, nil)

	// 7. run the function
	diffs, err = GatherDiffs(context.Background(), mockP2P, genesisBlock)
	if err != nil {
		t.Fatal(err)
	}

	// 8. compare to the expected
	expected = []*statechain.Diff{
		diff,
	}

	if !reflect.DeepEqual(expected, diffs) {
		t.Errorf("expected %v\nreceived %v", expected, diffs)
	}
}

func TestGenerateStateFromDiffs(t *testing.T) {
	t.Parallel()

	genesisState, err := ioutil.ReadFile("./test_data/tmp/state.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := ioutil.WriteFile("./test_data/tmp/state.txt", genesisState, os.ModePerm); err != nil {
			t.Logf("err returning genesis state file to it's original state\n%v", err)
		}
	}()
	expectedState, err := ioutil.ReadFile("./test_data/tmp/state3.txt")
	if err != nil {
		t.Fatal(err)
	}
	data1, err := ioutil.ReadFile("./test_data/1.patch")
	if err != nil {
		t.Fatal(err)
	}
	data2, err := ioutil.ReadFile("./test_data/2.patch")
	if err != nil {
		t.Fatal(err)
	}
	data3, err := ioutil.ReadFile("./test_data/3.patch")
	if err != nil {
		t.Fatal(err)
	}

	diff1 := statechain.NewDiff(&statechain.DiffProps{
		Data: string(data1),
	})
	diff2 := statechain.NewDiff(&statechain.DiffProps{
		Data: string(data2),
	})
	diff3 := statechain.NewDiff(&statechain.DiffProps{
		Data: string(data3),
	})

	diffs := []*statechain.Diff{
		diff1,
		diff2,
		diff3,
	}

	state, err := generateStateFromDiffs(context.TODO(), "fakeImage", genesisState, diffs)
	if err != nil {
		t.Fatal(err)
	}

	if string(expectedState) != string(state) {
		t.Errorf("expected %s\nreceived %s", string(expectedState), string(state))
	}
}

func TestGenerateCombinedDiffs(t *testing.T) {
	t.Parallel()

	expected, err := ioutil.ReadFile("./test_data/combined.patch")
	if err != nil {
		t.Fatal(err)
	}
	data1, err := ioutil.ReadFile("./test_data/1.patch")
	if err != nil {
		t.Fatal(err)
	}
	data2, err := ioutil.ReadFile("./test_data/2.patch")
	if err != nil {
		t.Fatal(err)
	}
	data3, err := ioutil.ReadFile("./test_data/3.patch")
	if err != nil {
		t.Fatal(err)
	}

	diff1 := statechain.NewDiff(&statechain.DiffProps{
		Data: string(data1),
	})
	diff2 := statechain.NewDiff(&statechain.DiffProps{
		Data: string(data2),
	})
	diff3 := statechain.NewDiff(&statechain.DiffProps{
		Data: string(data3),
	})

	diffs := []*statechain.Diff{
		diff1,
		diff2,
		diff3,
	}

	received, err := generateCombinedDiffs(context.TODO(), "fakeImage", diffs)
	if err != nil {
		t.Fatal(err)
	}

	if string(expected) != string(received) {
		t.Errorf("expected %s\n\n\nreceived %s", string(expected), string(received))
	}
}

func TestIsGenesisTransaction(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockP2P := mock_p2p.NewMockInterface(mockCtrl)
	txDeploy := statechain.NewTransaction(&statechain.TransactionProps{
		Method: "c3_deploy",
	})
	txAction := statechain.NewTransaction(&statechain.TransactionProps{
		Method: "c3_invokeMethod",
	})
	stateBlock := statechain.New(&statechain.BlockProps{})

	p2pErr := errors.New("p2p err")

	imageHash := "imageHash"
	mainBlock := mainchain.New(&mainchain.Props{})

	// 1. An err with p2pSvc should throw
	txs := []*statechain.Transaction{
		txAction,
		txAction,
		txDeploy,
		txAction,
	}

	mockP2P.
		EXPECT().
		FetchMostRecentStateBlock(imageHash, mainBlock).
		Return(nil, p2pErr)

	if _, _, _, err := isGenesisTransaction(mockP2P, mainBlock, imageHash, txs); err != p2pErr {
		t.Errorf("expected %v\nreceived %v", p2pErr, err)
	}

	// 2. A deploy tx with a previous state block should throw
	mockP2P.
		EXPECT().
		FetchMostRecentStateBlock(imageHash, mainBlock).
		Return(stateBlock, nil)

	if _, _, _, err := isGenesisTransaction(mockP2P, mainBlock, imageHash, txs); err == nil {
		t.Error("expected err but received nil")
	}

	// 3. A deploy tx without a previous state block should return the expected data
	mockP2P.
		EXPECT().
		FetchMostRecentStateBlock(imageHash, mainBlock).
		Return(nil, nil)

	expectedTxs := []*statechain.Transaction{
		txAction,
		txAction,
		txAction,
	}

	isGenesis, tx, remainingTxs, err := isGenesisTransaction(mockP2P, mainBlock, imageHash, txs)
	if err != nil || isGenesis != true || !reflect.DeepEqual(tx, txDeploy) || !reflect.DeepEqual(remainingTxs, expectedTxs) {
		t.Errorf("expected nil, true, %v, %v\nreceived %v, %v, %v, %v", txDeploy, expectedTxs, err, isGenesis, tx, remainingTxs)
	}

	// 4. A non deploy txs should just return
	txs = expectedTxs

	isGenesis, tx, remainingTxs, err = isGenesisTransaction(mockP2P, mainBlock, imageHash, txs)
	if err != nil || isGenesis != false || tx != nil || !reflect.DeepEqual(remainingTxs, expectedTxs) {
		t.Errorf("expected nil, false, nil, %v\nreceived %v, %v, %v, %v", expectedTxs, err, isGenesis, tx, remainingTxs)
	}
}

func TestOrderStatechainBlocks(t *testing.T) {
	t.Parallel()

	// 1. nil blocks returns an err
	if _, err := orderStatechainBlocks(nil); err == nil {
		t.Error("expected err but received nil")
	}

	// 2. a nil statechainblock returns an err
	if _, err := orderStatechainBlocks([]*statechain.Block{nil}); err == nil {
		t.Error("expected err but received nil")
	}

	// 3. an invalid blockNumber returns an err
	invalid := statechain.New(&statechain.BlockProps{
		BlockNumber: "foo",
	})
	if _, err := orderStatechainBlocks([]*statechain.Block{invalid}); err == nil {
		t.Error("expected err but received nil")
	}

	// 4. valid blocks should return, sorted
	b1 := statechain.New(&statechain.BlockProps{
		BlockNumber: "0x1",
	})
	b2 := statechain.New(&statechain.BlockProps{
		BlockNumber: "0x2",
	})
	b3 := statechain.New(&statechain.BlockProps{
		BlockNumber: "0x3",
	})
	b4 := statechain.New(&statechain.BlockProps{
		BlockNumber: "0x4",
	})
	b5 := statechain.New(&statechain.BlockProps{
		BlockNumber: "0x5",
	})

	inputs := []*statechain.Block{
		b2,
		b4,
		b3,
		b5,
		b1,
	}
	expected := []*statechain.Block{
		b1,
		b2,
		b3,
		b4,
		b5,
	}

	blocks, err := orderStatechainBlocks(inputs)
	if err != nil || !reflect.DeepEqual(expected, blocks) {
		t.Errorf("expected nil err and %v\nreceived %v and %v", expected, err, blocks)
	}
}

func TestGroupStateBlocksByImageHash(t *testing.T) {
	t.Parallel()

	// 1. a nil stateblocks map should return an empty map
	ret, err := groupStateBlocksByImageHash(nil)
	if err != nil || ret == nil || len(ret) != 0 {
		t.Errorf("expected nil err, non nil ret with len == 0\nreceived %v, %v, %v", err, ret, len(ret))
	}

	// 2. an initialized, empty map should return an empty map
	tmpIn := make(map[string]*statechain.Block)
	ret, err = groupStateBlocksByImageHash(tmpIn)
	if err != nil || ret == nil || len(ret) != 0 {
		t.Errorf("expected nil err, non nil ret with len == 0\nreceived %v, %v, %v", err, ret, len(ret))
	}

	// 3. an input with a nil block should return an err
	tmpIn = make(map[string]*statechain.Block)
	tmpIn["foo"] = nil
	if _, err = groupStateBlocksByImageHash(tmpIn); err == nil {
		t.Error("expected non nil err but received nil")
	}

	// 4. a proper input should return transactions grouped by image hash
	foo1Hash := "f1"
	foo1 := statechain.New(&statechain.BlockProps{
		BlockHash: &foo1Hash,
		ImageHash: "foo",
	})
	foo2Hash := "f2"
	foo2 := statechain.New(&statechain.BlockProps{
		BlockHash: &foo2Hash,
		ImageHash: "foo",
	})
	bar1Hash := "b1"
	bar1 := statechain.New(&statechain.BlockProps{
		BlockHash: &bar1Hash,
		ImageHash: "bar",
	})

	tmpIn = make(map[string]*statechain.Block)
	tmpIn[foo1Hash] = foo1
	tmpIn[foo2Hash] = foo2
	tmpIn[bar1Hash] = bar1

	expected := make(map[string][]*statechain.Block)
	expected["foo"] = []*statechain.Block{foo1, foo2}
	expected["bar"] = []*statechain.Block{bar1}

	ret, err = groupStateBlocksByImageHash(tmpIn)
	if err != nil || len(ret) != 2 || !reflect.DeepEqual(ret["bar"], expected["bar"]) {
		t.Errorf("expected nil err and %v\nreceived %v, %v", err, expected, ret)
	}

	if len(ret["foo"]) != 2 || !((ret["foo"][0] == foo1 && ret["foo"][1] == foo2) || (ret["foo"][0] == foo2 && ret["foo"][1] == foo1)) {
		t.Errorf("expected foo's %v\nreceeived %v", expected["foo"], ret["foo"])
	}
}

func TestCleanupFiles(t *testing.T) {
	t.Parallel()

	// 1. nil fileNames should return
	cleanupFiles(nil)

	// 2. an array of names should be cleaned
	f1, err := os.Create("./test_data/foo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f1.Name())
	if err = f1.Close(); err != nil {
		t.Fatal(err)
	}

	f2, err := os.Create("./test_data/bar")
	if err != nil {
		t.Fatal(err)
	}
	if err = f2.Close(); err != nil {
		t.Fatal(err)
	}

	fileNames := &[]string{
		f1.Name(),
		f2.Name(),
	}

	cleanupFiles(fileNames)

	if _, err = os.Stat(f1.Name()); !os.IsNotExist(err) {
		t.Errorf("expected file %s to not exist", f1.Name())
	}
	if _, err = os.Stat(f2.Name()); !os.IsNotExist(err) {
		t.Errorf("expected file %s to not exist", f2.Name())
	}

	// 3. TODO: test that an error is printed to stdout
}

func TestMakeTempFile(t *testing.T) {
	t.Parallel()

	file := "baz"
	fileName := fmt.Sprintf("%s/%s/%s", "foo", "bar", file)
	f, err := fileutil.CreateTempFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(t.Name())
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err = os.Stat(f.Name()); err != nil {
		t.Fatalf("expected nil err but received %v", err)
	}

	prefix := "/tmp/foo_bar"
	name := f.Name()

	if len(name) <= len(prefix) {
		t.Fatalf("expected filename of at least %d\nreceived %d", len(prefix), len(name))
	}
	if name[:len(prefix)] != prefix {
		t.Errorf("expected %s prefix\nreceived %s", prefix, name[:len(prefix)])
	}
	if name[len(name)-len(file):] != file {
		t.Errorf("expected %s ending\nreceived %s", file, name[len(name)-len(file):])
	}
}
