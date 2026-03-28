package paths

import (
	"os"
	"path/filepath"
)

func DataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(home, ".nxcurl")
	return d, nil
}

func EnsureDataDir() (string, error) {
	d, err := DataDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(d, "envs"), 0o755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(d, "collections"), 0o755); err != nil {
		return "", err
	}
	return d, nil
}

func HistoryPath() (string, error) {
	d, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "history.jsonl"), nil
}

func EnvFile(name string) (string, error) {
	d, err := EnsureDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "envs", name+".json"), nil
}
