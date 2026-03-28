package cli

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/chenchunaidu/nxcurl/internal/paths"
)

//go:embed agent_skill.md
var agentSkillMarkdown string

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Print agent-oriented SKILL.md and write ~/.nxcurl/docs/SKILL.md",
	Long:  "Writes the embedded agent skill to ~/.nxcurl/docs/SKILL.md and prints the same Markdown to stdout for agents and editors.",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := paths.EnsureDataDir()
		if err != nil {
			return err
		}
		dir := filepath.Join(base, "docs")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		p := filepath.Join(dir, "SKILL.md")
		if err := os.WriteFile(p, []byte(agentSkillMarkdown), 0o644); err != nil {
			return err
		}
		_, err = fmt.Fprint(os.Stdout, agentSkillMarkdown)
		return err
	},
}
