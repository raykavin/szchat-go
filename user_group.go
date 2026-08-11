package szchat

import (
	"context"
	"net/url"
)

// UserGroup ("grupo de usuários") is a permission group assignable to
// admins/agents, distinct from AdminGroup which is read-only from
// AdminAPI.ListGroups.
type UserGroup struct {
	ID                   string   `json:"_id,omitempty"`
	Name                 string   `json:"name"`
	Master               bool     `json:"master"`
	Permissions          []string `json:"permissions,omitempty"`
	ViewByTeams          bool     `json:"view_by_teams,omitempty"`
	GeneralSelectedTeams []string `json:"general_selected_teams,omitempty"`
	CreatedAt            string   `json:"created_at,omitempty"`
	UpdatedAt            string   `json:"updated_at,omitempty"`
}

// UserGroupAPI groups the /userGroup endpoints.
type UserGroupAPI struct {
	client *Client
}

// List returns a paginated list of user groups via GET /userGroup.
func (a *UserGroupAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[UserGroup], error) {
	var resp PaginatedResponse[UserGroup]
	if err := a.client.get(ctx, "/userGroup", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a user group via POST /userGroup.
func (a *UserGroupAPI) Create(ctx context.Context, group UserGroup) (*UserGroup, error) {
	var created UserGroup
	if err := a.client.post(ctx, "/userGroup", group, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a user group via PUT /userGroup/{id}.
func (a *UserGroupAPI) Update(ctx context.Context, id string, group UserGroup) (*UserGroup, error) {
	var updated UserGroup
	path := "/userGroup/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, group, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete deletes a user group via DELETE /userGroup/{id}.
func (a *UserGroupAPI) Delete(ctx context.Context, id string) error {
	path := "/userGroup/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}
