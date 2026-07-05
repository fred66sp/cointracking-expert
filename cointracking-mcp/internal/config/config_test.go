package config

import (
	"os"
	"testing"
)

func TestObfuscate(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"abc123def456xyz789uvw", "abc123***uvw"},
		{"", "***"},
		{"ab", "***"},
	}
	for _, c := range cases {
		if got := Obfuscate(c.in); got != c.want {
			t.Errorf("Obfuscate(%q) = %q, want %q", c.in, got, c.want)
		}
	}

	// Never reveal the full secret for any non-trivial input.
	secret := "secret_key_xyz_12345"
	if got := Obfuscate(secret); got == secret {
		t.Errorf("Obfuscate(%q) leaked the full secret", secret)
	}
}

func TestParseRequiresCredentials(t *testing.T) {
	t.Setenv("COINTRACKING_API_KEY", "")
	t.Setenv("COINTRACKING_API_SECRET", "")
	if _, err := Parse([]string{}); err == nil {
		t.Fatal("expected error when no credentials provided")
	}
}

func TestParseCLIFlagsOverrideEnv(t *testing.T) {
	t.Setenv("COINTRACKING_API_KEY", "env_key")
	t.Setenv("COINTRACKING_API_SECRET", "env_secret")
	cfg, err := Parse([]string{"--api-key", "cli_key", "--api-secret", "cli_secret"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cfg.APIKey != "cli_key" || cfg.APISecret != "cli_secret" {
		t.Fatalf("expected CLI flags to take precedence, got key=%q secret=%q", cfg.APIKey, cfg.APISecret)
	}
}

func TestParseEnvFallback(t *testing.T) {
	t.Setenv("COINTRACKING_API_KEY", "env_key")
	t.Setenv("COINTRACKING_API_SECRET", "env_secret")
	cfg, err := Parse([]string{})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cfg.APIKey != "env_key" || cfg.APISecret != "env_secret" {
		t.Fatalf("expected env fallback, got key=%q secret=%q", cfg.APIKey, cfg.APISecret)
	}
}

func TestParseDefaults(t *testing.T) {
	cfg, err := Parse([]string{"--api-key", "k", "--api-secret", "s"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cfg.Project != "default" || cfg.Tier != "pro" || cfg.CacheSize != 5000 || cfg.LogLevel != "info" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if cfg.RateLimit() != 20 {
		t.Fatalf("expected pro tier rate limit 20, got %d", cfg.RateLimit())
	}
}

func TestParseTierSizing(t *testing.T) {
	cfg, err := Parse([]string{"--api-key", "k", "--api-secret", "s", "--tier", "unlimited"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cfg.CacheSize != 150000 || cfg.RateLimit() != 60 {
		t.Fatalf("expected unlimited tier sizing, got cacheSize=%d rateLimit=%d", cfg.CacheSize, cfg.RateLimit())
	}
}

func TestParseInvalidValues(t *testing.T) {
	base := []string{"--api-key", "k", "--api-secret", "s"}
	cases := map[string][]string{
		"tier":       append(base, "--tier", "gold"),
		"project":    append(base, "--project", "has spaces"),
		"log-level":  append(base, "--log-level", "verbose"),
		"log-format": append(base, "--log-format", "yaml"),
	}
	for name, args := range cases {
		if _, err := Parse(args); err == nil {
			t.Errorf("%s: expected validation error, got none", name)
		}
	}
}

// --- ADR-040: credenciales por proyecto ---

func newEnvDirConfig(t *testing.T) (*Config, string) {
	t.Helper()
	dir := t.TempDir()
	cfg, err := Parse([]string{"--api-key", "proc-key", "--api-secret", "proc-secret",
		"--tier", "pro", "--project-env-dir", dir})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return cfg, dir
}

func TestResolveCredentialsProcessFallbacks(t *testing.T) {
	// Sin --project-env-dir: siempre proceso.
	cfg, err := Parse([]string{"--api-key", "proc-key", "--api-secret", "proc-secret"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	c, err := cfg.ResolveProjectCredentials("cualquiera")
	if err != nil || c.Source != "process" || c.APIKey != "proc-key" {
		t.Fatalf("expected process fallback, got %+v err=%v", c, err)
	}

	// Con dir pero sin fichero para el proyecto: también proceso.
	cfg2, _ := newEnvDirConfig(t)
	c2, err := cfg2.ResolveProjectCredentials("sin_fichero")
	if err != nil || c2.Source != "process" || c2.APIKey != "proc-key" {
		t.Fatalf("expected process fallback for missing file, got %+v err=%v", c2, err)
	}
}

func TestResolveCredentialsFromProjectEnv(t *testing.T) {
	cfg, dir := newEnvDirConfig(t)
	envContent := "# cuenta del cliente B\nCOINTRACKING_API_KEY = \"b-key\"\nCOINTRACKING_API_SECRET='b-secret'\nCOINTRACKING_TIER=unlimited\n"
	if err := os.WriteFile(dir+"/cliente_b.env", []byte(envContent), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	c, err := cfg.ResolveProjectCredentials("cliente_b")
	if err != nil {
		t.Fatalf("ResolveProjectCredentials: %v", err)
	}
	if c.Source != "project-env" || c.APIKey != "b-key" || c.APISecret != "b-secret" || c.Tier != "unlimited" {
		t.Fatalf("bad resolution: %+v", c)
	}
	if lim, _ := RateLimitForTier(c.Tier); lim != 60 {
		t.Fatalf("expected unlimited=60 calls/h, got %d", lim)
	}
}

func TestResolveCredentialsFailClosed(t *testing.T) {
	cfg, dir := newEnvDirConfig(t)

	// Fichero incompleto (falta el secret) -> error, nunca fallback silencioso.
	if err := os.WriteFile(dir+"/incompleto.env", []byte("COINTRACKING_API_KEY=x\n"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if _, err := cfg.ResolveProjectCredentials("incompleto"); err == nil {
		t.Fatal("expected error for incomplete env file (fail-closed), got none")
	}

	// Tier inválido en el fichero -> error.
	bad := "COINTRACKING_API_KEY=x\nCOINTRACKING_API_SECRET=y\nCOINTRACKING_TIER=gold\n"
	if err := os.WriteFile(dir+"/tier_malo.env", []byte(bad), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if _, err := cfg.ResolveProjectCredentials("tier_malo"); err == nil {
		t.Fatal("expected error for invalid tier in env file, got none")
	}

	// Línea sin '=' -> error de parseo.
	if err := os.WriteFile(dir+"/roto.env", []byte("esto no es un env\n"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if _, err := cfg.ResolveProjectCredentials("roto"); err == nil {
		t.Fatal("expected parse error for malformed env file, got none")
	}
}
