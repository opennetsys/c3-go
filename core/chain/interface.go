package chain

type Interface interface {
	AddMainBlock(block *mainchain.Block) *cid.CID
	Transactions() []*statechain.Transaction
	MainHead() (*mainchain.Block, error)
	StateHead(hash string) (*statechain.Block, error)
}
