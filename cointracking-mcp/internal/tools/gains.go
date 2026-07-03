package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetGainsIn struct {
	Price string `json:"price,omitempty" jsonschema:"Cost-basis selection method for unrealized gains: best, worst, oldest (FIFO), or newest (LIFO)."`
	BTC   int    `json:"btc,omitempty" jsonschema:"If 1, return values in BTC. If 0 (default), return values in account fiat currency."`
}

func getGainsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_get_gains",
		Description: "Returns realized and unrealized gains for the account. Useful for tax reporting and " +
			"P&L analysis. Choose the cost-basis method via price (FIFO/LIFO/best/worst).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, Title: "Get Gains"},
	}
}

func getGainsHandler(app *App) func(context.Context, *mcp.CallToolRequest, GetGainsIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in GetGainsIn) (*mcp.CallToolResult, any, error) {
		if in.Price != "" {
			switch in.Price {
			case "best", "worst", "oldest", "newest":
			default:
				return errResult(validationError("price debe ser \"best\", \"worst\", \"oldest\" o \"newest\", recibido %q", in.Price))
			}
		}
		if in.BTC != 0 && in.BTC != 1 {
			return errResult(validationError("btc debe ser 0 o 1, recibido %d", in.BTC))
		}

		p := params().
			str("price", in.Price).
			intp("btc", in.BTC)

		raw, _, err := cachedCall(app, "getGains", p)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(raw)
	}
}
