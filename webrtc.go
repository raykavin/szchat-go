package szchat

import (
	"context"
	"net/url"
)

// WebRTCConfig is the payload returned by WebRTCAPI.Get.
type WebRTCConfig struct {
	Enabled             bool   `json:"enabled"`
	URL                 string `json:"url,omitempty"`
	Width               int    `json:"width,omitempty"`
	Height              int    `json:"height,omitempty"`
	RefreshWindowButton bool   `json:"refreshWindowButton,omitempty"`
	CloseWindowButton   bool   `json:"closeWindowButton,omitempty"`
}

// WebRTCAPI groups the agent WebRTC configuration endpoint.
type WebRTCAPI struct {
	client *Client
}

// Get returns an agent's WebRTC (Callcenter) configuration via
// GET /user/agent/webrtc/{agent_id}.
func (a *WebRTCAPI) Get(ctx context.Context, agentID string) (*WebRTCConfig, error) {
	var cfg WebRTCConfig
	path := "/user/agent/webrtc/" + url.PathEscape(agentID)
	if err := a.client.get(ctx, path, nil, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
