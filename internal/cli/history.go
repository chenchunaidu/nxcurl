package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/chenchunaidu/nxcurl/internal/executor"
	"github.com/chenchunaidu/nxcurl/internal/history"
	"github.com/chenchunaidu/nxcurl/internal/output"
)

var historyLimit int

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "List or inspect past requests",
}

var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent history (newest first)",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := history.List(historyLimit)
		if err != nil {
			return err
		}
		if flagJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(entries)
		}
		output.PrintHistoryHuman(os.Stdout, entries)
		return nil
	},
}

var historyShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show one history entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		e, err := history.FindByID(args[0])
		if err != nil {
			return err
		}
		if flagJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(e)
		}
		fmt.Fprintf(os.Stdout, "id %s  %s  %s %s  status %d  %dms\n",
			e.ID, e.TS.Format("2006-01-02 15:04:05"), e.Method, e.URL, e.Status, e.DurationMS)
		if len(e.RequestHeaders) > 0 {
			fmt.Fprintln(os.Stdout, "\nRequest headers:")
			for k, v := range e.RequestHeaders {
				fmt.Fprintf(os.Stdout, "  %s: %s\n", k, v)
			}
		}
		if e.RequestBody != "" {
			fmt.Fprintln(os.Stdout, "\nRequest body:")
			fmt.Fprintln(os.Stdout, output.PrettyMaybeJSON(e.RequestBody))
		}
		fmt.Fprintln(os.Stdout, "\nResponse headers:")
		for k, v := range e.ResponseHeaders {
			fmt.Fprintf(os.Stdout, "  %s: %s\n", k, v)
		}
		fmt.Fprintln(os.Stdout, "\nResponse body:")
		fmt.Fprintln(os.Stdout, output.PrettyMaybeJSON(e.ResponseBody))
		return nil
	},
}

var historyReplayCmd = &cobra.Command{
	Use:   "replay <id>",
	Short: "Re-run a history entry (saves a new history row)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		e, err := history.FindByID(args[0])
		if err != nil {
			return err
		}
		spec := executor.RequestSpec{
			Method:  e.Method,
			URL:     e.URL,
			Headers: e.RequestHeaders,
			Body:    e.RequestBody,
			EnvName: flagEnv,
		}
		return runAndPrint(spec, e.Collection, e.RequestName)
	},
}

func init() {
	historyListCmd.Flags().IntVar(&historyLimit, "limit", 50, "max rows")
	historyCmd.AddCommand(historyListCmd, historyShowCmd, historyReplayCmd)
}
