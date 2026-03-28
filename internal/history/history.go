package history

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/chenchunaidu/nxcurl/internal/paths"
)

// Entry is one saved request/response pair.
type Entry struct {
	ID               string            `json:"id"`
	TS               time.Time         `json:"ts"`
	Method           string            `json:"method"`
	URL              string            `json:"url"`
	RequestHeaders   map[string]string `json:"request_headers,omitempty"`
	RequestBody      string            `json:"request_body,omitempty"`
	Status           int               `json:"status"`
	ResponseHeaders  map[string]string `json:"response_headers,omitempty"`
	ResponseBody     string            `json:"response_body,omitempty"`
	DurationMS       int64             `json:"duration_ms"`
	Collection       string            `json:"collection,omitempty"`
	RequestName      string            `json:"request_name,omitempty"`
}

func Append(e Entry) error {
	p, err := paths.HistoryPath()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return err
	}
	return nil
}

// List returns entries newest-first (last lines in file are newest — we read all and reverse).
func List(limit int) ([]Entry, error) {
	p, err := paths.HistoryPath()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var all []Entry
	sc := bufio.NewScanner(f)
	// Avoid unbounded memory on huge files: read in chunks if needed later.
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			continue
		}
		all = append(all, e)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	// Newest last in file → reverse
	for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
		all[i], all[j] = all[j], all[i]
	}
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}

func FindByID(id string) (*Entry, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	all, err := List(0)
	if err != nil {
		return nil, err
	}
	for _, e := range all {
		if e.ID == id {
			cp := e
			return &cp, nil
		}
	}
	return nil, errors.New("history entry not found")
}
