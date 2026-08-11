package szchat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// APIVersionResponse is the payload returned by Client.APIVersion.
type APIVersionResponse struct {
	Version string `json:"version"`
}

// APIVersion returns the API version reported by GET /api/version. This
// endpoint sits outside the versioned /api/v4 path this client otherwise
// targets and requires no authentication.
func (c *Client) APIVersion(ctx context.Context) (*APIVersionResponse, error) {
	root := strings.TrimSuffix(c.baseURL, apiVersionPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, root+"/api/version", nil)
	if err != nil {
		return nil, fmt.Errorf("szchat: building request: %w", err)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, parseAPIError(resp.StatusCode, resp.Body)
	}

	var out APIVersionResponse
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, &out); err != nil {
			return nil, fmt.Errorf("szchat: decoding response body: %w", err)
		}
	}
	return &out, nil
}
