package trie

// ValueIterator ...
type ValueIterator struct {
	value        *Value
	currentValue *Value
	idx          int
}

// NewIterator ...
func (val *Value) NewIterator() *ValueIterator {
	return &ValueIterator{value: val}
}

// Next ...
func (it *ValueIterator) Next() bool {
	if it.idx >= it.value.Size() {
		return false
	}

	it.currentValue = it.value.Get(it.idx)
	it.idx++

	return true
}

// Value ...
func (it *ValueIterator) Value() *Value {
	return it.currentValue
}

// Idx ...
func (it *ValueIterator) Idx() int {
	return it.idx
}
