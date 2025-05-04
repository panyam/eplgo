package epl

import (
	"sort"
)

func SortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Requires "sort" import
	return keys
}

func StringListEq(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, s := range s1 {
		if s2[i] != s {
			return false
		}
	}
	return true
}

func StringSetEq(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	found := map[string]bool{}
	for _, s := range s1 {
		found[s] = true
	}
	for _, s := range s2 {
		if !found[s] {
			return false
		}
	}
	return true
}

func DictZip[K comparable, V any](keys []K, vals []V) map[K]V {
	out := map[K]V{}
	l1 := len(keys)
	l2 := len(vals)
	for i := range min(l1, l2) {
		out[keys[i]] = vals[i]
	}
	return out
}

func Dict[K comparable, V any](args ...any) (out map[K]V) {
	out = map[K]V{}
	for i := 0; i < len(args); i += 2 {
		k := args[i].(K)
		v := args[i+1].(V)
		out[k] = v
	}
	return out
}
