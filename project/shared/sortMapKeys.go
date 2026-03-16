package shared

import (
	"cmp"
	"iter"
	"slices"
)

func SortedMap[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		keys := make([]K, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}
