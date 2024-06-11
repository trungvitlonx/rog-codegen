package util

import "sort"

// SortedMapKeys takes a map with keys of type string and returns a slice of those keys sorted lexicographically.
func SortedMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
