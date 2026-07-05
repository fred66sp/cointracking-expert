package tools

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/alfredogp/cointracking-mcp/internal/cache"
)

// cachedCall implements the read cascade from SPEC/03-cache-strategy.md:
// memory (L1) -> disk (L2) -> CoinTracking API. On a miss, it calls the API,
// caches the raw JSON response, and returns it. Errors are never cached.
func cachedCall(app *App, method string, params map[string]string) (json.RawMessage, cache.Source, error) {
	key := cache.Key(method, params)

	mgr := app.CacheManager()
	if entry, src, ok := mgr.Lookup(key); ok {
		app.Log.Debugf("%s(%v) -> CACHE HIT (%s, expires %s)", method, params, src, entry.ExpiresAt.Format(time.RFC3339))
		return json.RawMessage(entry.Value), src, nil
	}

	start := time.Now()
	raw, err := app.Client().RequestRaw(method, params)
	if err != nil {
		app.Log.Debugf("%s(%v) -> CACHE MISS -> API CALL FAILED (%s): %s", method, params, time.Since(start), err)
		return nil, "", err
	}
	app.Log.Debugf("%s(%v) -> CACHE MISS -> API CALL (%s) -> cached", method, params, time.Since(start))

	ttl := TTLFor(method)
	mgr.Store(key, method, raw, ttl, func(err error) {
		app.Log.Warnf("cache: async disk write failed for %s: %s", method, err)
	})

	return raw, cache.FromAPI, nil
}

// strParams converts a set of optional values into cache/API params,
// omitting zero values so cache keys stay normalized (SPEC 03 normalization).
type paramBuilder map[string]string

func params() paramBuilder { return paramBuilder{} }

func (p paramBuilder) str(k, v string) paramBuilder {
	if v != "" {
		p[k] = v
	}
	return p
}

func (p paramBuilder) intp(k string, v int) paramBuilder {
	if v != 0 {
		p[k] = strconv.Itoa(v)
	}
	return p
}

func (p paramBuilder) int64p(k string, v int64) paramBuilder {
	if v != 0 {
		p[k] = strconv.FormatInt(v, 10)
	}
	return p
}
