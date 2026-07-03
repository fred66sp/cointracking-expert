package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CloseProjectIn accepts an optional project_name for symmetry with
// SPEC/04-features.md, but this server instance only ever manages its own
// configured project (one project per process, per SPEC 01/03), so a
// mismatched name is rejected rather than silently ignored.
type CloseProjectIn struct {
	ProjectName string `json:"project_name,omitempty" jsonschema:"Project to close. If omitted, closes the current project. Must match this server's configured project."`
}

type CloseProjectOut struct {
	Project          string `json:"project"`
	CachedEntries    int    `json:"cached_entries"`
	EntriesPersisted int    `json:"entries_persisted"`
	Message          string `json:"message"`
}

func closeProjectTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_close_project",
		Description: "Signals that the audit of the current project is finished: flushes any pending disk " +
			"writes, clears the in-memory (L1) cache, and logs final stats. Call this when wrapping up an " +
			"audit/tax-prep session for a project so memory is freed; the cache reloads from disk next time.",
		Annotations: &mcp.ToolAnnotations{Title: "Close Project"},
	}
}

func closeProjectHandler(app *App) func(context.Context, *mcp.CallToolRequest, CloseProjectIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in CloseProjectIn) (*mcp.CallToolResult, any, error) {
		project := app.Project()
		if in.ProjectName != "" && in.ProjectName != project {
			return errResult(validationError(
				"este servidor gestiona el proyecto %q; no puede cerrar %q (una instancia = un proyecto)",
				project, in.ProjectName))
		}

		store := app.Store()
		cachedEntries := app.CacheManager().L1.Clear()
		store.Flush() // ensure pending async disk writes land before reporting/exiting
		persisted, err := store.Count()
		if err != nil {
			app.Log.Warnf("close_project: error counting persisted entries: %s", err)
		}

		app.Log.Infof("Proyecto %q cerrado. %d entradas en memoria liberadas, %d persistidas en disco.",
			project, cachedEntries, persisted)

		out := CloseProjectOut{
			Project:          project,
			CachedEntries:    cachedEntries,
			EntriesPersisted: persisted,
			Message: fmt.Sprintf("Proyecto cerrado. Caché persistida en %s. LRU limpiada de memoria.",
				app.CacheDir()+"/"+project),
		}
		raw, _ := json.Marshal(out)
		return jsonResult(raw)
	}
}
