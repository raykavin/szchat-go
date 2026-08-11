package szchat

import (
	"context"
	"net/url"
)

// Multichannel is a "link multicanal" configuration: a shareable link that
// lets a contact pick from a set of channels.
type Multichannel struct {
	ID              string   `json:"_id,omitempty"`
	Channels        []string `json:"channels"`
	ColorBackground string   `json:"colorBackground"`
	ColorButtons    string   `json:"colorButtons"`
	ColorButtonText string   `json:"colorButtonText"`
	Link            string   `json:"link,omitempty"`
	CreatedAt       string   `json:"created_at,omitempty"`
	UpdatedAt       string   `json:"updated_at,omitempty"`
}

// MultichannelAPI groups the /multichannel endpoints.
type MultichannelAPI struct {
	client *Client
}

// List returns a paginated list of multichannel links via GET /multichannel.
func (a *MultichannelAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[Multichannel], error) {
	var resp PaginatedResponse[Multichannel]
	if err := a.client.get(ctx, "/multichannel", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a multichannel link via POST /multichannel.
func (a *MultichannelAPI) Create(ctx context.Context, m Multichannel) (*Multichannel, error) {
	var created Multichannel
	if err := a.client.post(ctx, "/multichannel", m, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a multichannel link via PUT /multichannel/{id}.
func (a *MultichannelAPI) Update(ctx context.Context, id string, m Multichannel) (*Multichannel, error) {
	var updated Multichannel
	path := "/multichannel/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, m, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete deletes a multichannel link via DELETE /multichannel/{id}.
func (a *MultichannelAPI) Delete(ctx context.Context, id string) error {
	path := "/multichannel/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}
