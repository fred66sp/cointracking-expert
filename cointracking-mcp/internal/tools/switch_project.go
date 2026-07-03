package tools

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SwitchProjectIn struct {
	ProjectName string `json:"project_name" jsonschema:"Project to make active. Alphanumeric, '_' and '-' only. Isolates cache under cache-dir/<project>, same as --project at startup."`
}

func switchProjectTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_switch_project",
		Description: "Makes a different project active, live, without restarting the MCP server: flushes and " +
			"closes the current project's cache, then opens (or creates) the target project's cache under the " +
			"same cache-dir, isolated per SPEC/03-cache-strategy.md. Call this whenever the conversation's active " +
			"project changes (ADR-013), before any cointracking_get_* call, so data from different projects is " +
			"never mixed. Credentials, account tier, and the hourly rate limit are per-process and unaffected.",
		Annotations: &mcp.ToolAnnotations{Title: "Switch Project"},
	}
}

func switchProjectHandler(app *App) func(context.Context, *mcp.CallToolRequest, SwitchProjectIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in SwitchProjectIn) (*mcp.CallToolResult, any, error) {
		if in.ProjectName == "" {
			return errResult(validationError("project_name es obligatorio"))
		}

		result, err := app.SwitchProject(in.ProjectName)
		if err != nil {
			return errResult(err)
		}

		raw, _ := json.Marshal(result)
		return jsonResult(raw)
	}
}
