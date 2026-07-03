package tools

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/alfredogp/cointracking-mcp/internal/api"
	"github.com/alfredogp/cointracking-mcp/internal/config"
	"github.com/alfredogp/cointracking-mcp/internal/logging"
)

// TestCachePersistsAcrossRestart exercises the real read path end-to-end
// (App -> cachedCall -> fake CT API) and verifies SPEC 04's close_project
// contract: after closing and reopening the app against the same cache
// dir/project, the previously-fetched value is served from disk without a
// second API call.
func TestCachePersistsAcrossRestart(t *testing.T) {
	apiCalls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":1,"balance":{"BTC":"1.5"}}`))
	}))
	defer srv.Close()
	restore := api.SetAPIURLForTesting(srv.URL)
	defer restore()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		APIKey: "k", APISecret: "s", Tier: "pro", Project: "itest",
		CacheSize: 100, CacheDir: cacheDir, LogLevel: "error", LogFormat: "text", Timezone: "UTC",
	}
	log := logging.New(&bytes.Buffer{}, logging.ParseLevel(cfg.LogLevel), false, nil)

	app1, err := NewApp(cfg, log)
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}
	raw1, src1, err := cachedCall(app1, "getBalance", nil)
	if err != nil {
		t.Fatalf("cachedCall (first): %v", err)
	}
	if src1 != "FROM_API" {
		t.Fatalf("expected first call FROM_API, got %s", src1)
	}
	if apiCalls != 1 {
		t.Fatalf("expected 1 API call, got %d", apiCalls)
	}

	// close_project: flush L1, ensure disk has the entry.
	cleared := app1.CacheManager().L1.Clear()
	if cleared != 1 {
		t.Fatalf("expected 1 entry cleared from L1, got %d", cleared)
	}
	if err := app1.Close(); err != nil {
		t.Fatalf("app1.Close: %v", err)
	}

	// "Restart": open a fresh App against the same cache dir/project.
	app2, err := NewApp(cfg, log)
	if err != nil {
		t.Fatalf("NewApp (restart): %v", err)
	}
	defer app2.Close()

	raw2, src2, err := cachedCall(app2, "getBalance", nil)
	if err != nil {
		t.Fatalf("cachedCall (after restart): %v", err)
	}
	if src2 != "FROM_MEMORY" {
		// LoadAll() at startup should have warmed L1 from disk.
		t.Fatalf("expected FROM_MEMORY (warmed from disk at startup), got %s", src2)
	}
	if apiCalls != 1 {
		t.Fatalf("expected no additional API call after restart, got %d total calls", apiCalls)
	}
	if string(raw1) != string(raw2) {
		t.Fatalf("cached value mismatch across restart: %s vs %s", raw1, raw2)
	}
}

// TestProjectIsolation verifies SPEC 03's per-project cache isolation:
// two projects under the same cache dir never see each other's entries.
func TestProjectIsolation(t *testing.T) {
	apiCalls := map[string]int{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		apiCalls[r.Form.Get("nonce")]++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":1,"balance":{"BTC":"1.5"}}`))
	}))
	defer srv.Close()
	restore := api.SetAPIURLForTesting(srv.URL)
	defer restore()

	cacheDir := t.TempDir()
	log := logging.New(&bytes.Buffer{}, logging.ParseLevel("error"), false, nil)

	mk := func(project string) *App {
		cfg := &config.Config{
			APIKey: "k", APISecret: "s", Tier: "pro", Project: project,
			CacheSize: 100, CacheDir: cacheDir, LogLevel: "error", LogFormat: "text", Timezone: "UTC",
		}
		app, err := NewApp(cfg, log)
		if err != nil {
			t.Fatalf("NewApp(%s): %v", project, err)
		}
		t.Cleanup(func() { app.Close() })
		return app
	}

	appA := mk("client_a")
	appB := mk("client_b")

	if _, srcA, err := cachedCall(appA, "getBalance", nil); err != nil || srcA != "FROM_API" {
		t.Fatalf("appA first call: src=%s err=%v", srcA, err)
	}
	// appB must still miss (separate LRU + separate disk file) despite
	// identical method/params.
	if _, srcB, err := cachedCall(appB, "getBalance", nil); err != nil || srcB != "FROM_API" {
		t.Fatalf("appB first call: src=%s err=%v", srcB, err)
	}

	dbA := filepath.Join(cacheDir, "client_a", "cointracking-mcp.db")
	dbB := filepath.Join(cacheDir, "client_b", "cointracking-mcp.db")
	if dbA == dbB {
		t.Fatal("expected distinct db paths per project")
	}
}

// TestSwitchProject verifies the live project-switch tool: a single App
// moves between projects without restarting, each project's cache stays
// isolated, and switching back to a project it already visited serves from
// disk (no extra API call) instead of starting cold.
func TestSwitchProject(t *testing.T) {
	apiCalls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":1,"balance":{"BTC":"1.5"}}`))
	}))
	defer srv.Close()
	restore := api.SetAPIURLForTesting(srv.URL)
	defer restore()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		APIKey: "k", APISecret: "s", Tier: "pro", Project: "agp",
		CacheSize: 100, CacheDir: cacheDir, LogLevel: "error", LogFormat: "text", Timezone: "UTC",
	}
	log := logging.New(&bytes.Buffer{}, logging.ParseLevel(cfg.LogLevel), false, nil)

	app, err := NewApp(cfg, log)
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}
	defer app.Close()

	if app.Project() != "agp" {
		t.Fatalf("expected initial project agp, got %s", app.Project())
	}

	if _, src, err := cachedCall(app, "getBalance", nil); err != nil || src != "FROM_API" {
		t.Fatalf("agp first call: src=%s err=%v", src, err)
	}
	if apiCalls != 1 {
		t.Fatalf("expected 1 API call, got %d", apiCalls)
	}

	// Switching to an invalid name must be rejected without touching state.
	if _, err := app.SwitchProject("bad name!"); err == nil {
		t.Fatal("expected error for invalid project name")
	}
	if app.Project() != "agp" {
		t.Fatalf("project must be unchanged after a rejected switch, got %s", app.Project())
	}

	// Switching to the same project is a no-op that reports AlreadyActive.
	res, err := app.SwitchProject("agp")
	if err != nil {
		t.Fatalf("SwitchProject (same project): %v", err)
	}
	if !res.AlreadyActive {
		t.Fatal("expected AlreadyActive=true when switching to the current project")
	}

	// Switch to a new project: must not see agp's cached entry.
	res, err = app.SwitchProject("cliente_b")
	if err != nil {
		t.Fatalf("SwitchProject (cliente_b): %v", err)
	}
	if res.AlreadyActive {
		t.Fatal("expected AlreadyActive=false when switching projects")
	}
	if app.Project() != "cliente_b" {
		t.Fatalf("expected active project cliente_b, got %s", app.Project())
	}
	if _, src, err := cachedCall(app, "getBalance", nil); err != nil || src != "FROM_API" {
		t.Fatalf("cliente_b first call: src=%s err=%v", src, err)
	}
	if apiCalls != 2 {
		t.Fatalf("expected 2 API calls after switching to a fresh project, got %d", apiCalls)
	}

	// Switch back to agp: its cache was flushed to disk on the way out, so
	// this must warm from disk, not hit the API again.
	if _, err := app.SwitchProject("agp"); err != nil {
		t.Fatalf("SwitchProject (back to agp): %v", err)
	}
	if _, src, err := cachedCall(app, "getBalance", nil); err != nil || src != "FROM_MEMORY" {
		t.Fatalf("expected FROM_MEMORY switching back to agp, got src=%s err=%v", src, err)
	}
	if apiCalls != 2 {
		t.Fatalf("expected no additional API call switching back to a previously-visited project, got %d", apiCalls)
	}
}
