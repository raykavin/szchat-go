package szchat

import (
	"context"
	"net/url"
)

// Tabulation is a classification label applied when closing an attendance.
type Tabulation struct {
	ID        string `json:"_id,omitempty"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type tabulationRequest struct {
	Name string `json:"name"`
}

// TabulationAPI groups the /tabulations endpoints.
type TabulationAPI struct {
	client *Client
}

// List returns a paginated list of tabulations via GET /tabulations.
func (a *TabulationAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[Tabulation], error) {
	var resp PaginatedResponse[Tabulation]
	if err := a.client.get(ctx, "/tabulations", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a tabulation via POST /tabulations.
func (a *TabulationAPI) Create(ctx context.Context, name string) (*Tabulation, error) {
	var tabulation Tabulation
	if err := a.client.post(ctx, "/tabulations", tabulationRequest{Name: name}, &tabulation); err != nil {
		return nil, err
	}
	return &tabulation, nil
}

// Update renames a tabulation via PUT /tabulations/{id}.
func (a *TabulationAPI) Update(ctx context.Context, id, name string) (*Tabulation, error) {
	var tabulation Tabulation
	path := "/tabulations/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, tabulationRequest{Name: name}, &tabulation); err != nil {
		return nil, err
	}
	return &tabulation, nil
}

// Delete deletes a tabulation via DELETE /tabulations/{id}.
func (a *TabulationAPI) Delete(ctx context.Context, id string) error {
	path := "/tabulations/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}
