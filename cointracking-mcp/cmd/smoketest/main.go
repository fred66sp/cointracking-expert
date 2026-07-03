// Command smoketest is a throwaway MCP client used to sanity-check the
// server end-to-end (tools/list + a no-API-call tool) during development.
// Not part of the shipped binaries.
package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "smoketest", Version: "0.0.1"}, nil)

	cmd := exec.Command("./dist/cointracking-mcp.exe",
		"--api-key", "testkey123456789",
		"--api-secret", "testsecret123456789",
		"--project", "smoketest",
		"--cache-dir", "./cache_smoketest",
	)
	transport := &mcp.CommandTransport{Command: cmd}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer session.Close()

	toolsRes, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("ListTools: %v", err)
	}
	fmt.Printf("Registered %d tools:\n", len(toolsRes.Tools))
	for _, t := range toolsRes.Tools {
		fmt.Printf("  - %s\n", t.Name)
	}

	res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "cointracking_cache_stats", Arguments: map[string]any{}})
	if err != nil {
		log.Fatalf("CallTool cache_stats: %v", err)
	}
	fmt.Println("\ncointracking_cache_stats result:")
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			fmt.Println(tc.Text)
		}
	}
	if res.IsError {
		log.Fatal("cache_stats returned IsError=true")
	}

	// Validation error path: getGroupedBalance without required `group`.
	res2, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "cointracking_get_grouped_balance", Arguments: map[string]any{}})
	if err != nil {
		log.Fatalf("CallTool grouped_balance: %v", err)
	}
	fmt.Println("\ncointracking_get_grouped_balance (missing group) result:")
	for _, c := range res2.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			fmt.Println(tc.Text)
		}
	}
	if !res2.IsError {
		log.Fatal("expected IsError=true for missing required 'group' param")
	}

	fmt.Println("\nSMOKETEST OK")
}
