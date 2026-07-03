package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetHistoricalSummaryIn struct {
	BTC          int    `json:"btc,omitempty" jsonschema:"If 1, return values in BTC. If 0 (default), return values in account fiat currency."`
	Start        int64  `json:"start,omitempty" jsonschema:"Start time as UNIX timestamp in SECONDS."`
	End          int64  `json:"end,omitempty" jsonschema:"End time as UNIX timestamp in SECONDS."`
	FiatCurrency string `json:"fiat_currency,omitempty" jsonschema:"Fiat currency for values, e.g. EUR, USD."`
}

func getHistoricalSummaryTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_get_historical_summary",
		Description: "Returns a historical summary of the entire portfolio aggregated by year/month. Useful " +
			"for tracking portfolio value evolution over time and for tax-year breakdowns.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, Title: "Get Historical Portfolio Summary"},
	}
}

func getHistoricalSummaryHandler(app *App) func(context.Context, *mcp.CallToolRequest, GetHistoricalSummaryIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in GetHistoricalSummaryIn) (*mcp.CallToolResult, any, error) {
		if in.BTC != 0 && in.BTC != 1 {
			return errResult(validationError("btc debe ser 0 o 1, recibido %d", in.BTC))
		}

		p := params().
			intp("btc", in.BTC).
			int64p("start", in.Start).
			int64p("end", in.End).
			str("fiat_currency", in.FiatCurrency)

		raw, _, err := cachedCall(app, "getHistoricalSummary", p)
		if err != nil {
			return errResult(err)
		}
		return jsonResult(raw)
	}
}
