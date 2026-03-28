package executor

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chenchunaidu/nxcurl/internal/env"
)

type RequestSpec struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	EnvName string
}

type Result struct {
	ID              string
	Method          string
	URL             string
	RequestHeaders  map[string]string
	RequestBody     string
	Status          int
	ResponseHeaders map[string]string
	ResponseBody    []byte
	Duration        time.Duration
}

func Run(spec RequestSpec) (*Result, error) {
	vars, err := env.Load(spec.EnvName)
	if err != nil {
		return nil, err
	}
	method := strings.ToUpper(strings.TrimSpace(spec.Method))
	if method == "" {
		method = http.MethodGet
	}
	rawURL := env.Subst(strings.TrimSpace(spec.URL), vars)
	if rawURL == "" {
		return nil, fmt.Errorf("url is required")
	}

	reqHeaders := map[string]string{}
	for k, v := range spec.Headers {
		kk := strings.TrimSpace(k)
		if kk == "" {
			continue
		}
		reqHeaders[kk] = env.Subst(v, vars)
	}
	bodyStr := env.Subst(spec.Body, vars)

	var bodyReader io.Reader
	if bodyStr != "" {
		bodyReader = strings.NewReader(bodyStr)
	}
	req, err := http.NewRequest(method, rawURL, bodyReader)
	if err != nil {
		return nil, err
	}
	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}
	if bodyStr != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 120 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	dur := time.Since(start)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respHeaders := map[string]string{}
	for k, vals := range resp.Header {
		if len(vals) > 0 {
			respHeaders[k] = strings.Join(vals, ", ")
		}
	}

	id, err := shortID()
	if err != nil {
		return nil, err
	}

	return &Result{
		ID:              id,
		Method:          method,
		URL:             rawURL,
		RequestHeaders:  reqHeaders,
		RequestBody:     bodyStr,
		Status:          resp.StatusCode,
		ResponseHeaders: respHeaders,
		ResponseBody:    respBody,
		Duration:        dur,
	}, nil
}

func shortID() (string, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
