package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/chenchunaidu/nxcurl/internal/importers"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import requests from HAR or Postman Collection v2",
}

var importHARCmd = &cobra.Command{
	Use:   "har <file.har>",
	Short: "Import HTTP Archive 1.x",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reqs, err := importers.FromHAR(args[0])
		if err != nil {
			return err
		}
		dir := importers.DefaultImportDir("har")
		out, err := importers.WriteCollection(dir, reqs)
		if err != nil {
			return err
		}
		fmt.Printf("wrote %d requests under %s\n", len(reqs), out)
		return nil
	},
}

var importPostmanCmd = &cobra.Command{
	Use:   "postman <collection.json>",
	Short: "Import Postman Collection v2.0 / v2.1",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reqs, err := importers.FromPostman(args[0])
		if err != nil {
			return err
		}
		dir := importers.DefaultImportDir("postman")
		out, err := importers.WriteCollection(dir, reqs)
		if err != nil {
			return err
		}
		fmt.Printf("wrote %d requests under %s\n", len(reqs), out)
		return nil
	},
}

func init() {
	importCmd.AddCommand(importHARCmd, importPostmanCmd)
}
