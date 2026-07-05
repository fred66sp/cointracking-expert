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
// cache manager, config, and logger. cfg/cache/store — and, since ADR-040,
// client/rate too (per-project credentials) — are swapped live by
// SwitchProject (one project active at a time, per SPEC 01/03), so handlers
// must go through the accessor methods below instead of reading fields
// directly — that's what keeps a switch mid-request from racing a read.
type App struct {
	Log *logging.Logger

	mu     sync.RWMutex
	cfg    *config.Config
	cache  *cache.Manager
	store  *cache.Store
	client *api.Client
	rate   *api.RateTracker
	creds  *config.ProjectCredentials
}

// NewApp wires the client, rate tracker, and L1/L2 cache per the configured
// project, per SPEC/03-cache-strategy.md (cache isolated under
// {CACHE_PERSIST_DIR}/{PROJECT_NAME}/). Credentials resolve per project when
// --project-env-dir is set (ADR-040); otherwise the process credentials apply.
func NewApp(cfg *config.Config, log *logging.Logger) (*App, error) {
	creds, err := cfg.ResolveProjectCredentials(cfg.Project)
	if err != nil {
		return nil, err
	}
	limit, err := config.RateLimitForTier(creds.Tier)
	if err != nil {
		return nil, err
	}
	if creds.Source == "project-env" {
		log.Infof("Credenciales por proyecto (ADR-040): %q usa %s (API key %s, tier %s)",
			cfg.Project, cfg.ProjectEnvDir, config.Obfuscate(creds.APIKey), creds.Tier)
	}
	rate := api.NewRateTracker(limit)
	client := api.NewClient(creds.APIKey, creds.APISecret, rate)

	mgr, store, _, err := openProjectCache(cfg, log)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:    cfg,
		client: client,
		rate:   rate,
		creds:  creds,
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

// Client returns the API client for the currently active project's account
// (per-project credentials, ADR-040; process credentials by default).
func (a *App) Client() *api.Client {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.client
}

// Rate returns the rate tracker for the currently active account.
func (a *App) Rate() *api.RateTracker {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.rate
}

// ErrProjectActive is returned by WithProjectLockedIfNotActive when name is
// the currently active project.
var ErrProjectActive = fmt.Errorf("project is currently active")

// WithProjectLockedIfNotActive runs fn while holding the same lock
// SwitchProject uses, after checking name isn't the active project — closing
// the TOCTOU window a bare "check app.Project(), then act" would have (found
// 2026-07-05, independent robustness review): without this, a concurrent
// SwitchProject(name) landing between the check and, say, a delete could
// activate name right before its on-disk cache gets removed out from under
// the just-opened SQLite handle. Held for the duration of fn, so fn should
// be reasonably quick (e.g. removeAllWithRetry's backoff tops out in the
// hundreds of ms) — it blocks other project-switching tools meanwhile.
func (a *App) WithProjectLockedIfNotActive(name string, fn func() error) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if name == a.cfg.Project {
		return ErrProjectActive
	}
	return fn()
}

// SwitchProjectResult reports what SwitchProject did, for the MCP tool's
// response. CredentialsSource/APIKeyObfuscated make the active account
// auditable in every switch without ever exposing a secret (ADR-040/011).
type SwitchProjectResult struct {
	Project               string `json:"project"`
	PreviousProject       string `json:"previous_project,omitempty"`
	AlreadyActive         bool   `json:"already_active"`
	EntriesClearedPrev    int    `json:"entries_cleared_previous_project"`
	EntriesLoadedFromDisk int    `json:"entries_loaded_from_disk"`
	CredentialsSource     string `json:"credentials_source"`
	APIKeyObfuscated      string `json:"api_key_obfuscated"`
	Message               string `json:"message"`
}

// SwitchProject moves the server from whatever project is currently active
// to name, live: flushes and closes the current project's cache, then opens
// (or creates) name's cache the same way NewApp does at startup.
//
// Which ACCOUNT the tools query after the switch depends on --project-env-dir
// (ADR-040): if {project-env-dir}/{name}.env exists, the target project gets
// its own API client (and fresh rate tracker — CoinTracking's hourly limit is
// per key); otherwise the process credentials keep applying and the current
// client/rate are reused (preserving the consumed hourly window). Without
// --project-env-dir every project talks to the SAME CoinTracking account and
// only cached data is isolated — in that mode, do NOT use project switching
// to audit two different accounts.
//
// The target's credentials are resolved and validated BEFORE anything is
// closed (a malformed .env aborts cleanly without touching state), and the
// client swap happens only after the target cache opened successfully.
func (a *App) SwitchProject(name string) (*SwitchProjectResult, error) {
	if err := config.ValidateProjectName(name); err != nil {
		return nil, err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if name == a.cfg.Project {
		return &SwitchProjectResult{
			Project:           name,
			AlreadyActive:     true,
			CredentialsSource: a.creds.Source,
			APIKeyObfuscated:  config.Obfuscate(a.creds.APIKey),
			Message:           fmt.Sprintf("El proyecto %q ya estaba activo; no se ha tocado la caché.", name),
		}, nil
	}

	// Fail-early (ADR-040): resolver credenciales del destino antes de cerrar
	// nada. Un .env malformado o incompleto aborta aquí, con el proyecto
	// actual intacto.
	newCreds, err := a.cfg.ResolveProjectCredentials(name)
	if err != nil {
		return nil, err
	}
	newLimit, err := config.RateLimitForTier(newCreds.Tier)
	if err != nil {
		return nil, err
	}

	prevCfg := a.cfg
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
		// The previous project's store is already closed at this point.
		// Leaving a.cache/a.store pointing at it would mean every
		// subsequent tool call fails against a closed SQLite handle
		// instead of a clear error. Try to reopen the previous project
		// so the server stays in a known-good state instead of a broken
		// one — the caller still finds out the switch itself failed.
		rollbackMgr, rollbackStore, _, rollbackErr := openProjectCache(prevCfg, a.Log)
		if rollbackErr == nil {
			a.cache = rollbackMgr
			a.store = rollbackStore
			a.Log.Warnf("switch_project: failed to open %q (%s); rolled back to %q", name, err, prevProject)
			return nil, fmt.Errorf("switching to project %q: %w (permaneces en %q)", name, err, prevProject)
		}
		a.Log.Errorf("switch_project: failed to open %q (%s) AND failed to roll back to %q (%s) — server has no usable cache until restarted",
			name, err, prevProject, rollbackErr)
		return nil, fmt.Errorf("switching to project %q: %w (y no se pudo volver a %q: %v — reinicia el servidor MCP)",
			name, err, prevProject, rollbackErr)
	}

	a.cfg = &newCfg
	a.cache = mgr
	a.store = store

	// Swap de cuenta solo si las credenciales resueltas difieren; si son las
	// mismas se conservan cliente y tracker (mantiene la ventana horaria ya
	// consumida — el límite de CoinTracking es por API key).
	if newCreds.APIKey != a.creds.APIKey || newCreds.APISecret != a.creds.APISecret {
		a.rate = api.NewRateTracker(newLimit)
		a.client = api.NewClient(newCreds.APIKey, newCreds.APISecret, a.rate)
		a.Log.Infof("Cuenta cambiada con el proyecto (ADR-040): %q usa API key %s (fuente: %s, tier %s)",
			name, config.Obfuscate(newCreds.APIKey), newCreds.Source, newCreds.Tier)
	}
	a.creds = newCreds

	a.Log.Infof("Proyecto cambiado: %q -> %q (%d entradas liberadas de memoria, %d cargadas de disco)",
		prevProject, name, prevCleared, loaded)

	return &SwitchProjectResult{
		Project:               name,
		PreviousProject:       prevProject,
		EntriesClearedPrev:    prevCleared,
		EntriesLoadedFromDisk: loaded,
		CredentialsSource:     newCreds.Source,
		APIKeyObfuscated:      config.Obfuscate(newCreds.APIKey),
		Message: fmt.Sprintf("Proyecto activo cambiado de %q a %q. Caché aislada en %s. Cuenta: API key %s (%s).",
			prevProject, name, filepath.Join(newCfg.CacheDir, name),
			config.Obfuscate(newCreds.APIKey), newCreds.Source),
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
