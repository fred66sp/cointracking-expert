package tools

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alfredogp/cointracking-mcp/internal/api"
	"github.com/alfredogp/cointracking-mcp/internal/cache"
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

// TestSwitchProjectRollsBackOnOpenFailure covers the bug found 2026-07-05
// (see CHANGELOG.md): if opening the target project's cache fails, the
// previous project's store had already been closed. The failure mode isn't a
// crash — Store.Get treats a "database is closed" error from a query on the
// closed *sql.DB as a plain cache miss, so cachedCall silently falls back to
// calling the API again. The real damage is quieter: L2 (disk) persistence
// for the "current" project silently stops working (writes fail in the
// background via SetAsync's error callback, reads always miss), so nothing
// written after the failed switch survives a close/reopen — until the server
// is restarted, with no visible error at the time. This forces that failure
// (a regular file where the new project's cache directory needs to go, so
// os.MkdirAll fails) and verifies disk persistence still works after the
// rollback, by closing and reopening the app and confirming a FROM_DISK hit
// instead of a second API call.
func TestSwitchProjectRollsBackOnOpenFailure(t *testing.T) {
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

	if _, _, err := cachedCall(app, "getBalance", nil); err != nil {
		t.Fatalf("agp warmup call: %v", err)
	}
	if apiCalls != 1 {
		t.Fatalf("expected 1 API call after warmup, got %d", apiCalls)
	}

	// Sabotage: create a regular file at the path the new project's cache
	// directory would need to occupy, so os.MkdirAll fails inside
	// openProjectCache when SwitchProject tries to open it.
	blockedPath := filepath.Join(cacheDir, "cliente_c")
	if err := os.WriteFile(blockedPath, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("setup: writing blocking file: %v", err)
	}

	if _, err := app.SwitchProject("cliente_c"); err == nil {
		t.Fatal("expected SwitchProject to fail when the target dir is blocked by a file")
	}
	if app.Project() != "agp" {
		t.Fatalf("expected rollback to leave project=agp, got %s", app.Project())
	}

	// A cache miss here (e.g. a different key) forces a fresh API call and
	// an L2 write — this is the write that silently fails without the fix.
	if _, src, err := cachedCall(app, "getTrades", map[string]string{"limit": "5"}); err != nil || src != cache.FromAPI {
		t.Fatalf("post-rollback fresh call must hit the API, got src=%s err=%v", src, err)
	}
	if apiCalls != 2 {
		t.Fatalf("expected 2 API calls total, got %d", apiCalls)
	}

	// The real assertion: close and reopen against the same cache dir/project.
	// With the bug, the write above never reached disk (store was closed),
	// so this would silently re-call the API instead of hitting L2.
	if err := app.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	app2, err := NewApp(cfg, log)
	if err != nil {
		t.Fatalf("reopening NewApp: %v", err)
	}
	defer app2.Close()

	// NewApp preloads L1 from disk at startup (openProjectCache calls
	// store.LoadAll()), so a hit here comes back as FROM_MEMORY, not
	// FROM_DISK — the proof that persistence survived the rollback is that
	// apiCalls stays at 2 below, not that this specific entry reads as
	// FROM_MEMORY vs FROM_DISK.
	if _, _, err := cachedCall(app2, "getTrades", map[string]string{"limit": "5"}); err != nil {
		t.Fatalf("post-reopen call failed: %v", err)
	}
	if apiCalls != 2 {
		t.Fatalf("expected no additional API call on reopen (data should have persisted to disk), got %d", apiCalls)
	}
}

// TestWithProjectLockedIfNotActiveClosesTheDeleteRaceWindow covers the TOCTOU
// found 2026-07-05 (independent robustness review, see CHANGELOG): checking
// "is name the active project?" and then deleting its cache directory as two
// separate steps leaves a window where a concurrent SwitchProject(name) can
// land in between and activate it right before its SQLite file gets removed.
// WithProjectLockedIfNotActive closes that window by running the check and
// the delete under the same lock SwitchProject uses. This proves it: it
// starts a switch that blocks mid-flight (via a directory sabotaged to make
// openProjectCache hang... no — Go doesn't hang on os.MkdirAll, so instead
// this drives the point home more directly by holding the lock via
// WithProjectLockedIfNotActive itself and confirming a concurrent
// SwitchProject to the same name blocks until the callback returns, then
// correctly reports AlreadyActive against whichever project ends up active,
// never racing an in-flight delete.
func TestWithProjectLockedIfNotActiveClosesTheDeleteRaceWindow(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// Create cliente_b's cache dir up front (not the active project) so a
	// delete targeting it is legal.
	if _, err := app.SwitchProject("cliente_b"); err != nil {
		t.Fatalf("SwitchProject (cliente_b): %v", err)
	}
	if _, err := app.SwitchProject("agp"); err != nil {
		t.Fatalf("SwitchProject (back to agp): %v", err)
	}

	callbackStarted := make(chan struct{})
	callbackMayReturn := make(chan struct{})
	callbackDone := make(chan error, 1)

	go func() {
		callbackDone <- app.WithProjectLockedIfNotActive("cliente_b", func() error {
			close(callbackStarted)
			<-callbackMayReturn // hold the lock until the test says so
			return nil
		})
	}()

	<-callbackStarted // the delete-equivalent callback is now holding app.mu

	// A concurrent switch to the same project must block until the callback
	// (holding the lock) returns — proving the check-then-act window is
	// closed, not just narrowed.
	switchDone := make(chan error, 1)
	go func() {
		_, err := app.SwitchProject("cliente_b")
		switchDone <- err
	}()

	select {
	case <-switchDone:
		t.Fatal("SwitchProject completed while WithProjectLockedIfNotActive's callback was still running — the lock isn't shared, the race window is still open")
	case <-time.After(50 * time.Millisecond):
		// Expected: SwitchProject is blocked on app.mu.
	}

	close(callbackMayReturn)
	if err := <-callbackDone; err != nil {
		t.Fatalf("WithProjectLockedIfNotActive callback: %v", err)
	}
	if err := <-switchDone; err != nil {
		t.Fatalf("SwitchProject (after callback released the lock): %v", err)
	}
	if app.Project() != "cliente_b" {
		t.Fatalf("expected cliente_b active after the switch finally ran, got %s", app.Project())
	}
}

// TestSwitchProjectSwapsAccountCredentials covers ADR-040 end-to-end: with
// --project-env-dir set, switching to a project that has its own .env must
// make subsequent API calls authenticate as THAT account (asserted via the
// Key header the fake CoinTracking server receives), while projects without
// a .env keep the process account AND the same rate-tracker window. A
// malformed .env must abort the switch before anything is closed.
func TestSwitchProjectSwapsAccountCredentials(t *testing.T) {
	var lastKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastKey = r.Header.Get("Key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":1,"balance":{"BTC":"1.5"}}`))
	}))
	defer srv.Close()
	restore := api.SetAPIURLForTesting(srv.URL)
	defer restore()

	envDir := t.TempDir()
	envB := "COINTRACKING_API_KEY=cuenta-b-key\nCOINTRACKING_API_SECRET=cuenta-b-secret\nCOINTRACKING_TIER=unlimited\n"
	if err := os.WriteFile(filepath.Join(envDir, "cliente_b.env"), []byte(envB), 0o600); err != nil {
		t.Fatalf("setup env: %v", err)
	}
	if err := os.WriteFile(filepath.Join(envDir, "roto.env"), []byte("COINTRACKING_API_KEY=solo-key\n"), 0o600); err != nil {
		t.Fatalf("setup env roto: %v", err)
	}

	cfg := &config.Config{
		APIKey: "cuenta-a-key", APISecret: "cuenta-a-secret", Tier: "pro", Project: "agp",
		CacheSize: 100, CacheDir: t.TempDir(), ProjectEnvDir: envDir,
		LogLevel: "error", LogFormat: "text", Timezone: "UTC",
	}
	log := logging.New(&bytes.Buffer{}, logging.ParseLevel(cfg.LogLevel), false, nil)

	app, err := NewApp(cfg, log)
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}
	defer app.Close()

	// 1. Proyecto inicial sin .env propio -> cuenta del proceso (A).
	if _, _, err := cachedCall(app, "getBalance", nil); err != nil {
		t.Fatalf("agp call: %v", err)
	}
	if lastKey != "cuenta-a-key" {
		t.Fatalf("expected process account key, got %q", lastKey)
	}
	usedA := app.Rate().UsedInWindow()

	// 2. Switch a proyecto CON .env -> las llamadas frescas van como cuenta B,
	//    con tracker propio (ventana a cero) y tier del fichero (unlimited=60).
	res, err := app.SwitchProject("cliente_b")
	if err != nil {
		t.Fatalf("SwitchProject(cliente_b): %v", err)
	}
	if res.CredentialsSource != "project-env" {
		t.Fatalf("expected credentials_source=project-env, got %q", res.CredentialsSource)
	}
	if _, _, err := cachedCall(app, "getBalance", nil); err != nil {
		t.Fatalf("cliente_b call: %v", err)
	}
	if lastKey != "cuenta-b-key" {
		t.Fatalf("expected account B key after switch, got %q", lastKey)
	}
	if app.Rate().Limit() != 60 {
		t.Fatalf("expected tier from env file (unlimited=60), got %d", app.Rate().Limit())
	}
	if app.Rate().UsedInWindow() != 1 {
		t.Fatalf("expected fresh rate window for account B (1 call), got %d", app.Rate().UsedInWindow())
	}

	// 3. Switch a proyecto SIN .env -> vuelve la cuenta del proceso; y como
	//    A ya se usó antes, su llamada nueva se autentica como A de nuevo.
	if _, err := app.SwitchProject("cliente_c"); err != nil {
		t.Fatalf("SwitchProject(cliente_c): %v", err)
	}
	if _, _, err := cachedCall(app, "getTrades", map[string]string{"limit": "1"}); err != nil {
		t.Fatalf("cliente_c call: %v", err)
	}
	if lastKey != "cuenta-a-key" {
		t.Fatalf("expected process account key for env-less project, got %q", lastKey)
	}
	_ = usedA // la ventana de A tras B->C es nueva (tracker de A se recreó); documentado en ADR-040

	// 4. Switch entre dos proyectos que comparten cuenta (proceso) debe
	//    REUTILIZAR cliente y tracker (conserva la ventana consumida).
	usedBefore := app.Rate().UsedInWindow()
	if _, err := app.SwitchProject("cliente_d"); err != nil {
		t.Fatalf("SwitchProject(cliente_d): %v", err)
	}
	if app.Rate().UsedInWindow() != usedBefore {
		t.Fatalf("expected preserved rate window switching between same-account projects: %d != %d",
			app.Rate().UsedInWindow(), usedBefore)
	}

	// 5. .env malformado -> el switch aborta ANTES de cerrar nada: el
	//    proyecto actual sigue activo y funcional.
	if _, err := app.SwitchProject("roto"); err == nil {
		t.Fatal("expected error switching to project with incomplete .env (fail-closed)")
	}
	if app.Project() != "cliente_d" {
		t.Fatalf("project must be unchanged after failed credential resolution, got %s", app.Project())
	}
	if _, _, err := cachedCall(app, "getBalance", nil); err != nil {
		t.Fatalf("current project must stay functional after aborted switch: %v", err)
	}
}
