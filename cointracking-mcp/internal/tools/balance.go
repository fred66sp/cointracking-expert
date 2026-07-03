package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetBalanceIn struct{}

func getBalanceTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_get_balance",
		Description: "Returns the current balance for every coin in the account. For each coin: currency, " +
			"total balance, amount on exchange, amount on wallet, BTC value, and fiat value in the account " +
			"currency. Use this for an instant portfolio overview.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, Title: "Get Current Balance"},
	}
}

func getBalanceHandler(app *App) func(context.Context, *mcp.CallToolRequest, GetBalanceIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, _ GetBalanceIn) (*mcp.CallToolResult, any, error) {
		raw, _, err := cachedCall(app, "getBalance", nil)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(raw)
	}
}
