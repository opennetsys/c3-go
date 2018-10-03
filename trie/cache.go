package trie

import (
	hashutil "github.com/c3systems/c3-go/common/hashutil"
)

// Cache ...
type Cache struct {
	nodes   map[string]*Node
	db      Database
	IsDirty bool
}

// NewCache ...
func NewCache(db Database) *Cache {
	return &Cache{
		db:    db,
		nodes: make(map[string]*Node),
	}
}

// Put ...
func (cache *Cache) Put(v interface{}) interface{} {
	return cache.putValue(v, false)
}

// putValue ...
func (cache *Cache) putValue(v interface{}, force bool) interface{} {
	value := NewValue(v)

	enc := value.Encode()
	if len(enc) >= 32 || force {
		sha := cache.hashBytes(enc)

		cache.nodes[string(sha)] = NewNode(sha, value, true)
		cache.IsDirty = true

		return sha
	}

	return v
}

// Get ...
func (cache *Cache) Get(key []byte) *Value {
	if cache.nodes[string(key)] != nil {
		return cache.nodes[string(key)].Value
	}

	data, _ := cache.db.Get(key)
	value := NewValueFromBytes(data)
	cache.nodes[string(key)] = NewNode(key, value, false)

	return value
}

// Delete ...
func (cache *Cache) Delete(key []byte) {
	delete(cache.nodes, string(key))

	cache.db.Delete(key)
}

// Commit ...
func (cache *Cache) Commit() {
	if !cache.IsDirty {
		return
	}

	for key, node := range cache.nodes {
		if node.Dirty {
			cache.db.Put([]byte(key), node.Value.Encode())
			node.Dirty = false
		}
	}

	cache.IsDirty = false

	// empty cache nodes if exceeds size limit
	if len(cache.nodes) > 200 {
		cache.nodes = make(map[string]*Node)
	}
}

// Undo ...
func (cache *Cache) Undo() {
	for key, node := range cache.nodes {
		if node.Dirty {
			delete(cache.nodes, key)
		}
	}

	cache.IsDirty = false
}

func (cache *Cache) hashBytes(v []byte) []byte {
	h := hashutil.Hash(v)
	return h[:]
}
