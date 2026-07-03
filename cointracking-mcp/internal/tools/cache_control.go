package tools

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/alfredogp/cointracking-mcp/internal/cache"
)

// --- cointracking_invalidate_cache ---

type InvalidateCacheIn struct {
	Pattern string `json:"pattern,omitempty" jsonschema:"Pattern of entries to invalidate: a method name, 'method*' prefix, comma-separated list, or '*' for everything. Use after the user modifies data in the CoinTracking web UI, so stale cached data isn't reused."`
}

type InvalidateCacheOut struct {
	Invalidated int `json:"invalidated"`
}

func invalidateCacheTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_invalidate_cache",
		Description: "Invalidates cached entries matching a pattern (e.g. 'getTrades*', '*' for everything). " +
			"Call this after guiding the user to edit/delete/add/reimport data in CoinTracking's web UI, before " +
			"re-querying — otherwise stale cached figures could be mixed with fresh ones.",
		Annotations: &mcp.ToolAnnotations{Title: "Invalidate Cache"},
	}
}

func invalidateCacheHandler(app *App) func(context.Context, *mcp.CallToolRequest, InvalidateCacheIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in InvalidateCacheIn) (*mcp.CallToolResult, any, error) {
		removed, err := app.CacheManager().Invalidate(cache.MatchPattern(in.Pattern))
		if err != nil {
			return errResult(err)
		}
		app.Log.Infof("invalidateCache(%q) -> %d entries removed", in.Pattern, removed)
		raw, _ := json.Marshal(InvalidateCacheOut{Invalidated: removed})
		return jsonResult(raw)
	}
}

// --- cointracking_cache_stats ---

type CacheStatsIn struct{}

type methodCallStat struct {
	Method   string `json:"method"`
	Count    int    `json:"count"`
	LastCall string `json:"last_call"`
}

type CacheStatsOut struct {
	Size         int              `json:"size"`
	MaxSize      int              `json:"max_size"`
	Hits         int64            `json:"hits"`
	Misses       int64            `json:"misses"`
	TotalCalls   int64            `json:"total_calls"`
	HitRate      float64          `json:"hit_rate"`
	RateLimit    int              `json:"rate_limit_per_hour"`
	UsedThisHour int              `json:"used_this_hour"`
	CallsToAPI   []methodCallStat `json:"calls_to_api"`
}

func cacheStatsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_cache_stats",
		Description: "Returns cache and API-usage statistics: cache size/hit-rate, and calls made to the " +
			"CoinTracking API by method with their timestamps and current hourly rate-limit consumption. Useful " +
			"to decide whether re-querying is worth it or whether cached data should be reused.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, Title: "Cache Stats"},
	}
}

func cacheStatsHandler(app *App) func(context.Context, *mcp.CallToolRequest, CacheStatsIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, _ CacheStatsIn) (*mcp.CallToolResult, any, error) {
		s := app.CacheManager().L1.Stats()
		total := s.Hits + s.Misses
		var hitRate float64
		if total > 0 {
			hitRate = float64(s.Hits) / float64(total)
		}

		methodStats := app.Rate.MethodStats()
		calls := make([]methodCallStat, 0, len(methodStats))
		for _, m := range methodStats {
			calls = append(calls, methodCallStat{Method: m.Method, Count: m.Count, LastCall: m.LastCall.UTC().Format("2006-01-02T15:04:05Z")})
		}

		out := CacheStatsOut{
			Size: s.Size, MaxSize: s.MaxSize, Hits: s.Hits, Misses: s.Misses,
			TotalCalls: total, HitRate: hitRate,
			RateLimit: app.Rate.Limit(), UsedThisHour: app.Rate.UsedInWindow(),
			CallsToAPI: calls,
		}
		raw, _ := json.Marshal(out)
		return jsonResult(raw)
	}
}
