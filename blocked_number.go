package szchat

import (
	"context"
	"net/url"
)

// BlockedNumber is a phone number blocked from starting new attendances.
type BlockedNumber struct {
	ID        string `json:"_id,omitempty"`
	Number    string `json:"number"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// BlockedNumberAPI groups the /blockedNumbers endpoints.
type BlockedNumberAPI struct {
	client *Client
}

// List returns a paginated list of blocked numbers via GET /blockedNumbers.
func (a *BlockedNumberAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[BlockedNumber], error) {
	var resp PaginatedResponse[BlockedNumber]
	if err := a.client.get(ctx, "/blockedNumbers", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create blocks a number via POST /blockedNumbers.
func (a *BlockedNumberAPI) Create(ctx context.Context, number string) (*BlockedNumber, error) {
	var created BlockedNumber
	req := struct {
		Number string `json:"number"`
	}{Number: number}
	if err := a.client.post(ctx, "/blockedNumbers", req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a blocked number via PUT /blockedNumbers/{id}.
func (a *BlockedNumberAPI) Update(ctx context.Context, id, number string) (*BlockedNumber, error) {
	var updated BlockedNumber
	req := struct {
		Number string `json:"number"`
	}{Number: number}
	path := "/blockedNumbers/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete unblocks a number via DELETE /blockedNumbers/{id}.
func (a *BlockedNumberAPI) Delete(ctx context.Context, id string) error {
	path := "/blockedNumbers/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}
