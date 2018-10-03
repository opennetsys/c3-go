package trie

// CompareIntSlice is a helper function for comparing int slices
func CompareIntSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// MatchingNibbleLength returns the amount of nibbles that match each other from 0 ...
func MatchingNibbleLength(a, b []int) int {
	i := 0
	for CompareIntSlice(a[:i+1], b[:i+1]) && i < len(b) {
		i++
	}

	return i
}
