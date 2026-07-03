package config

import "testing"

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
