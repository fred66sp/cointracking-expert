// Package tools implements the MCP tool handlers exposed by cointracking-mcp:
// the 6 CoinTracking API tools plus cache/project control (SPEC 02, 03, 04).
package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alfredogp/cointracking-mcp/internal/api"
	"github.com/alfredogp/cointracking-mcp/internal/cache"
	"github.com/alfredogp/cointracking-mcp/internal/config"
	"github.com/alfredogp/cointracking-mcp/internal/logging"
)

// App holds the shared state every tool handler needs: the API client, the
// cache manager, config, and logger. Credentials/tier/rate limiter are fixed
// for the process lifetime, but cfg/cache/store are swapped live by
// SwitchProject (one project active at a time, per SPEC 01/03), so handlers
// must go through the accessor methods below instead of reading fields
// directly — that's what keeps a switch mid-request from racing a read.
type App struct {
	Client *api.Client
	Rate   *api.RateTracker
	Log    *logging.Logger

	mu    sync.RWMutex
	cfg   *config.Config
	cache *cache.Manager
	store *cache.Store
}

// NewApp wires the client, rate tracker, and L1/L2 cache per the configured
// project, per SPEC/03-cache-strategy.md (cache isolated under
// {CACHE_PERSIST_DIR}/{PROJECT_NAME}/).
func NewApp(cfg *config.Config, log *logging.Logger) (*App, error) {
	rate := api.NewRateTracker(cfg.RateLimit())
	client := api.NewClient(cfg.APIKey, cfg.APISecret, rate)

	mgr, store, _, err := openProjectCache(cfg, log)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:    cfg,
		Client: client,
		Rate:   rate,
		cache:  mgr,
		Log:    log,
		store:  store,
	}, nil
}

// openProjectCache creates (if needed) cfg.CacheDir/cfg.Project, opens its
// SQLite store, purges expired rows, and warms an L1 LRU from disk. Shared by
// NewApp (startup) and SwitchProject (live project change) so both open a
// project's cache exactly the same way.
func openProjectCache(cfg *config.Config, log *logging.Logger) (*cache.Manager, *cache.Store, int, error) {
	projectDir := filepath.Join(cfg.CacheDir, cfg.Project)
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		return nil, nil, 0, fmt.Errorf("creating cache dir %s: %w", projectDir, err)
	}
	dbPath := filepath.Join(projectDir, "cointracking-mcp.db")
	store, err := cache.OpenStore(dbPath)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("opening cache store: %w", err)
	}

	purged, err := store.PurgeExpired()
	if err != nil {
		log.Warnf("cache: error purging expired entries at startup: %s", err)
	} else if purged > 0 {
		log.Infof("Cache: purged %d expired entries from disk", purged)
	}

	l1 := cache.NewLRU(cfg.CacheSize)
	loaded, err := store.LoadAll()
	if err != nil {
		log.Warnf("cache: error loading entries from disk: %s", err)
	} else {
		for key, entry := range loaded {
			l1.Set(key, entry)
		}
		log.Infof("Loaded %d cached entries from disk", len(loaded))
	}

	return cache.NewManager(l1, store), store, len(loaded), nil
}

// Project returns the currently active project name.
func (a *App) Project() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cfg.Project
}

// CacheDir returns the configured base cache directory (fixed for the
// process; only the project subdirectory under it changes on a switch).
func (a *App) CacheDir() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cfg.CacheDir
}

// CacheManager returns the cache manager for the currently active project.
func (a *App) CacheManager() *cache.Manager {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cache
}

// Store returns the disk store for the currently active project.
func (a *App) Store() *cache.Store {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.store
}

// SwitchProjectResult reports what SwitchProject did, for the MCP tool's
// response.
type SwitchProjectResult struct {
	Project               string `json:"project"`
	PreviousProject       string `json:"previous_project,omitempty"`
	AlreadyActive         bool   `json:"already_active"`
	EntriesClearedPrev    int    `json:"entries_cleared_previous_project"`
	EntriesLoadedFromDisk int    `json:"entries_loaded_from_disk"`
	Message               string `json:"message"`
}

// SwitchProject moves the server from whatever project is currently active
// to name, live: flushes and closes the current project's cache, then opens
// (or creates) name's cache the same way NewApp does at startup. Credentials,
// tier, and rate limiter are per-process and untouched — only the cache
// (and therefore which account's data cointracking_get_* tools return)
// changes. This is what lets the ADR-013 "active project" gate switch
// projects mid-conversation without restarting the MCP server or touching
// .mcp.json.
func (a *App) SwitchProject(name string) (*SwitchProjectResult, error) {
	if err := config.ValidateProjectName(name); err != nil {
		return nil, err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if name == a.cfg.Project {
		return &SwitchProjectResult{
			Project:       name,
			AlreadyActive: true,
			Message:       fmt.Sprintf("El proyecto %q ya estaba activo; no se ha tocado la caché.", name),
		}, nil
	}

	prevProject := a.cfg.Project
	prevCleared := a.cache.L1.Clear()
	a.store.Flush()
	if err := a.store.Close(); err != nil {
		a.Log.Warnf("switch_project: error closing store for project %q: %s", prevProject, err)
	}

	newCfg := *a.cfg
	newCfg.Project = name

	mgr, store, loaded, err := openProjectCache(&newCfg, a.Log)
	if err != nil {
		return nil, fmt.Errorf("switching to project %q: %w", name, err)
	}

	a.cfg = &newCfg
	a.cache = mgr
	a.store = store

	a.Log.Infof("Proyecto cambiado: %q -> %q (%d entradas liberadas de memoria, %d cargadas de disco)",
		prevProject, name, prevCleared, loaded)

	return &SwitchProjectResult{
		Project:               name,
		PreviousProject:       prevProject,
		EntriesClearedPrev:    prevCleared,
		EntriesLoadedFromDisk: loaded,
		Message: fmt.Sprintf("Proyecto activo cambiado de %q a %q. Caché aislada en %s.",
			prevProject, name, filepath.Join(newCfg.CacheDir, name)),
	}, nil
}

// Close flushes and closes the disk store (graceful shutdown / close_project).
func (a *App) Close() error {
	a.mu.RLock()
	store := a.store
	a.mu.RUnlock()
	if store == nil {
		return nil
	}
	return store.Close()
}

// TTLFor returns the cache TTL for a given CoinTracking API method, per
// SPEC/03-cache-strategy.md.
func TTLFor(method string) time.Duration {
	switch method {
	case "getTrades":
		return config.TTLTrades * time.Second
	case "getBalance":
		return config.TTLBalance * time.Second
	case "getGroupedBalance":
		return config.TTLGroupedBalance * time.Second
	case "getHistoricalSummary":
		return config.TTLSummary * time.Second
	case "getHistoricalCurrency":
		return config.TTLCurrency * time.Second
	case "getGains":
		return config.TTLGains * time.Second
	default:
		return 10 * time.Minute
	}
}
