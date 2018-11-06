package eosclient

import (
	eos "github.com/eoscanada/eos-go"
)

// Config ...
type Config struct {
	URL string
}

// Client ...
type Client struct {
	client *eos.API
}

// NewClient ...
func NewClient(config *Config) *Client {
	client := eos.New(config.URL)

	return &Client{
		client: client,
	}
}

// Info ...
func (s *Client) Info() (*eos.InfoResp, error) {
	return s.client.GetInfo()
}

// AccountInfo ...
func (s *Client) AccountInfo(account string) (*eos.AccountResp, error) {
	acct := eos.AccountName(account)
	return s.client.GetAccount(acct)
}
