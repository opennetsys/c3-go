package eosclient

import (
	"fmt"
	"testing"
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
