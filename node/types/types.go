package types

// NewAddressResponse ...
type NewAddressResponse struct {
	Address string
}

// SendTxResponse ...
type SendTxResponse struct {
	TxHash *string
}

// GetInfoResponse ...
type GetInfoResponse struct {
	BlockHeight string
}

// CFG ...
type CFG struct {
	URI     string
	Peer    string
	DataDir string
}
