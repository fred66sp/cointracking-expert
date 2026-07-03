package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// Key computes a deterministic cache key for (method, params), per
// SPEC/03-cache-strategy.md: params sorted alphabetically, empty values
// omitted, then sha256-hashed.
func Key(method string, params map[string]string) string {
	names := make([]string, 0, len(params))
	for k, v := range params {
		if v == "" {
			continue
		}
		names = append(names, k)
	}
	sort.Strings(names)

	var sb strings.Builder
	sb.WriteString(method)
	for _, k := range names {
		sb.WriteByte('|')
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(strings.ToLower(params[k]))
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}
