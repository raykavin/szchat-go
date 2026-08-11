package szchat

import (
	"context"
	"net/url"
)

// AgentPauseStartRequest is the payload for AgentPauseAPI.Start.
type AgentPauseStartRequest struct {
	PauseID             string `json:"pause_id"`
	IsPausePersonalized bool   `json:"isPausePersonalized,omitempty"`
	IsNewPause          bool   `json:"is_new_pause,omitempty"`
}

// AgentPauseStartResponse is the payload returned by AgentPauseAPI.Start.
type AgentPauseStartResponse struct {
	Message string `json:"message"`
	Pause   Pause  `json:"pause"`
}

// AgentPauseProgress is the payload returned by AgentPauseAPI.Progress.
type AgentPauseProgress struct {
	Status int `json:"status"`
	Pause  struct {
		Pause
		StartTime      string `json:"startTime,omitempty"`
		CumulativeTime string `json:"cumulative_time,omitempty"`
	} `json:"pause"`
}

// AgentPauseAPI groups the authenticated agent's own pause endpoints
// (/user/agents/pauses*), distinct from PauseAPI which manages the tenant's
// pause reason catalog.
type AgentPauseAPI struct {
	client *Client
}

// List returns the pause reason catalog available to the authenticated
// agent via GET /user/agents/pauses.
func (a *AgentPauseAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[Pause], error) {
	var resp PaginatedResponse[Pause]
	if err := a.client.get(ctx, "/user/agents/pauses", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Start begins a pause for the authenticated agent via
// POST /user/agents/pauses/start.
func (a *AgentPauseAPI) Start(ctx context.Context, req AgentPauseStartRequest) (*AgentPauseStartResponse, error) {
	var resp AgentPauseStartResponse
	if err := a.client.post(ctx, "/user/agents/pauses/start", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Stop ends the authenticated agent's current pause via
// POST /user/agents/pauses/stop.
func (a *AgentPauseAPI) Stop(ctx context.Context, pauseID string) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	req := struct {
		PauseID string `json:"pause_id"`
	}{PauseID: pauseID}
	if err := a.client.post(ctx, "/user/agents/pauses/stop", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}

// Progress returns the authenticated agent's current pause progress via
// GET /user/agents/pauses/progress.
func (a *AgentPauseAPI) Progress(ctx context.Context, newVersion bool) (*AgentPauseProgress, error) {
	q := url.Values{}
	if newVersion {
		q.Set("new_version", "true")
	}
	var resp AgentPauseProgress
	if err := a.client.get(ctx, "/user/agents/pauses/progress", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
