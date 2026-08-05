package szchat

import (
	"context"
	"net/url"
)

// Admin is a SZChat administrator account.
type Admin struct {
	ID                        string `json:"_id,omitempty"`
	Type                      string `json:"type,omitempty"`
	Name                      string `json:"name"`
	Email                     string `json:"email"`
	EmailForgotPassword       string `json:"email_forgot_password,omitempty"`
	GroupID                   string `json:"groupId"`
	Language                  string `json:"language,omitempty"`
	SessionToken              string `json:"session_token,omitempty"`
	LoggedAt                  string `json:"logged_at,omitempty"`
	Status                    string `json:"status,omitempty"`
	HasAuthToken              bool   `json:"hasAuthToken,omitempty"`
	SessionAgent              string `json:"session_agent,omitempty"`
	EnableLoginWithRemoteAuth bool   `json:"enable_login_with_remote_auth,omitempty"`
	CreatedAt                 string `json:"created_at,omitempty"`
	UpdatedAt                 string `json:"updated_at,omitempty"`
}

// AdminRequest is the payload for AdminAPI.Create.
type AdminRequest struct {
	Email                     string `json:"email"`
	Name                      string `json:"name"`
	Password                  string `json:"password,omitempty"`
	EmailForgotPassword       string `json:"email_forgot_password,omitempty"`
	EnableLoginWithRemoteAuth *bool  `json:"enable_login_with_remote_auth,omitempty"`
	GroupID                   string `json:"groupId"`
}

// AdminGroup is a permission group administrators can belong to.
type AdminGroup struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Master      int      `json:"master,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

// AdminUpdateResponse is the ack returned by AdminAPI.Update.
type AdminUpdateResponse struct {
	Success bool `json:"success"`
	Date    struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"date"`
}

// AdminAPI groups the /admins endpoints.
type AdminAPI struct {
	client *Client
}

// List returns a paginated list of administrators via GET /admins.
func (a *AdminAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[Admin], error) {
	var resp PaginatedResponse[Admin]
	if err := a.client.get(ctx, "/admins", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListGroups returns a paginated list of administrator groups via
// GET /admins/groups.
func (a *AdminAPI) ListGroups(ctx context.Context, opts ListOptions) (*PaginatedResponse[AdminGroup], error) {
	var resp PaginatedResponse[AdminGroup]
	if err := a.client.get(ctx, "/admins/groups", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates an administrator via POST /admins.
func (a *AdminAPI) Create(ctx context.Context, req AdminRequest) (*Admin, error) {
	var admin Admin
	if err := a.client.post(ctx, "/admins", req, &admin); err != nil {
		return nil, err
	}
	return &admin, nil
}

// Get returns a single administrator via GET /admins/{id}.
func (a *AdminAPI) Get(ctx context.Context, id string) (*Admin, error) {
	var admin Admin
	if err := a.client.get(ctx, "/admins/"+url.PathEscape(id), nil, &admin); err != nil {
		return nil, err
	}
	return &admin, nil
}

// Update partially updates an administrator via PUT /admins/{id}.
func (a *AdminAPI) Update(ctx context.Context, id string, fields map[string]any) (*AdminUpdateResponse, error) {
	var resp AdminUpdateResponse
	if err := a.client.put(ctx, "/admins/"+url.PathEscape(id), fields, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes an administrator via DELETE /admins/{id}.
func (a *AdminAPI) Delete(ctx context.Context, id string) error {
	return a.client.delete(ctx, "/admins/"+url.PathEscape(id), nil)
}
