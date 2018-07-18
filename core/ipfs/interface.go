package ipfs

// Interface ...
type Interface interface {
	Get(hash, outdir string) error
	AddDir(dir string) (string, error)
	Refs(hash string, recursive bool) (<-chan string, error)
}
