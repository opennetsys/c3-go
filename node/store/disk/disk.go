package disk

// TODO: everything...

// Props ...
type Props struct {
}

// Service ...
type Service struct {
	props Props
}

// New ...
func New(props *Props) (*Service, error) {
	return nil, nil
}

// Props ...
func (s Service) Props() Props {
	return s.props
}

func (s Service) HasTx(hash string) (bool, error) {
	return false, nil
}

func (s Service) GetTx(hash string) (*statechain.Transaction, error) {
	return nil, nil
}

func (s Service) GetTxs(hashes []string) ([]*statechain.Transaction, error) {
	return nil, nil
}

func (s Service) RemoveTx(hash string) error {
	return nil
}

// RemoveTxs ...
func (s Service) RemoveTxs(hashes []string) error {
	return nil
}

// AddTx ...
func (s Service) AddTx(tx *statechain.Transaction) error {
	return nil
}

// GatherTransactions ...
func (s Service) GatherTransactions() (*[]statechain.Transaction, error) {
	return nil
}
