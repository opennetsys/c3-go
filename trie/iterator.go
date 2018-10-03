package trie

// TrieIterator ...
type TrieIterator struct {
	trie  *Trie
	key   string
	value string

	shas   [][]byte
	values []string

	lastNode []byte
}

// NewIterator ...
func (t *Trie) NewIterator() *TrieIterator {
	return &TrieIterator{trie: t}
}

func (it *TrieIterator) workNode(currentNode *Value) {
	if currentNode.Size() == 2 {
		k := CompactDecode(currentNode.Get(0).Str())

		if currentNode.Get(1).Str() == "" {
			it.workNode(currentNode.Get(1))
		} else {
			if k[len(k)-1] == 16 {
				it.values = append(it.values, currentNode.Get(1).Str())
			} else {
				it.shas = append(it.shas, currentNode.Get(1).Bytes())
				it.getNode(currentNode.Get(1).Bytes())
			}
		}
	} else {
		for i := 0; i < currentNode.Size(); i++ {
			if i == 16 && currentNode.Get(i).Size() != 0 {
				it.values = append(it.values, currentNode.Get(i).Str())
			} else {
				if currentNode.Get(i).Str() == "" {
					it.workNode(currentNode.Get(i))
				} else {
					val := currentNode.Get(i).Str()
					if val != "" {
						it.shas = append(it.shas, currentNode.Get(1).Bytes())
						it.getNode([]byte(val))
					}
				}
			}
		}
	}
}

func (it *TrieIterator) getNode(node []byte) {
	currentNode := it.trie.cache.Get(node)
	it.workNode(currentNode)
}

// Collect ...
func (it *TrieIterator) Collect() [][]byte {
	if it.trie.Root == "" {
		return nil
	}

	it.getNode(NewValue(it.trie.Root).Bytes())

	return it.shas
}

// Purge ...
func (it *TrieIterator) Purge() int {
	shas := it.Collect()
	for _, sha := range shas {
		it.trie.cache.Delete(sha)
	}

	return len(it.values)
}

// Key ...
func (it *TrieIterator) Key() string {
	return ""
}

// Value ...
func (it *TrieIterator) Value() string {
	return ""
}

// EachCallback ...
type EachCallback func(key string, node *Value)

// Each ...
func (it *TrieIterator) Each(cb EachCallback) {
	it.fetchNode(nil, NewValue(it.trie.Root).Bytes(), cb)
}

func (it *TrieIterator) fetchNode(key []int, node []byte, cb EachCallback) {
	it.iterateNode(key, it.trie.cache.Get(node), cb)
}

func (it *TrieIterator) iterateNode(key []int, currentNode *Value, cb EachCallback) {
	if currentNode.Size() == 2 {
		k := CompactDecode(currentNode.Get(0).Str())

		pk := append(key, k...)
		if currentNode.Get(1).Size() != 0 && currentNode.Get(1).Str() == "" {
			it.iterateNode(pk, currentNode.Get(1), cb)
		} else {
			if k[len(k)-1] == 16 {
				cb(DecodeCompact(pk), currentNode.Get(1))
			} else {
				it.fetchNode(pk, currentNode.Get(1).Bytes(), cb)
			}
		}
	} else {
		for i := 0; i < currentNode.Size(); i++ {
			pk := append(key, i)
			if i == 16 && currentNode.Get(i).Size() != 0 {
				cb(DecodeCompact(pk), currentNode.Get(i))
			} else {
				if currentNode.Get(i).Size() != 0 && currentNode.Get(i).Str() == "" {
					it.iterateNode(pk, currentNode.Get(i), cb)
				} else {
					val := currentNode.Get(i).Str()
					if val != "" {
						it.fetchNode(pk, []byte(val), cb)
					}
				}
			}
		}
	}
}
