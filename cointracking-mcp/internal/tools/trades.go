package tools

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetTradesIn struct {
	Limit       int    `json:"limit,omitempty" jsonschema:"Maximum number of trades to return. Strongly recommended for large historical accounts."`
	Order       string `json:"order,omitempty" jsonschema:"Sort order by time: ASC or DESC. DESC (default) = newest first."`
	Start       int64  `json:"start,omitempty" jsonschema:"Start time as UNIX timestamp in SECONDS (not ms). Filters trades on or after this time."`
	End         int64  `json:"end,omitempty" jsonschema:"End time as UNIX timestamp in SECONDS (not ms). Filters trades on or before this time."`
	TradePrices int    `json:"trade_prices,omitempty" jsonschema:"If 1, include the trade value in BTC and the account fiat currency."`
}

func getTradesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_get_trades",
		Description: "Returns all trades and transactions from CoinTracking, including trades, deposits, " +
			"withdrawals, staking, mining, airdrops, DeFi, etc. Each entry includes buy/sell amounts and " +
			"currencies, fees, type, exchange, group, comment, imported_from, time, and trade_id. For accounts " +
			"with long history, ALWAYS narrow with start/end (UNIX seconds) and a limit to avoid timeouts and " +
			"rate-limit waste.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, Title: "Get Trades"},
	}
}

func getTradesHandler(app *App) func(context.Context, *mcp.CallToolRequest, GetTradesIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in GetTradesIn) (*mcp.CallToolResult, any, error) {
		if in.Order != "" {
			order := strings.ToUpper(in.Order)
			if order != "ASC" && order != "DESC" {
				return errResult(validationError("order debe ser \"ASC\" o \"DESC\", recibido %q", in.Order))
			}
			in.Order = order
		}
		if in.TradePrices != 0 && in.TradePrices != 1 {
			return errResult(validationError("trade_prices debe ser 0 o 1, recibido %d", in.TradePrices))
		}

		p := params().
			intp("limit", in.Limit).
			str("order", in.Order).
			int64p("start", in.Start).
			int64p("end", in.End).
			intp("trade_prices", in.TradePrices)

		raw, _, err := cachedCall(app, "getTrades", p)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(raw)
	}
}
