package merkle

import (
	"encoding/json"

	"github.com/c3systems/c3/common/hexutil"

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
func BuildFromObjects(chainObjects []merkletree.Content, kind string) (*Tree, error) {
	if chainObjects == nil || len(chainObjects) == 0 {
		hashes := []string{}
		mrRootHash := hexutil.EncodeString("")
		return &Tree{
			props: TreeProps{
				MerkleTreeRootHash: &mrRootHash,
				Kind:               kind,
				Hashes:             hashes,
			},
		}, nil
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
	return t.MarshalJSON()
}

// Deserialize ...
func (t *Tree) Deserialize(data []byte) error {
	return t.UnmarshalJSON(data)
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
	if t.props.MerkleTreeRootHash != nil {
		return *t.props.MerkleTreeRootHash, nil
	}

	var tmpContent []merkletree.Content
	for _, str := range t.props.Hashes {
		tmpContent = append(tmpContent, testContent{
			x: str,
		})
	}

	tmpTree, err := BuildFromObjects(tmpContent, t.props.Kind)
	if err != nil {
		return "", err
	}
	if tmpTree == nil {
		return "", ErrNilMerkleTree
	}
	if tmpTree.Props().MerkleTreeRootHash == nil {
		return "", ErrNilMerkleTreeRootHash
	}

	return *tmpTree.props.MerkleTreeRootHash, nil
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

// MarshalJSON ...
func (t *Tree) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.props)
}

// UnmarshalJSON ...
func (t *Tree) UnmarshalJSON(data []byte) error {
	var props TreeProps
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	t.props = props

	return nil
}
