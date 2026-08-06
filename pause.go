package szchat

import (
	"context"
	"net/url"
)

// Pause is an agent break/pause reason configuration.
type Pause struct {
	ID                 string `json:"_id,omitempty"`
	Name               string `json:"name"`
	Description        string `json:"description,omitempty"`
	InitialTime        int    `json:"initialTime,omitempty"`
	MaxTime            int    `json:"maxTime,omitempty"`
	Active             bool   `json:"active,omitempty"`
	Productive         bool   `json:"productive,omitempty"`
	Message            string `json:"message,omitempty"`
	Cumulative         bool   `json:"cumulative,omitempty"`
	Supervisioned      bool   `json:"supervisioned,omitempty"`
	ContinueAttendance bool   `json:"continue_attendance,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
}

// PauseAPI groups the /pauses endpoints.
type PauseAPI struct {
	client *Client
}

// List returns a paginated list of pauses via GET /pauses.
func (a *PauseAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[Pause], error) {
	var resp PaginatedResponse[Pause]
	if err := a.client.get(ctx, "/pauses", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a pause via POST /pauses.
func (a *PauseAPI) Create(ctx context.Context, pause Pause) (*Pause, error) {
	var created Pause
	if err := a.client.post(ctx, "/pauses", pause, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a pause via PUT /pauses/{id}.
func (a *PauseAPI) Update(ctx context.Context, id string, pause Pause) (*Pause, error) {
	var updated Pause
	path := "/pauses/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, pause, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete deletes a pause via DELETE /pauses/{id}.
func (a *PauseAPI) Delete(ctx context.Context, id string) error {
	path := "/pauses/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}
