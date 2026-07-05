package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/alfredogp/cointracking-mcp/internal/config"
)

type DeleteProjectIn struct {
	ProjectName string `json:"project_name" jsonschema:"Project whose cache directory to delete permanently. Must not be the currently active project."`
}

type DeleteProjectOut struct {
	Project string `json:"project"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

func deleteProjectTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cointracking_delete_project",
		Description: "Permanently deletes a project's on-disk cache directory (cache-dir/<project>). The MCP " +
			"server owns its own cache files, so use this instead of deleting cache-dir/<project> by hand: an " +
			"external rm can race the server's SQLite handle and fail with 'file busy' on Windows. Refuses to " +
			"delete the currently active project — switch away first (cointracking_switch_project).",
		Annotations: &mcp.ToolAnnotations{Title: "Delete Project Cache", DestructiveHint: boolPtr(true)},
	}
}

func deleteProjectHandler(app *App) func(context.Context, *mcp.CallToolRequest, DeleteProjectIn) (*mcp.CallToolResult, any, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, in DeleteProjectIn) (*mcp.CallToolResult, any, error) {
		if in.ProjectName == "" {
			return errResult(validationError("project_name es obligatorio"))
		}
		if err := config.ValidateProjectName(in.ProjectName); err != nil {
			return errResult(err)
		}
		projectDir := filepath.Join(app.CacheDir(), in.ProjectName)
		if _, err := os.Stat(projectDir); os.IsNotExist(err) {
			return errResult(validationError("no existe caché para el proyecto %q en %s", in.ProjectName, projectDir))
		}

		// The active-project check and the delete happen under the same lock
		// SwitchProject uses (WithProjectLockedIfNotActive), so a concurrent
		// switch to in.ProjectName can't land in between and get its cache
		// pulled out from under it (TOCTOU fixed 2026-07-05).
		err := app.WithProjectLockedIfNotActive(in.ProjectName, func() error {
			return removeAllWithRetry(projectDir)
		})
		if err == ErrProjectActive {
			return errResult(validationError(
				"%q es el proyecto activo; cambia a otro con cointracking_switch_project antes de borrarlo",
				in.ProjectName))
		}
		if err != nil {
			return errResult(fmt.Errorf("borrando caché de %q: %w", in.ProjectName, err))
		}

		app.Log.Infof("Proyecto %q borrado (caché en %s eliminada)", in.ProjectName, projectDir)

		out := DeleteProjectOut{
			Project: in.ProjectName,
			Path:    projectDir,
			Message: fmt.Sprintf("Caché de %q borrada de %s.", in.ProjectName, projectDir),
		}
		raw, _ := json.Marshal(out)
		return jsonResult(raw)
	}
}

// removeAllWithRetry retries os.RemoveAll a few times with backoff: on
// Windows, a file an in-process SQLite handle just released can stay
// transiently locked (AV/indexer scan) for a few hundred ms after Close().
func removeAllWithRetry(dir string) error {
	var err error
	for attempt, delay := 0, 25*time.Millisecond; attempt < 6; attempt++ {
		if err = os.RemoveAll(dir); err == nil {
			return nil
		}
		time.Sleep(delay)
		delay *= 2
	}
	return err
}

func boolPtr(b bool) *bool { return &b }
