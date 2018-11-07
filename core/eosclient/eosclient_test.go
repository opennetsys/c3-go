package eosclient

import (
	"fmt"
	"testing"

	hexutil "github.com/c3systems/c3-go/common/hexutil"
)

func TestNewClient(t *testing.T) {
	_ = NewClient(&Config{
		URL: "http://api.kylin.alohaeos.com",
	})
}

func TestInfo(t *testing.T) {
	client := NewClient(&Config{
		URL: "http://api.kylin.alohaeos.com",
	})

	info, err := client.Info()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(info)
}

func TestAccountInfo(t *testing.T) {
	client := NewClient(&Config{
		URL: "http://api.kylin.alohaeos.com",
	})

	info, err := client.AccountInfo("helloworld54")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(info)
}

func TestGetTransaction(t *testing.T) {
	client := NewClient(&Config{
		URL:   "https://api-kylin.eosasia.one",
		Debug: true,
	})

	tx, err := client.GetTransaction("3d43785ceca9a919e73b547487d9da6dad246f05425e513035e373c67310bc47")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(tx)
}

func TestPushAction(t *testing.T) {
	client := NewClient(&Config{
		URL: "http://api.kylin.alohaeos.com",
	})

	randhex := hexutil.RandomHex(64)
	root := "0x" + randhex

	action := &Action{
		ActionName:  "chkpointroot",
		AccountName: "helloworld54",
		Permissions: "helloworld54@active",
		ActionData: &CheckpointData{
			Root: root,
		},
	}

	wifPrivateKey := "5Jh9tD4Fp1EpVn3EzEW6ura5NV3NddY8NNBcfpCZTvPDsKd9i5c"
	client.SetSigner(wifPrivateKey)

	resp, err := client.PushAction(action)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(resp)
}

func TestCheckpointRoot(t *testing.T) {
	client := NewCheckpointClient(&CheckpointConfig{
		URL:           "https://api-kylin.eosasia.one",
		ActionName:    "chkpointroot",
		AccountName:   "helloworld54",
		Permissions:   "helloworld54@active",
		WifPrivateKey: "5Jh9tD4Fp1EpVn3EzEW6ura5NV3NddY8NNBcfpCZTvPDsKd9i5c",
	})

	randhex := hexutil.RandomHex(64)
	root := "0x" + randhex

	fmt.Println(root)
	resp, err := client.CheckpointRoot(root)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.TransactionID)
	fmt.Println(resp)
}
