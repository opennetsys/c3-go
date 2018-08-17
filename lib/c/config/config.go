package main

import "C"
import (
	"github.com/c3systems/c3-go/config"
)

// ServerHost ...
//export ServerHost
func ServerHost() *C.char {
	return C.CString(config.ServerHost)
}

// ServerPort ...
//export ServerPort
func ServerPort() C.int {
	return C.int(config.ServerPort)
}

// DefaultStoreDirectory is the default directory where the file system store will live.
//export DefaultStoreDirectory
func DefaultStoreDirectory() *C.char {
	return C.CString(config.DefaultStoreDirectory)
}

// TempContainerStatePath ...
//export TempContainerStatePath
func TempContainerStatePath() *C.char {
	return C.CString(config.TempContainerStatePath)
}

// TempContainerStateFileName ...
//export TempContainerStateFileName
func TempContainerStateFileName() *C.char {
	return C.CString(config.TempContainerStateFileName)
}

// TempContainerStateFilePath ...
//export TempContainerStateFilePath
func TempContainerStateFilePath() *C.char {
	return C.CString(config.TempContainerStateFilePath)
}

// DockerRegistryPort ...
//export DockerRegistryPort
func DockerRegistryPort() C.int {
	return C.int(config.DockerRegistryPort)
}

// DefaultBlockDifficulty ...
//export DefaultBlockDifficulty
func DefaultBlockDifficulty() C.int {
	return C.int(config.DefaultBlockDifficulty)
}

// MinedBlockVerificationTimeout in seconds
//export MinedBlockVerificationTimeout
func MinedBlockVerificationTimeout() C.double {
	return C.double(config.MinedBlockVerificationTimeout.Seconds())
}

// IPFSTimeout ...
//export IPFSTimeout
func IPFSTimeout() C.double {
	return C.double(config.IPFSTimeout.Seconds())
}

func main() {}
