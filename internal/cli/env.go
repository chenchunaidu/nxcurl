package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/chenchunaidu/nxcurl/internal/paths"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments (~/.nxcurl/envs/<name>.json)",
}

var envInitCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Create an empty environment file (flat JSON object of string keys)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.EnvFile(args[0])
		if err != nil {
			return err
		}
		if _, err := os.Stat(p); err == nil {
			return fmt.Errorf("already exists: %s", p)
		}
		if _, err := paths.EnsureDataDir(); err != nil {
			return err
		}
		b, err := json.MarshalIndent(map[string]string{}, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile(p, b, 0o644)
	},
}

var envPathCmd = &cobra.Command{
	Use:   "path <name>",
	Short: "Print absolute path to an environment file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.EnvFile(args[0])
		if err != nil {
			return err
		}
		fmt.Println(p)
		return nil
	},
}

func init() {
	envCmd.AddCommand(envInitCmd, envPathCmd)
}
