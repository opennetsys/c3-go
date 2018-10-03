package state

// StateObject ...
type StateObject interface {
	State() *Trie
	Sync()
	Undo()
}
