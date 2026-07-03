package cache

import "strings"

// MatchPattern builds a match function for InvalidateFunc from a pattern
// string per SPEC/03-cache-strategy.md: "*" matches everything, "method*"
// prefix-matches, "a,b,c" is a comma-separated list of exact/prefix terms,
// and a bare method name matches exactly.
func MatchPattern(pattern string) func(method, key string) bool {
	if pattern == "" || pattern == "*" {
		return func(string, string) bool { return true }
	}
	terms := strings.Split(pattern, ",")
	for i, t := range terms {
		terms[i] = strings.TrimSpace(t)
	}
	return func(method, _ string) bool {
		for _, t := range terms {
			if t == "*" {
				return true
			}
			if strings.HasSuffix(t, "*") {
				if strings.HasPrefix(method, strings.TrimSuffix(t, "*")) {
					return true
				}
				continue
			}
			if method == t {
				return true
			}
		}
		return false
	}
}
