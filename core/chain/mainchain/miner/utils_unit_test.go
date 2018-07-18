// +build unit

package miner

import (
	"context"
	"reflect"
	"testing"

	"github.com/c3systems/c3/common/c3crypto"

	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/p2p/mock"

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
		hash1 string = hexutil.AddLeader("01")
		hash2 string = hexutil.AddLeader("1")
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
			err:      hexutil.ErrNotHexString,
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
				hashHex:    hexutil.AddLeader("01"),
				difficulty: 1,
			},
			expected: true,
			err:      nil,
		},
		test{
			input: input{
				hashHex:    hexutil.AddLeader("01"),
				difficulty: 2,
			},
			expected: false,
			err:      nil,
		},
		test{
			input: input{
				hashHex:    hexutil.AddLeader("1"),
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
			err:      hexutil.ErrNotHexString,
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
	diffs, err := gatherDiffs(context.Background(), mockP2P, block)
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
	diffs, err = gatherDiffs(context.Background(), mockP2P, genesisBlock)
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
	t.Skip()

	// TODO
}

func TestGenerateCombinedDiffs(t *testing.T) {
	t.Skip()

	// TODO
}

func TestIsGenesisTransaction(t *testing.T) {
	t.Skip()

	// TODO
}

func TestOrderStatechainBlocks(t *testing.T) {
	t.Skip()

	// TODO
}

func TestGroupStateBlocksByImageHash(t *testing.T) {
	t.Skip()

	// TODO
}

func TestCleanupFiles(t *testing.T) {
	t.Skip()

	// TODO
}

func TestMakeTempFile(t *testing.T) {
	t.Skip()

	// TODO
}
