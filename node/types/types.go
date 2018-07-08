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

// Config ...
type Config struct {
	URI     string
	Peer    string
	DataDir string
}
