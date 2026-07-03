// Package config parses and validates cointracking-mcp startup configuration:
// CLI flags for everything, with environment-variable fallback only for credentials.
package config

import (
	"flag"
	"fmt"
	"os"
	"regexp"
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
	logLevel := fs.String("log-level", "info", "Log level: debug|info|warn|error")
	logFormat := fs.String("log-format", "text", "Log format: text|json")
	timezone := fs.String("timezone", "UTC", "IANA timezone for log timestamps")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		Tier:      *tier,
		Project:   *project,
		CacheDir:  *cacheDir,
		LogLevel:  *logLevel,
		LogFormat: *logFormat,
		Timezone:  *timezone,
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
