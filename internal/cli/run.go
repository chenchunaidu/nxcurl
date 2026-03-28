package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/chenchunaidu/nxcurl/internal/executor"
	"github.com/chenchunaidu/nxcurl/internal/history"
	"github.com/chenchunaidu/nxcurl/internal/output"
)

var (
	runMethod string
	runHeader []string
	runData   string
)

var runCmd = &cobra.Command{
	Use:   "run <url>",
	Short: "Execute an HTTP request (curl-style)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		headers, err := parseHeaders(runHeader)
		if err != nil {
			return err
		}
		method := runMethod
		if method == "" {
			method = "GET"
		}
		spec := executor.RequestSpec{
			Method:  method,
			URL:     args[0],
			Headers: headers,
			Body:    runData,
			EnvName: flagEnv,
		}
		return runAndPrint(spec, "", "")
	},
}

var sendCmd = &cobra.Command{
	Use:   "send <request.json>",
	Short: "Run a saved request definition (from import or hand-written JSON)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}
		var sr struct {
			Name    string            `json:"name"`
			Method  string            `json:"method"`
			URL     string            `json:"url"`
			Headers map[string]string `json:"headers"`
			Body    string            `json:"body"`
		}
		if err := json.Unmarshal(b, &sr); err != nil {
			return err
		}
		spec := executor.RequestSpec{
			Method:  sr.Method,
			URL:     sr.URL,
			Headers: sr.Headers,
			Body:    sr.Body,
			EnvName: flagEnv,
		}
		return runAndPrint(spec, "", sr.Name)
	},
}

func init() {
	runCmd.Flags().StringVarP(&runMethod, "request", "X", "", "HTTP method (GET, POST, ...)")
	runCmd.Flags().StringArrayVarP(&runHeader, "header", "H", nil, "header 'Name: value' (repeatable)")
	runCmd.Flags().StringVarP(&runData, "data", "d", "", "request body")
	runCmd.Flags().BoolVar(&flagNoHistory, "no-history", false, "do not save this request to history")
	sendCmd.Flags().BoolVar(&flagNoHistory, "no-history", false, "do not save this request to history")
}

func parseHeaders(lines []string) (map[string]string, error) {
	m := map[string]string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx <= 0 {
			return nil, fmt.Errorf("invalid header %q (want 'Name: value')", line)
		}
		k := strings.TrimSpace(line[:idx])
		v := strings.TrimSpace(line[idx+1:])
		m[k] = v
	}
	return m, nil
}

func runAndPrint(spec executor.RequestSpec, collection, requestName string) error {
	res, err := executor.Run(spec)
	if err != nil {
		return err
	}
	if !flagNoHistory {
		_ = history.Append(history.Entry{
			ID:              res.ID,
			TS:              time.Now().UTC(),
			Method:          res.Method,
			URL:             res.URL,
			RequestHeaders:  res.RequestHeaders,
			RequestBody:     res.RequestBody,
			Status:          res.Status,
			ResponseHeaders: res.ResponseHeaders,
			ResponseBody:    string(res.ResponseBody),
			DurationMS:      res.Duration.Milliseconds(),
			Collection:      collection,
			RequestName:     requestName,
		})
	}
	if flagJSON {
		return output.PrintAgentJSON(os.Stdout, res)
	}
	output.PrintHuman(os.Stdout, res)
	return nil
}
