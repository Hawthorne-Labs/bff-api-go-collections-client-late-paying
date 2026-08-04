package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ForwardRequest forwards a request to the core-api and returns the response.
func (c *CoreClient) ForwardRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	targetURL := c.BaseURL.String() + path
	req, err := http.NewRequestWithContext(nil, method, targetURL, reqBody) // context set by caller
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add trace propagation headers
	req.Header.Set("X-Trace-Id", "trace-local-dev")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("forward request to core: %w", err)
	}
	return resp, nil
}
