package importers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chenchunaidu/nxcurl/internal/paths"
)

// SavedRequest is a portable request definition (collection item).
type SavedRequest struct {
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

// FromHAR reads a HAR 1.x file and returns requests.
func FromHAR(path string) ([]SavedRequest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Log struct {
			Entries []struct {
				Request struct {
					Method   string `json:"method"`
					URL      string `json:"url"`
					Headers  []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"headers"`
					PostData *struct {
						Text string `json:"text"`
					} `json:"postData"`
				} `json:"request"`
			} `json:"entries"`
		} `json:"log"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	var out []SavedRequest
	for i, e := range doc.Log.Entries {
		r := e.Request
		if strings.TrimSpace(r.URL) == "" {
			continue
		}
		h := map[string]string{}
		for _, hh := range r.Headers {
			n := strings.TrimSpace(hh.Name)
			if n == "" || strings.HasPrefix(strings.ToLower(n), ":") {
				continue
			}
			ln := strings.ToLower(n)
			if ln == "content-length" || ln == "host" || ln == "connection" {
				continue
			}
			h[n] = hh.Value
		}
		body := ""
		if r.PostData != nil {
			body = r.PostData.Text
		}
		out = append(out, SavedRequest{
			Name:    fmt.Sprintf("har-%d", i+1),
			Method:  r.Method,
			URL:     r.URL,
			Headers: h,
			Body:    body,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no usable entries in HAR")
	}
	return out, nil
}

// FromPostman reads a Postman Collection v2.0/v2.1 JSON file.
func FromPostman(path string) ([]SavedRequest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var root map[string]any
	if err := json.Unmarshal(b, &root); err != nil {
		return nil, err
	}
	items, _ := root["item"].([]any)
	if len(items) == 0 {
		return nil, fmt.Errorf("no items in Postman collection")
	}
	var out []SavedRequest
	var walk func(prefix string, node any)
	walk = func(prefix string, node any) {
		m, ok := node.(map[string]any)
		if !ok {
			return
		}
		if req, ok := m["request"].(map[string]any); ok {
			name, _ := m["name"].(string)
			if prefix != "" {
				name = prefix + "/" + name
			}
			sr := postmanRequestToSaved(name, req)
			if sr.URL != "" {
				out = append(out, sr)
			}
			return
		}
		name, _ := m["name"].(string)
		childPrefix := prefix
		if name != "" {
			if prefix == "" {
				childPrefix = name
			} else {
				childPrefix = prefix + "/" + name
			}
		}
		if arr, ok := m["item"].([]any); ok {
			for _, ch := range arr {
				walk(childPrefix, ch)
			}
		}
	}
	for _, it := range items {
		walk("", it)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no requests found in Postman collection")
	}
	return out, nil
}

func postmanRequestToSaved(name string, req map[string]any) SavedRequest {
	method, _ := req["method"].(string)
	urlStr := ""
	switch u := req["url"].(type) {
	case string:
		urlStr = u
	case map[string]any:
		if raw, ok := u["raw"].(string); ok {
			urlStr = raw
		}
	}
	h := map[string]string{}
	if hdrs, ok := req["header"].([]any); ok {
		for _, x := range hdrs {
			hm, ok := x.(map[string]any)
			if !ok {
				continue
			}
			k, _ := hm["key"].(string)
			v, _ := hm["value"].(string)
			if k != "" && !strings.EqualFold(k, "content-length") {
				h[k] = v
			}
		}
	}
	body := ""
	if bm, ok := req["body"].(map[string]any); ok {
		if mode, _ := bm["mode"].(string); mode == "raw" {
			body, _ = bm["raw"].(string)
		}
	}
	return SavedRequest{Name: name, Method: method, URL: urlStr, Headers: h, Body: body}
}

// WriteCollection saves requests as JSON files under ~/.nxcurl/collections/<dir>/.
func WriteCollection(dir string, reqs []SavedRequest) (string, error) {
	base, err := paths.EnsureDataDir()
	if err != nil {
		return "", err
	}
	colDir := filepath.Join(base, "collections", dir)
	if err := os.MkdirAll(colDir, 0o755); err != nil {
		return "", err
	}
	for i, r := range reqs {
		safe := sanitizeFilename(r.Name)
		if safe == "" {
			safe = fmt.Sprintf("req-%d", i+1)
		}
		fn := filepath.Join(colDir, safe+".json")
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(fn, b, 0o644); err != nil {
			return "", err
		}
	}
	return colDir, nil
}

func sanitizeFilename(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	out := strings.Trim(b.String(), "_")
	if len(out) > 80 {
		out = out[:80]
	}
	return out
}

// DefaultImportDir returns a timestamped directory name for imports.
func DefaultImportDir(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, time.Now().UTC().Format("20060102-150405"))
}
