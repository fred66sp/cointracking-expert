package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetHistoricalCurrencyIn struct {
	Currency     string `json:"currency" jsonschema:"Ticker of the currency to query (e.g. BTC, ETH, USDC)."`
	Start        int64  `json:"start,omitempty" jsonschema:"Start time as UNIX timestamp in SECONDS."`
	End          int64  `json:"end,omitempty" jsonschema:"End time as UNIX timestamp in SECONDS."`
	FiatCurrency string `json:"fiat_currency,omitempty" jsonschema:"Fiat currency for values, e.g. EUR, USD."`
}

func getHistoricalCurrencyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_get_historical_currency",
		Description: "Returns the historical balance and value of a single specific currency in the account " +
			"over time. Use this to chart the holdings and valuation of one coin (e.g. how BTC balance and " +
			"value evolved year by year).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, Title: "Get Historical Currency"},
	}
}

func getHistoricalCurrencyHandler(app *App) func(context.Context, *mcp.CallToolRequest, GetHistoricalCurrencyIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in GetHistoricalCurrencyIn) (*mcp.CallToolResult, any, error) {
		if in.Currency == "" {
			return errResult(validationError("currency es obligatorio"))
		}

		p := params().
			str("currency", in.Currency).
			int64p("start", in.Start).
			int64p("end", in.End).
			str("fiat_currency", in.FiatCurrency)

		raw, _, err := cachedCall(app, "getHistoricalCurrency", p)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(raw)
	}
}
