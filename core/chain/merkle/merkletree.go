package merkle

import (
	"encoding/json"
	"errors"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3-go/common/coder"
	"github.com/c3systems/c3-go/common/hexutil"
	loghooks "github.com/c3systems/c3-go/log/hooks"

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
		mrRootHash := "0x0"
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
		log.Println("[merkle] unknown kind")
		return nil, ErrUnknownKind
	}

	log.Println("[merkle] building tree")

	for _, chainObject := range chainObjects {
		log.Printf("[merkle] chain object %v", chainObject)

		if chainObject == nil {
			log.Println("[merkle] nil chain object")
			return nil, ErrNilChainObject
		}

		list = append(list, chainObject)

		hash, err := chainObject.CalculateHashBytes()
		if err != nil {
			log.Printf("[merkle] error calculating hash bytes; %s", err)
			return nil, err
		}

		log.Printf("[merkle] appending hash %s", hexutil.EncodeToString(hash))
		hashes = append(hashes, hexutil.EncodeToString(hash))
	}

	t, err := merkletree.NewTree(list)
	if err != nil {
		log.Printf("[merkle] error creating new merkle tree; %s", err)
		return nil, err
	}

	mr := t.MerkleRoot()
	mrRootHash := hexutil.EncodeToString(mr)
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
	tmp := BuildCoderFromTree(t)
	bytes, err := tmp.Marshal()
	if err != nil {
		return nil, err
	}

	return coder.AppendCode(bytes), nil
}

// Deserialize ...
func (t *Tree) Deserialize(data []byte) error {
	if data == nil {
		return errors.New("nil bytes")
	}
	if t == nil {
		return errors.New("nil tree")
	}

	_, bytes, err := coder.StripCode(data)
	if err != nil {
		return err
	}

	props, err := BuildTreePropsFromBytes(bytes)
	if err != nil {
		return err
	}

	t.props = *props

	return nil
}

// SerializeString ...
func (t *Tree) SerializeString() (string, error) {
	data, err := t.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeToString(data), nil
}

// DeserializeString ...
func (t *Tree) DeserializeString(hexStr string) error {
	if t == nil {
		return ErrNilMerkleTree
	}

	b, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return t.Deserialize(b)
}

// CalculateHash ...
// note: this function doesn't actually calculate the hash but returns the root hash
func (t *Tree) CalculateHash() (string, error) {
	if t.props.MerkleTreeRootHash != nil {
		return *t.props.MerkleTreeRootHash, nil
	}

	var tmpContent []merkletree.Content
	log.Printf("[merkle] calculate hash - hashes length: %v", len(t.props.Hashes))
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

	return reflect.DeepEqual(tHash, oHash), nil
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
		if allowedKind == kind {
			return true
		}
	}

	return false
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

// BuildCoderFromTree ...
func BuildCoderFromTree(t *Tree) *coder.MerkleTree {
	tmp := &coder.MerkleTree{
		Kind:   t.props.Kind,
		Hashes: t.props.Hashes,
	}

	// note: is there a better way to handle nil with protobuff?
	if t.props.MerkleTreeRootHash != nil {
		tmp.MerkleTreeRootHash = *t.props.MerkleTreeRootHash
	}

	return tmp
}

// BuildTreePropsFromBytes ...
func BuildTreePropsFromBytes(data []byte) (*TreeProps, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	c, err := BuildCoderFromBytes(data)
	if err != nil {
		return nil, err
	}

	return BuildTreePropsFromCoder(c)
}

// BuildCoderFromBytes ...
func BuildCoderFromBytes(data []byte) (*coder.MerkleTree, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	tmp := new(coder.MerkleTree)
	if err := tmp.Unmarshal(data); err != nil {
		return nil, err
	}
	if tmp == nil {
		return nil, errors.New("nil output")
	}

	return tmp, nil
}

// BuildTreePropsFromCoder ...
func BuildTreePropsFromCoder(tmp *coder.MerkleTree) (*TreeProps, error) {
	if tmp == nil {
		return nil, errors.New("nil coder")
	}

	props := &TreeProps{
		Kind: tmp.Kind,
	}

	// note: is there a better way to handle nil with protobuff?
	if tmp.MerkleTreeRootHash != "" {
		s := tmp.MerkleTreeRootHash
		props.MerkleTreeRootHash = &s
	}
	if tmp.Hashes != nil && len(tmp.Hashes) > 0 {
		props.Hashes = tmp.Hashes
	}

	return props, nil
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
