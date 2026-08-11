package szchat

import "context"

// PlaceholderResolveRequest is the payload for PlaceholderAPI.Resolve.
type PlaceholderResolveRequest struct {
	ContactID          string   `json:"contact_id"`
	AgentID            string   `json:"agent_id,omitempty"`
	SessionID          string   `json:"session_id,omitempty"`
	PlaceholdersParams []string `json:"placeholders_params"`
}

// PlaceholderAPI groups the message placeholder resolution endpoint.
type PlaceholderAPI struct {
	client *Client
}

// Resolve resolves a list of placeholders (e.g. "{{NAME}}") to their values
// for a given contact/agent/session via POST /user/agent/placeholders. The
// returned slice preserves the order of req.PlaceholdersParams.
func (a *PlaceholderAPI) Resolve(ctx context.Context, req PlaceholderResolveRequest) ([]string, error) {
	var values []string
	if err := a.client.post(ctx, "/user/agent/placeholders", req, &values); err != nil {
		return nil, err
	}
	return values, nil
}
