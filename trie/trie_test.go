package trie

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

const longWord = "1234567890abcdefghijklmnopqrstuvwxxzABCEFGHIJKLMNOPQRSTUVWXYZ"

type MemDatabase struct {
	db map[string][]byte
}

func NewMemDatabase() (*MemDatabase, error) {
	db := &MemDatabase{db: make(map[string][]byte)}
	return db, nil
}
func (db *MemDatabase) Put(key []byte, value []byte) {
	db.db[string(key)] = value
}
func (db *MemDatabase) Get(key []byte) ([]byte, error) {
	return db.db[string(key)], nil
}
func (db *MemDatabase) Delete(key []byte) error {
	delete(db.db, string(key))
	return nil
}
func (db *MemDatabase) Print()              {}
func (db *MemDatabase) Close()              {}
func (db *MemDatabase) LastKnownTD() []byte { return nil }

func New() (*MemDatabase, *Trie) {
	db, _ := NewMemDatabase()
	return db, NewTrie(db, "")
}

func TestTrieSync(t *testing.T) {
	db, trie := New()
	trie.Update("foo", longWord)
	if len(db.db) != 0 {
		t.Error("Expected no data in database")
	}
	trie.Sync()
	if len(db.db) == 0 {
		t.Error("Expected data to be persisted")
	}
}

func TestTrieDirtyTracking(t *testing.T) {
	_, trie := New()
	trie.Update("foo", longWord)
	if !trie.cache.IsDirty {
		t.Error("Expected trie to be dirty")
	}
	trie.Sync()
	if trie.cache.IsDirty {
		t.Error("Expected trie not to be dirty")
	}
	trie.Update("test", longWord)
	trie.cache.Undo()
	if trie.cache.IsDirty {
		t.Error("Expected trie not to be dirty")
	}
}

func TestTrieReset(t *testing.T) {
	_, trie := New()
	trie.Update("foo", longWord)
	if len(trie.cache.nodes) == 0 {
		t.Error("Expected cached nodes")
	}
	trie.cache.Undo()
	if len(trie.cache.nodes) != 0 {
		t.Error("Expected no nodes after undo")
	}
}

func TestTrieGet(t *testing.T) {
	_, trie := New()
	trie.Update("foo", longWord)
	x := trie.Get("foo")
	if x != longWord {
		t.Errorf("expected %s, got %s", longWord, x)
	}
}

func TestTrieUpdating(t *testing.T) {
	_, trie := New()
	trie.Update("foo", longWord)
	trie.Update("foo", longWord+"1")
	x := trie.Get("foo")
	if x != longWord+"1" {
		t.Errorf("expected %s, got %s", longWord+"1", x)
	}
}

func TestTrieCmp(t *testing.T) {
	_, trie1 := New()
	_, trie2 := New()
	trie1.Update("foo", longWord)
	trie2.Update("foo", longWord)
	if !trie1.Cmp(trie2) {
		t.Error("Expected tries to be equal")
	}
	trie1.Update("foo", longWord)
	trie2.Update("bar", longWord)
	if trie1.Cmp(trie2) {
		t.Errorf("Expected tries not to be equal %x %x", trie1.Root, trie2.Root)
	}
}

func TestTrieDelete(t *testing.T) {
	_, trie := New()
	trie.Update("foo", longWord)
	exp := trie.Root
	trie.Update("bar", longWord)
	trie.Delete("bar")
	if !reflect.DeepEqual(exp, trie.Root) {
		t.Errorf("Expected tries to be equal %x : %x", exp, trie.Root)
	}
	trie.Update("bar", longWord)
	exp = trie.Root
	trie.Update("qux", longWord)
	trie.Delete("qux")
	if !reflect.DeepEqual(exp, trie.Root) {
		t.Errorf("Expected tries to be equal %x : %x", exp, trie.Root)
	}
}

func TestTrieDeleteWithValue(t *testing.T) {
	_, trie := New()
	trie.Update("f", longWord)
	exp := trie.Root
	trie.Update("fo", longWord)
	trie.Update("foo", longWord)
	trie.Delete("fo")
	trie.Delete("foo")
	if !reflect.DeepEqual(exp, trie.Root) {
		t.Errorf("Expected tries to be equal %x : %x", exp, trie.Root)
	}
}

func TestTriePurge(t *testing.T) {
	_, trie := New()
	trie.Update("f", longWord)
	trie.Update("fo", longWord)
	trie.Update("foo", longWord)
	lenBefore := len(trie.cache.nodes)
	it := trie.NewIterator()
	if num := it.Purge(); num != 3 {
		t.Errorf("Expected purge to return 3, got %d", num)
	}
	if lenBefore == len(trie.cache.nodes) {
		t.Errorf("Expected cached nodes to be deleted")
	}
}

func h(str string) string {
	d, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return string(d)
}

func get(in string) (out string) {
	if len(in) > 2 && in[:2] == "0x" {
		out = h(in[2:])
	} else {
		out = in
	}
	return
}

type Test struct {
	Name string
	In   map[string]string
	Root string
}

const MaxTest = 1000

func TestDelete(t *testing.T) {
	_, trie := New()
	trie.Update("a", "jeffreytestlongstring")
	trie.Update("aa", "otherstring")
	trie.Update("aaa", "othermorestring")
	trie.Update("aabbbbccc", "hithere")
	trie.Update("abbcccdd", "hstanoehutnaheoustnh")
	trie.Update("rnthaoeuabbcccdd", "hstanoehutnaheoustnh")
	trie.Update("rneuabbcccdd", "hstanoehutnaheoustnh")
	trie.Update("rneuabboeusntahoeucccdd", "hstanoehutnaheoustnh")
	trie.Update("rnxabboeusntahoeucccdd", "hstanoehutnaheoustnh")
	trie.Delete("aaboaestnuhbccc")
	trie.Delete("a")
	trie.Update("a", "nthaonethaosentuh")
	trie.Update("c", "shtaosntehua")
	trie.Delete("a")
	trie.Update("aaaa", "testmegood")
	fmt.Println("aa =>", trie.Get("aa"))
	_, t2 := New()
	trie.NewIterator().Each(func(key string, v *Value) {
		if key == "aaaa" {
			t2.Update(key, v.Str())
		} else {
			t2.Update(key, v.Str())
		}
	})
	a := NewValue(trie.Root).Bytes()
	b := NewValue(t2.Root).Bytes()
	fmt.Printf("o: %x\nc: %x\n", a, b)
}

func TestRndCase(t *testing.T) {
	_, trie := New()

	data := []struct{ k, v string }{
		{"0000000000000000000000000000000000000000000000000000000000000001", "a07573657264617461000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000003", "8453bb5b31"},
		{"0000000000000000000000000000000000000000000000000000000000000004", "850218711a00"},
		{"0000000000000000000000000000000000000000000000000000000000000005", "9462d7705bd0b3ecbc51a8026a25597cb28a650c79"},
		{"0000000000000000000000000000000000000000000000000000000000000010", "947e70f9460402290a3e487dae01f610a1a8218fda"},
		{"0000000000000000000000000000000000000000000000000000000000000111", "01"},
		{"0000000000000000000000000000000000000000000000000000000000000112", "a053656e6174650000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000113", "a053656e6174650000000000000000000000000000000000000000000000000000"},
		{"53656e6174650000000000000000000000000000000000000000000000000000", "94977e3f62f5e1ed7953697430303a3cfa2b5b736e"},
	}
	for _, e := range data {
		trie.Update(string(Hex2Bytes(e.k)), string(Hex2Bytes(e.v)))
	}

	fmt.Printf("root after update %x\n", trie.Root)
	trie.NewIterator().Each(func(k string, v *Value) {
		fmt.Printf("%x %x\n", k, v.Bytes())
	})

	data = []struct{ k, v string }{
		{"0000000000000000000000000000000000000000000000000000000000000112", ""},
		{"436974697a656e73000000000000000000000000000000000000000000000001", ""},
		{"436f757274000000000000000000000000000000000000000000000000000002", ""},
		{"53656e6174650000000000000000000000000000000000000000000000000000", ""},
		{"436f757274000000000000000000000000000000000000000000000000000000", ""},
		{"53656e6174650000000000000000000000000000000000000000000000000001", ""},
		{"0000000000000000000000000000000000000000000000000000000000000113", ""},
		{"436974697a656e73000000000000000000000000000000000000000000000000", ""},
		{"436974697a656e73000000000000000000000000000000000000000000000002", ""},
		{"436f757274000000000000000000000000000000000000000000000000000001", ""},
		{"0000000000000000000000000000000000000000000000000000000000000111", ""},
		{"53656e6174650000000000000000000000000000000000000000000000000002", ""},
	}

	for _, e := range data {
		trie.Delete(string(Hex2Bytes(e.k)))
	}

	fmt.Printf("root after delete %x\n", trie.Root)

	trie.NewIterator().Each(func(k string, v *Value) {
		fmt.Printf("%x %x\n", k, v.Bytes())
	})

	fmt.Printf("%x\n", trie.Get(string(Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"))))
}
