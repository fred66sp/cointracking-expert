// Package config parses and validates cointracking-mcp startup configuration:
// CLI flags for everything, with environment-variable fallback only for credentials.
package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var projectNameRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// TierLimits maps an account tier to its cache size and hourly API rate limit.
type TierLimits struct {
	CacheSize int
	RateLimit int
}

var tiers = map[string]TierLimits{
	"pro":       {CacheSize: 5000, RateLimit: 20},
	"expert":    {CacheSize: 50000, RateLimit: 20},
	"unlimited": {CacheSize: 150000, RateLimit: 60},
}

// TTLs are hardcoded per spec 06 (no CLI flag), in seconds.
const (
	TTLTrades         = 3600
	TTLBalance        = 600
	TTLGroupedBalance = 600
	TTLSummary        = 7200
	TTLCurrency       = 7200
	TTLGains          = 3600
)

type Config struct {
	APIKey    string
	APISecret string

	Tier      string
	Project   string
	CacheSize int
	CacheDir  string

	// ProjectEnvDir, si no está vacío, habilita credenciales por proyecto
	// (ADR-040): al activar un proyecto se busca {ProjectEnvDir}/{project}.env
	// con COINTRACKING_API_KEY / COINTRACKING_API_SECRET (y COINTRACKING_TIER
	// opcional). Vacío = todas las llamadas usan las credenciales del proceso.
	ProjectEnvDir string

	LogLevel  string
	LogFormat string

	Timezone string
}

// Parse builds a Config from CLI args and environment variables, applying the
// precedence documented in SPEC/06-configuration.md: CLI flag > envvar > error
// (credentials only); CLI flag > default (everything else).
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("cointracking-mcp", flag.ContinueOnError)

	apiKey := fs.String("api-key", "", "CoinTracking API key (or COINTRACKING_API_KEY)")
	apiSecret := fs.String("api-secret", "", "CoinTracking API secret (or COINTRACKING_API_SECRET)")
	tier := fs.String("tier", "pro", "Account tier: pro|expert|unlimited")
	project := fs.String("project", "default", "Project name (isolates cache)")
	cacheMaxSize := fs.Int("cache-max-size", 0, "Max cache entries (default: per tier)")
	cacheDir := fs.String("cache-dir", "./cache", "Cache persistence directory")
	projectEnvDir := fs.String("project-env-dir", "", "Directory with per-project credential files <project>.env (ADR-040); empty = single-account mode")
	logLevel := fs.String("log-level", "info", "Log level: debug|info|warn|error")
	logFormat := fs.String("log-format", "text", "Log format: text|json")
	timezone := fs.String("timezone", "UTC", "IANA timezone for log timestamps")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		Tier:          *tier,
		Project:       *project,
		CacheDir:      *cacheDir,
		ProjectEnvDir: *projectEnvDir,
		LogLevel:      *logLevel,
		LogFormat:     *logFormat,
		Timezone:      *timezone,
	}

	// Credentials: CLI flag > envvar > fatal error.
	cfg.APIKey = *apiKey
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("COINTRACKING_API_KEY")
	}
	cfg.APISecret = *apiSecret
	if cfg.APISecret == "" {
		cfg.APISecret = os.Getenv("COINTRACKING_API_SECRET")
	}
	if cfg.APIKey == "" || cfg.APISecret == "" {
		return nil, fmt.Errorf("--api-key y --api-secret requeridos (o variables de entorno COINTRACKING_API_KEY / COINTRACKING_API_SECRET)")
	}

	limits, ok := tiers[cfg.Tier]
	if !ok {
		return nil, fmt.Errorf("--tier inválido: %q (válidos: pro, expert, unlimited)", cfg.Tier)
	}

	if err := ValidateProjectName(cfg.Project); err != nil {
		return nil, err
	}

	switch cfg.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return nil, fmt.Errorf("--log-level inválido: %q (válidos: debug, info, warn, error)", cfg.LogLevel)
	}

	switch cfg.LogFormat {
	case "text", "json":
	default:
		return nil, fmt.Errorf("--log-format inválido: %q (válidos: text, json)", cfg.LogFormat)
	}

	cfg.CacheSize = *cacheMaxSize
	if cfg.CacheSize <= 0 {
		cfg.CacheSize = limits.CacheSize
	}

	return cfg, nil
}

// ValidateProjectName checks a project name against the same rule enforced
// at startup for --project, so cointracking_switch_project can reject a bad
// name before touching any cache directory.
func ValidateProjectName(name string) error {
	if !projectNameRe.MatchString(name) {
		return fmt.Errorf("nombre de proyecto inválido: %q (solo alfanumérico, _ y -)", name)
	}
	return nil
}

// RateLimit returns the hourly API call limit for the configured tier.
func (c *Config) RateLimit() int {
	return tiers[c.Tier].RateLimit
}

// RateLimitForTier returns the hourly API limit for an arbitrary tier name
// (per-project tier override, ADR-040).
func RateLimitForTier(tier string) (int, error) {
	limits, ok := tiers[tier]
	if !ok {
		return 0, fmt.Errorf("tier inválido: %q (válidos: pro, expert, unlimited)", tier)
	}
	return limits.RateLimit, nil
}

// ProjectCredentials is the result of resolving which CoinTracking account a
// project uses (ADR-040).
type ProjectCredentials struct {
	APIKey    string
	APISecret string
	Tier      string
	// Source is "process" (fallback/default) or "project-env" (loaded from
	// {ProjectEnvDir}/{project}.env).
	Source string
}

// ResolveProjectCredentials decides which credentials a project uses, per
// ADR-040. Fail-closed by design: a credentials file that EXISTS but is
// incomplete or unreadable is an error, never a silent fallback to the
// process account — querying the wrong CoinTracking account is exactly the
// accident this feature exists to prevent.
func (c *Config) ResolveProjectCredentials(project string) (*ProjectCredentials, error) {
	fromProcess := &ProjectCredentials{
		APIKey: c.APIKey, APISecret: c.APISecret, Tier: c.Tier, Source: "process",
	}
	if c.ProjectEnvDir == "" {
		return fromProcess, nil
	}
	envPath := filepath.Join(c.ProjectEnvDir, project+".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return fromProcess, nil
	}
	vars, err := parseEnvFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("leyendo credenciales de proyecto %s: %w", envPath, err)
	}
	key, secret := vars["COINTRACKING_API_KEY"], vars["COINTRACKING_API_SECRET"]
	if key == "" || secret == "" {
		return nil, fmt.Errorf(
			"%s existe pero está incompleto (COINTRACKING_API_KEY y COINTRACKING_API_SECRET son obligatorios) — "+
				"se aborta en vez de degradar a la cuenta del proceso (ADR-040 fail-closed)", envPath)
	}
	tier := vars["COINTRACKING_TIER"]
	if tier == "" {
		tier = c.Tier
	}
	if _, err := RateLimitForTier(tier); err != nil {
		return nil, fmt.Errorf("en %s: %w", envPath, err)
	}
	return &ProjectCredentials{APIKey: key, APISecret: secret, Tier: tier, Source: "project-env"}, nil
}

// parseEnvFile reads a minimal KEY=VALUE file: blank lines and #-comments
// ignored, optional surrounding single/double quotes stripped. No expansion,
// no "export" prefix — deliberately strict and predictable.
func parseEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	vars := make(map[string]string)
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("línea %d: se esperaba KEY=VALUE, encontrado %q", i+1, line)
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if len(v) >= 2 && (v[0] == '"' && v[len(v)-1] == '"' || v[0] == '\'' && v[len(v)-1] == '\'') {
			v = v[1 : len(v)-1]
		}
		vars[k] = v
	}
	return vars, nil
}

// Obfuscate renders a credential as "first6***last3" for safe logging
// (per SPEC/06-configuration.md). Shorter secrets degrade gracefully.
func Obfuscate(secret string) string {
	const headLen, tailLen = 6, 3
	if len(secret) <= headLen+tailLen {
		if len(secret) <= 3 {
			return "***"
		}
		return secret[:2] + "***"
	}
	return secret[:headLen] + "***" + secret[len(secret)-tailLen:]
}
