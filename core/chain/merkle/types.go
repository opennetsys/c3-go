package merkle

// Tree ...
type Tree struct {
	props TreeProps
}

// TreeProps ...
type TreeProps struct {
	MerkleTreeRootHash *string
	Kind               string
	Hashes             []string
}
