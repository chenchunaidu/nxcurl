package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/chenchunaidu/nxcurl/internal/executor"
	"github.com/chenchunaidu/nxcurl/internal/history"
)

// AgentJSON is stable structured output for LLM agents.
type AgentJSON struct {
	ID               string            `json:"id"`
	Method           string            `json:"method"`
	URL              string            `json:"url"`
	Status           int               `json:"status"`
	DurationMS       int64             `json:"duration_ms"`
	RequestHeaders   map[string]string `json:"request_headers,omitempty"`
	RequestBody      string            `json:"request_body,omitempty"`
	ResponseHeaders  map[string]string `json:"response_headers,omitempty"`
	ResponseBody     string            `json:"response_body"`
	ResponseBodyJSON any               `json:"response_body_json,omitempty"`
}

func PrintHuman(w io.Writer, r *executor.Result) {
	fmt.Fprintf(w, "→ %s %s\n", r.Method, r.URL)
	fmt.Fprintf(w, "← %d in %s\n", r.Status, r.Duration.Round(time.Millisecond))
	if len(r.RequestHeaders) > 0 {
		fmt.Fprintln(w, "\nRequest headers:")
		printSortedKV(w, r.RequestHeaders)
	}
	if r.RequestBody != "" {
		fmt.Fprintln(w, "\nRequest body:")
		fmt.Fprintln(w, PrettyMaybeJSON(r.RequestBody))
	}
	fmt.Fprintln(w, "\nResponse headers:")
	printSortedKV(w, r.ResponseHeaders)
	fmt.Fprintln(w, "\nBody:")
	fmt.Fprintln(w, PrettyMaybeJSON(string(r.ResponseBody)))
}

func PrintAgentJSON(w io.Writer, r *executor.Result) error {
	bodyStr := string(r.ResponseBody)
	aj := AgentJSON{
		ID:              r.ID,
		Method:          r.Method,
		URL:             r.URL,
		Status:          r.Status,
		DurationMS:      r.Duration.Milliseconds(),
		RequestHeaders:  r.RequestHeaders,
		RequestBody:     r.RequestBody,
		ResponseHeaders: r.ResponseHeaders,
		ResponseBody:    bodyStr,
	}
	var raw any
	if json.Unmarshal(r.ResponseBody, &raw) == nil {
		aj.ResponseBodyJSON = raw
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(aj)
}

func printSortedKV(w io.Writer, m map[string]string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "  %s: %s\n", k, m[k])
	}
}

// PrettyMaybeJSON indents JSON when the input is valid JSON; otherwise returns s trimmed.
func PrettyMaybeJSON(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "(empty)"
	}
	var v any
	if json.Unmarshal([]byte(s), &v) != nil {
		return s
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return s
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

func PrintHistoryHuman(w io.Writer, entries []history.Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No history yet.")
		return
	}
	for _, e := range entries {
		fmt.Fprintf(w, "%s  %s  %s  → %d  %s\n",
			e.ID, e.TS.Format(time.RFC3339), e.Method, e.Status, truncate(e.URL, 72))
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Stdout returns os.Stdout; split for tests.
var Stdout io.Writer = os.Stdout
