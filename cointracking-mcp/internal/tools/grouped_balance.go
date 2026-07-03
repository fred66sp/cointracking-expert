package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetGroupedBalanceIn follows reference-js-repo.md's tools/get-grouped-balance.ts,
// which requires `group`. (SPEC/02-api-mapping.md lists this tool as
// parameterless, but that contradicts both the real CoinTracking API and the
// JS reference implementation being ported; the reference implementation
// takes precedence here.)
type GetGroupedBalanceIn struct {
	Group          string `json:"group" jsonschema:"Grouping dimension: exchange, type, or currency."`
	ExcludeDepWith int    `json:"exclude_dep_with,omitempty" jsonschema:"If 1, exclude pure deposits/withdrawals from the grouping."`
	Type           string `json:"type,omitempty" jsonschema:"Optional transaction type filter (e.g. 'Trade', 'Staking', 'Airdrop')."`
}

func getGroupedBalanceTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_get_grouped_balance",
		Description: "Returns balances grouped by exchange, type, or currency. Useful to see, for example, " +
			"total holdings per exchange, or distribution by transaction type (staking vs trades vs airdrops).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, Title: "Get Grouped Balance"},
	}
}

func getGroupedBalanceHandler(app *App) func(context.Context, *mcp.CallToolRequest, GetGroupedBalanceIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in GetGroupedBalanceIn) (*mcp.CallToolResult, any, error) {
		switch in.Group {
		case "exchange", "type", "currency":
		default:
			return errResult(validationError("group debe ser \"exchange\", \"type\" o \"currency\", recibido %q", in.Group))
		}
		if in.ExcludeDepWith != 0 && in.ExcludeDepWith != 1 {
			return errResult(validationError("exclude_dep_with debe ser 0 o 1, recibido %d", in.ExcludeDepWith))
		}

		p := params().
			str("group", in.Group).
			intp("exclude_dep_with", in.ExcludeDepWith).
			str("type", in.Type)

		raw, _, err := cachedCall(app, "getGroupedBalance", p)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(raw)
	}
}
