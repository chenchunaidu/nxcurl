package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/chenchunaidu/nxcurl/internal/paths"
)

// catalog mirrors tars-style tool metadata: one JSON file agents can read.
var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Print path to the agent-oriented tool catalog JSON (and refresh it)",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := paths.EnsureDataDir()
		if err != nil {
			return err
		}
		dir := filepath.Join(base, "catalog")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		p := filepath.Join(dir, "tools.json")
		doc := map[string]any{
			"tools": []map[string]any{
				{
					"name":        "nxcurl",
					"description": "HTTP client CLI: run requests, history, env substitution {{VAR}}, import HAR/Postman; use --json for structured output for agents.",
					"usage": `nxcurl run <url> [-X METHOD] [-H 'Name: value']... [-d body] [-e env] [--json]
nxcurl send <request.json> [-e env] [--json]
nxcurl history list [--limit N] [--json]
nxcurl history show <id> [--json]
nxcurl history replay <id> [-e env] [--json]
nxcurl import har <file.har>
nxcurl import postman <collection.json>
nxcurl env init <name> && nxcurl env path <name>
Environment variables in URL/headers/body: {{API_KEY}} from ~/.nxcurl/envs/<name>.json`,
				},
			},
		}
		b, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(p, b, 0o644); err != nil {
			return err
		}
		fmt.Println(p)
		return nil
	},
}
