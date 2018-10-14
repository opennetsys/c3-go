package trie

import "sync"

// Trie structure. Trie will cache values first and data is only persisted when sync is invoked.
type Trie struct {
	mut      sync.RWMutex
	prevRoot interface{}
	Root     interface{}
	cache    *Cache
}

// 17 is chosen because there's 16 characters in hex, plus 1 delimiter
var maxSize = 17

func copyRoot(root interface{}) interface{} {
	var prevRootCopy interface{}
	if b, ok := root.([]byte); ok {
		prevRootCopy = CopyBytes(b)
	} else {
		prevRootCopy = root
	}

	return prevRootCopy
}

// NewTrie ...
func NewTrie(db Database, Root interface{}) *Trie {
	r := copyRoot(Root)
	p := copyRoot(Root)

	return &Trie{
		cache:    NewCache(db),
		Root:     r,
		prevRoot: p,
	}
}

// Update ...
func (t *Trie) Update(key, value string) {
	t.mut.Lock()
	defer t.mut.Unlock()

	k := CompactHexDecode(key)

	root := t.updateState(t.Root, k, value)
	switch root.(type) {
	case string:
		t.Root = root
	case []byte:
		t.Root = root
	default:
		t.Root = t.cache.putValue(root, true)
	}
}

// Get ...
func (t *Trie) Get(key string) string {
	t.mut.RLock()
	defer t.mut.RUnlock()

	k := CompactHexDecode(key)
	c := NewValue(t.getState(t.Root, k))

	return c.Str()
}

// Delete ...
func (t *Trie) Delete(key string) {
	t.mut.Lock()
	defer t.mut.Unlock()

	k := CompactHexDecode(key)

	root := t.deleteState(t.Root, k)
	switch root.(type) {
	case string:
		t.Root = root
	case []byte:
		t.Root = root
	default:
		t.Root = t.cache.putValue(root, true)
	}
}

// Cmp ...
func (t *Trie) Cmp(tr *Trie) bool {
	return NewValue(t.Root).Cmp(NewValue(tr.Root))
}

// Copy ...
func (t *Trie) Copy() *Trie {
	tr := NewTrie(t.cache.db, t.Root)
	for key, node := range t.cache.nodes {
		tr.cache.nodes[key] = node.Copy()
	}

	return tr
}

// Sync ...
func (t *Trie) Sync() {
	t.cache.Commit()
	t.prevRoot = copyRoot(t.Root)
}

// Undo ...
func (t *Trie) Undo() {
	t.cache.Undo()
	t.Root = t.prevRoot
}

// Cache ...
func (t *Trie) Cache() *Cache {
	return t.cache
}

func (t *Trie) getState(node interface{}, key []int) interface{} {
	n := NewValue(node)
	if len(key) == 0 || n.IsNil() || n.Size() == 0 {
		return node
	}

	currentNode := t.getNode(node)
	length := currentNode.Size()

	if length == 0 {
		return ""
	} else if length == 2 {
		k := CompactDecode(currentNode.Get(0).Str())
		v := currentNode.Get(1).Raw()

		if len(key) >= len(k) && CompareIntSlice(k, key[:len(k)]) {
			return t.getState(v, key[len(k):])
		}

		return ""
	} else if length == maxSize {
		return t.getState(currentNode.Get(key[0]).Raw(), key[1:])
	}

	panic("unexpecd return")
}

func (t *Trie) getNode(node interface{}) *Value {
	n := NewValue(node)

	if !n.Get(0).IsNil() {
		return n
	}

	str := n.Str()
	if len(str) == 0 {
		return n
	} else if len(str) < 32 {
		return NewValueFromBytes([]byte(str))
	}

	data := t.cache.Get(n.Bytes())
	return data
}

// updateState ...
func (t *Trie) updateState(node interface{}, key []int, value string) interface{} {
	return t.insertState(node, key, value)
}

// Put ...
func (t *Trie) Put(node interface{}) interface{} {
	return t.cache.Put(node)
}

// EmptyStringSlice ...
func EmptyStringSlice(l int) []interface{} {
	slice := make([]interface{}, l)
	for i := 0; i < l; i++ {
		slice[i] = ""
	}

	return slice
}

// insertState ...
func (t *Trie) insertState(node interface{}, key []int, value interface{}) interface{} {
	if len(key) == 0 {
		return value
	}

	n := NewValue(node)
	if node == nil || n.Size() == 0 {
		newNode := []interface{}{CompactEncode(key), value}

		return t.Put(newNode)
	}

	currentNode := t.getNode(node)
	// check for special 2 slice type node
	if currentNode.Size() == 2 {
		k := CompactDecode(currentNode.Get(0).Str())
		v := currentNode.Get(1).Raw()

		if CompareIntSlice(k, key) {
			newNode := []interface{}{CompactEncode(key), value}
			return t.Put(newNode)
		}

		var newHash interface{}
		matchingLength := MatchingNibbleLength(key, k)
		if matchingLength == len(k) {
			newHash = t.insertState(v, key[matchingLength:], value)
		} else {
			// expand 2 length slice to a maxSize length slice
			oldNode := t.insertState("", k[matchingLength+1:], v)
			newNode := t.insertState("", key[matchingLength+1:], value)
			scaledSlice := EmptyStringSlice(maxSize)
			scaledSlice[k[matchingLength]] = oldNode
			scaledSlice[key[matchingLength]] = newNode

			newHash = t.Put(scaledSlice)
		}

		if matchingLength == 0 {
			// end of chain, return
			return newHash
		}

		newNode := []interface{}{CompactEncode(key[:matchingLength]), newHash}
		return t.Put(newNode)
	}
	// copy
	newNode := EmptyStringSlice(maxSize)
	for i := 0; i < maxSize; i++ {
		cpy := currentNode.Get(i).Raw()
		if cpy != nil {
			newNode[i] = cpy
		}
	}

	newNode[key[0]] = t.insertState(currentNode.Get(key[0]).Raw(), key[1:], value)

	return t.Put(newNode)
}

func (t *Trie) deleteState(node interface{}, key []int) interface{} {
	if len(key) == 0 {
		return ""
	}

	if node == nil || NewValue(node).Size() == 0 {
		return ""
	}

	currentNode := t.getNode(node)
	if currentNode.Size() == 2 {
		k := CompactDecode(currentNode.Get(0).Str())
		v := currentNode.Get(1).Raw()

		if CompareIntSlice(k, key) {
			return ""
		} else if CompareIntSlice(key[:len(k)], k) {
			hash := t.deleteState(v, key[len(k):])
			child := t.getNode(hash)

			var newNode []interface{}
			if child.Size() == 2 {
				newKey := append(k, CompactDecode(child.Get(0).Str())...)
				newNode = []interface{}{CompactEncode(newKey), child.Get(1).Raw()}
			} else {
				newNode = []interface{}{currentNode.Get(0).Str(), hash}
			}

			return t.Put(newNode)
		}
		return node
	}

	// copy the current node over to the new node and replace the first nibble in the key
	n := EmptyStringSlice(maxSize)
	var newNode []interface{}

	for i := 0; i < maxSize; i++ {
		cpy := currentNode.Get(i).Raw()
		if cpy != nil {
			n[i] = cpy
		}
	}

	n[key[0]] = t.deleteState(n[key[0]], key[1:])
	amount := -1
	for i := 0; i < maxSize; i++ {
		if n[i] != "" {
			if amount == -1 {
				amount = i
			} else {
				amount = -2
			}
		}
	}
	if amount == 16 {
		newNode = []interface{}{CompactEncode([]int{16}), n[amount]}
	} else if amount >= 0 {
		child := t.getNode(n[amount])
		if child.Size() == maxSize {
			newNode = []interface{}{CompactEncode([]int{amount}), n[amount]}
		} else if child.Size() == 2 {
			key := append([]int{amount}, CompactDecode(child.Get(0).Str())...)
			newNode = []interface{}{CompactEncode(key), child.Get(1).Str()}
		}
	} else {
		newNode = n
	}

	return t.Put(newNode)
}
