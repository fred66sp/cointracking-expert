// Command cointracking-mcp is an MCP server exposing read-only CoinTracking
// API tools (trades, balances, historicals, gains) with a per-project,
// two-level (memory+disk) cache. See SPEC/ for the full design.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/alfredogp/cointracking-mcp/internal/config"
	"github.com/alfredogp/cointracking-mcp/internal/logging"
	"github.com/alfredogp/cointracking-mcp/internal/tools"
)

const version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "cointracking-mcp: fatal: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		loc = time.UTC
		fmt.Fprintf(os.Stderr, "cointracking-mcp: warning: --timezone %q inválido, usando UTC\n", cfg.Timezone)
	}
	log := logging.New(os.Stderr, logging.ParseLevel(cfg.LogLevel), cfg.LogFormat == "json", loc)

	log.Infof("cointracking-mcp v%s starting", version)
	log.Infof("API Key: %s (ofuscado)", config.Obfuscate(cfg.APIKey))
	log.Infof("Project: %s", cfg.Project)
	log.Infof("Tier: %s (%d calls/hour)", cfg.Tier, cfg.RateLimit())
	log.Infof("Cache: %d entries max, %s/%s", cfg.CacheSize, cfg.CacheDir, cfg.Project)
	log.Infof("Log level: %s", cfg.LogLevel)

	app, err := tools.NewApp(cfg, log)
	if err != nil {
		return fmt.Errorf("initializing app: %w", err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			log.Warnf("error closing cache store: %s", err)
		}
	}()

	server := mcp.NewServer(&mcp.Implementation{Name: "cointracking-mcp", Version: version}, nil)
	tools.RegisterAll(server, app)

	log.Infof("Listening on stdio (MCP)")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
