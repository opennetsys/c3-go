package miner

// Interface ...
// TODO: finish
type Interface interface {
	Props() Props
	SpawnMiner() error
}
