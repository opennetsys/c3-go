package hashing

import (
	"crypto/sha512"
	"hash"

	"github.com/c3systems/c3-go/common/hexutil"
)

// New ...
func New() hash.Hash {
	// note: 512_256 is a bit faster on x86_64 machines
	return sha512.New512_256()
}

// Hash ...
func Hash(data []byte) [sha512.Size256]byte {
	return sha512.Sum512_256(data)
}

// HashToHexString ...
func HashToHexString(data []byte) string {
	b := Hash(data)
	return hexutil.EncodeToString(b[:])
}

// IsEqual ...
func IsEqual(hexHash string, data []byte) bool {
	return hexHash == HashToHexString(data)
}
