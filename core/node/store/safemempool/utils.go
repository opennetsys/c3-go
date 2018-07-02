package safemempool

import "fmt"

func buildKey(hash string) string {
	return fmt.Sprintf("tx_%s", hash)
}

func buildKeys(hashes []string) []string {
	var keys []string
	for _, hash := range hashes {
		keys = append(keys, fmt.Sprintf("tx_%s", hash))
	}

	return keys
}
