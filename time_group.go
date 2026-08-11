package szchat

import (
	"context"
	"net/url"
)

// TimeGroupRange is a single day/week/month/time range within a TimeGroup.
type TimeGroupRange struct {
	Hour         string `json:"hour"`
	Description  string `json:"description"`
	InitialDay   string `json:"initialDay"`
	FinalDay     string `json:"finalDay"`
	InitialMonth int    `json:"initialMonth"`
	FinalMonth   int    `json:"finalMonth"`
	InitialWeek  int    `json:"initialWeek"`
	FinalWeek    int    `json:"finalWeek"`
	InitialTime  string `json:"initialTime"`
	FinalTime    string `json:"finalTime"`
}

// TimeGroup ("grupo de horários") is a named set of time ranges used to
// gate channel/flow availability.
type TimeGroup struct {
	ID        string           `json:"_id,omitempty"`
	Name      string           `json:"name"`
	Group     []TimeGroupRange `json:"group"`
	CreatedAt string           `json:"created_at,omitempty"`
	UpdatedAt string           `json:"updated_at,omitempty"`
}

// TimeGroupAPI groups the /timeGroup endpoints.
type TimeGroupAPI struct {
	client *Client
}

// List returns a paginated list of time groups via GET /timeGroup.
func (a *TimeGroupAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[TimeGroup], error) {
	var resp PaginatedResponse[TimeGroup]
	if err := a.client.get(ctx, "/timeGroup", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a time group via POST /timeGroup.
func (a *TimeGroupAPI) Create(ctx context.Context, group TimeGroup) (*TimeGroup, error) {
	var created TimeGroup
	if err := a.client.post(ctx, "/timeGroup", group, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a time group via PUT /timeGroup/{id}.
func (a *TimeGroupAPI) Update(ctx context.Context, id string, group TimeGroup) (*TimeGroup, error) {
	var updated TimeGroup
	path := "/timeGroup/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, group, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete deletes a time group via DELETE /timeGroup/{id}.
func (a *TimeGroupAPI) Delete(ctx context.Context, id string) error {
	path := "/timeGroup/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}
