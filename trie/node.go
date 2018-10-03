package trie

// Node ...
type Node struct {
	Key   []byte
	Value *Value
	Dirty bool
}

// NewNode ...
func NewNode(key []byte, val *Value, dirty bool) *Node {
	return &Node{
		Key:   key,
		Value: val,
		Dirty: dirty,
	}
}

// Copy ...
func (n *Node) Copy() *Node {
	cp := *n
	return &cp
}
