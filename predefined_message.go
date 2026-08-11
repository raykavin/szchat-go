package szchat

import (
	"context"
	"net/url"
)

// PredefinedMessagePart is a single part of a predefined message. Type is
// set by the server (observed values: "text") and is not sent on
// create/update.
type PredefinedMessagePart struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message"`
}

// PredefinedMessage is a canned reply agents can insert into a
// conversation.
type PredefinedMessage struct {
	ID          string                  `json:"_id,omitempty"`
	Message     []PredefinedMessagePart `json:"message"`
	Description string                  `json:"description"`
	CreatedAt   string                  `json:"created_at,omitempty"`
	UpdatedAt   string                  `json:"updated_at,omitempty"`
}

// PredefinedMessageAPI groups the /predefined_messages endpoints.
type PredefinedMessageAPI struct {
	client *Client
}

// List returns a paginated list of predefined messages via
// GET /predefined_messages.
func (a *PredefinedMessageAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[PredefinedMessage], error) {
	var resp PaginatedResponse[PredefinedMessage]
	if err := a.client.get(ctx, "/predefined_messages", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a single predefined message via GET /predefined_messages/{id}.
func (a *PredefinedMessageAPI) Get(ctx context.Context, id string) (*PredefinedMessage, error) {
	var msg PredefinedMessage
	path := "/predefined_messages/" + url.PathEscape(id)
	if err := a.client.get(ctx, path, nil, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// Create creates a predefined message via POST /predefined_messages.
func (a *PredefinedMessageAPI) Create(ctx context.Context, msg PredefinedMessage) (*PredefinedMessage, error) {
	var created PredefinedMessage
	if err := a.client.post(ctx, "/predefined_messages", msg, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a predefined message via PUT /predefined_messages/{id}.
func (a *PredefinedMessageAPI) Update(ctx context.Context, id string, msg PredefinedMessage) (*PredefinedMessage, error) {
	var updated PredefinedMessage
	path := "/predefined_messages/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, msg, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete deletes a predefined message via DELETE /predefined_messages/{id}.
func (a *PredefinedMessageAPI) Delete(ctx context.Context, id string) error {
	path := "/predefined_messages/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}
