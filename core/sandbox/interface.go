package sandbox

// Interface ...
type Interface interface {
	Play(config *PlayConfig) ([]byte, error)
	CommitPlay(config *PlayConfig) (string, error)
}
