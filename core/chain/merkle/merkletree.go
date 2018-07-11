package merkle

import (
	"bytes"
	"encoding/gob"

	"github.com/c3systems/c3/common/hexutil"
	chaintypes "github.com/c3systems/c3/core/chain/types"

	"github.com/c3systems/merkletree"
)

// New ...
func New(props *TreeProps) (*Tree, error) {
	if props == nil {
		return nil, ErrNilProps
	}

	return &Tree{
		props: *props,
	}, nil
}

// BuildFromObjects ...
// note: using this method to keep with our New function accepting props
// TODO: ensure that all of the kinds are the same?
// TODO ensure that the kind is the true kind?
func BuildFromObjects(chainObjects []chaintypes.ChainObject, kind string) (*Tree, error) {
	if chainObjects == nil || len(chainObjects) == 0 {
		return nil, ErrNilChainObjects
	}

	var (
		list   []merkletree.Content
		hashes []string
		err    error
	)

	if ok := checkKind(kind); !ok {
		return nil, ErrUnknownKind
	}

	for _, chainObject := range chainObjects {
		if chainObject == nil {
			return nil, ErrNilChainObject
		}

		list = append(list, chainObject)

		hash, err := chainObject.CalculateHashBytes()
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, string(hash))
	}

	t, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}

	mr := t.MerkleRoot()
	mrRootHash := hexutil.EncodeString(string(mr))
	return &Tree{
		props: TreeProps{
			MerkleTreeRootHash: &mrRootHash,
			Kind:               kind,
			Hashes:             hashes,
		},
	}, nil
}

// Props ...
func (t *Tree) Props() TreeProps {
	return t.props
}

// Serialize ...
func (t *Tree) Serialize() ([]byte, error) {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(t.props)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Deserialize ...
func (t *Tree) Deserialize(data []byte) error {
	if t == nil {
		return ErrNilMerkleTree
	}

	var tmpProps TreeProps
	b := bytes.NewBuffer(data)
	if err := gob.NewDecoder(b).Decode(&tmpProps); err != nil {
		return err
	}

	t.props = tmpProps
	return nil
}

// SerializeString ...
func (t *Tree) SerializeString() (string, error) {
	data, err := t.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(data)), nil
}

// DeserializeString ...
func (t *Tree) DeserializeString(hexStr string) error {
	if t == nil {
		return ErrNilMerkleTree
	}

	str, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return t.Deserialize([]byte(str))
}

// CalculateHash ...
// note: this function doesn't actually calculate the hash but returns the root hash
func (t *Tree) CalculateHash() (string, error) {
	if t.props.MerkleTreeRootHash == nil {
		return "", ErrNilMerkleTreeRootHash
	}

	return *t.props.MerkleTreeRootHash, nil
}

// CalculateHashBytes ...
func (t *Tree) CalculateHashBytes() ([]byte, error) {
	hash, err := t.CalculateHash()
	if err != nil {
		return nil, err
	}

	return []byte(hash), nil
}

// Equals ...
func (t *Tree) Equals(other merkletree.Content) (bool, error) {
	tHash, err := t.CalculateHashBytes()
	if err != nil {
		return false, err
	}

	oHash, err := other.CalculateHashBytes()
	if err != nil {
		return false, err
	}

	return string(tHash) == string(oHash), nil
}

// SetHash ...
func (t *Tree) SetHash() error {
	if t == nil {
		return ErrNilMerkleTree
	}

	hash, err := t.CalculateHash()
	if err != nil {
		return err
	}

	t.props.MerkleTreeRootHash = &hash

	return nil
}

func checkKind(kind string) bool {
	for _, allowedKind := range allowedKinds {
		if allowedKind != kind {
			return false
		}
	}

	return true
}
