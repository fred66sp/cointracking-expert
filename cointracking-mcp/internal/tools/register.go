package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers the 6 CoinTracking data tools plus the 3 cache/
// project control tools on server, bound to app's shared state.
func RegisterAll(server *mcp.Server, app *App) {
	mcp.AddTool(server, getTradesTool(), getTradesHandler(app))
	mcp.AddTool(server, getBalanceTool(), getBalanceHandler(app))
	mcp.AddTool(server, getGroupedBalanceTool(), getGroupedBalanceHandler(app))
	mcp.AddTool(server, getHistoricalSummaryTool(), getHistoricalSummaryHandler(app))
	mcp.AddTool(server, getHistoricalCurrencyTool(), getHistoricalCurrencyHandler(app))
	mcp.AddTool(server, getGainsTool(), getGainsHandler(app))

	mcp.AddTool(server, invalidateCacheTool(), invalidateCacheHandler(app))
	mcp.AddTool(server, cacheStatsTool(), cacheStatsHandler(app))
	mcp.AddTool(server, closeProjectTool(), closeProjectHandler(app))
	mcp.AddTool(server, switchProjectTool(), switchProjectHandler(app))
}
