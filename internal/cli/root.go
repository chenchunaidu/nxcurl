package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/chenchunaidu/nxcurl/internal/version"
)

var (
	flagEnv      string
	flagJSON     bool
	flagNoHistory bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "nxcurl",
	Short:   "HTTP CLI for humans and coding agents (Postman-like history, envs, imports)",
	Version: version.Version,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagEnv, "env", "e", "", "environment name (variables from ~/.nxcurl/envs/<name>.json)")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "machine-readable JSON output on stdout (for agents)")
	rootCmd.AddCommand(runCmd, sendCmd, historyCmd, envCmd, importCmd, catalogCmd)
}
