# Mock
## Install
```bash
$ go get github.com/golang/mock/gomock
$ go install github.com/golang/mock/mockgen
```

## Usage
In the main package directory: `$ mockgen -source=interface.go -destination=mock/mock.go`

## Testing
```go
func foo(p2pSvc ptp.Interface) (*statechain.Diff, error) {
  // foo does stuff, here...
}

func TestFoo(t *testing.T) {
	// 1. mock the p2p service
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

  mockP2P := mock_p2p.NewMockInterface(mockCtrl)
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

	// 2. run the function
	diffs, err := foo(mockP2P)
	if err != nil {
		t.Fatal(err)
  }
    
  // etc...
} 
```
