package szchat

import (
	"context"
	"net/url"
)

// WordFilter is a filtered word entry, applied either to agent messages
// ("agent") or to messages that should not (re)trigger a bot flow
// ("contact").
type WordFilter struct {
	ID        string `json:"_id,omitempty"`
	Word      string `json:"word"`
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// WordFilterAgentAPI groups the /wordFilter/agents endpoints, which filter
// words agents are not allowed to send.
type WordFilterAgentAPI struct {
	client *Client
}

// List returns a paginated list of filtered agent words via
// GET /wordFilter/agents.
func (a *WordFilterAgentAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[WordFilter], error) {
	var resp PaginatedResponse[WordFilter]
	if err := a.client.get(ctx, "/wordFilter/agents", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create adds a filtered agent word via POST /wordFilter/agents.
func (a *WordFilterAgentAPI) Create(ctx context.Context, word string) (*WordFilter, error) {
	var created WordFilter
	req := struct {
		Word string `json:"word"`
	}{Word: word}
	if err := a.client.post(ctx, "/wordFilter/agents", req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a filtered agent word via PUT /wordFilter/agents/{id}.
func (a *WordFilterAgentAPI) Update(ctx context.Context, id, word string) (*WordFilter, error) {
	var updated WordFilter
	req := struct {
		Word string `json:"word"`
	}{Word: word}
	path := "/wordFilter/agents/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes a filtered agent word via DELETE /wordFilter/agents/{id}.
func (a *WordFilterAgentAPI) Delete(ctx context.Context, id string) error {
	path := "/wordFilter/agents/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}

// WordFilterContactAPI groups the /wordFilter/contacts endpoints, which
// filter words that should not (re)trigger a bot flow when sent by a
// contact.
type WordFilterContactAPI struct {
	client *Client
}

// List returns a paginated list of filtered contact words via
// GET /wordFilter/contacts.
func (a *WordFilterContactAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[WordFilter], error) {
	var resp PaginatedResponse[WordFilter]
	if err := a.client.get(ctx, "/wordFilter/contacts", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create adds a filtered contact word via POST /wordFilter/contacts.
func (a *WordFilterContactAPI) Create(ctx context.Context, word string) (*WordFilter, error) {
	var created WordFilter
	req := struct {
		Word string `json:"word"`
	}{Word: word}
	if err := a.client.post(ctx, "/wordFilter/contacts", req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a filtered contact word via PUT /wordFilter/contacts/{id}.
func (a *WordFilterContactAPI) Update(ctx context.Context, id, word string) (*WordFilter, error) {
	var updated WordFilter
	req := struct {
		Word string `json:"word"`
	}{Word: word}
	path := "/wordFilter/contacts/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes a filtered contact word via DELETE /wordFilter/contacts/{id}.
func (a *WordFilterContactAPI) Delete(ctx context.Context, id string) error {
	path := "/wordFilter/contacts/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}
