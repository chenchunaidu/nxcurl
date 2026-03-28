package env

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/chenchunaidu/nxcurl/internal/paths"
)

var varRe = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.-]+)\s*\}\}`)

// Load reads ~/.nxcurl/envs/<name>.json as a flat string map.
func Load(name string) (map[string]string, error) {
	if name == "" {
		return nil, nil
	}
	p, err := paths.EnvFile(name)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if m == nil {
		return map[string]string{}, nil
	}
	return m, nil
}

// Subst replaces {{KEY}} in s using vars (missing keys become empty string).
func Subst(s string, vars map[string]string) string {
	return varRe.ReplaceAllStringFunc(s, func(match string) string {
		sub := varRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		key := strings.TrimSpace(sub[1])
		if vars == nil {
			return ""
		}
		return vars[key]
	})
}
