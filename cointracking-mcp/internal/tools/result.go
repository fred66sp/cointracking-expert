package tools

import (
	"bytes"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/alfredogp/cointracking-mcp/internal/api"
)

// jsonResult pretty-prints raw as MCP text content, mirroring
// utils/formatting.ts's formatJson + textResult.
func jsonResult(raw json.RawMessage) (*mcp.CallToolResult, any, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		buf.Reset()
		buf.Write(raw) // fall back to raw bytes if not re-indentable
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: buf.String()}},
	}, nil, nil
}

// errResult reports a tool-level error as text content with IsError=true,
// per reference-js-repo.md's formatToolError (never as an MCP protocol error,
// so the agent can see and self-correct).
func errResult(err error) (*mcp.CallToolResult, any, error) {
	msg := "Error: " + err.Error()
	if ctErr, ok := err.(*api.CTError); ok {
		msg = "CoinTracking API error: " + ctErr.Message
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}, nil, nil
}
