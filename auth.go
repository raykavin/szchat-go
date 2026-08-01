package szchat

import (
	"context"
	"net/url"
)

// User represents an authenticated agent or admin, as returned by the
// authentication endpoints.
type User struct {
	ID           string   `json:"_id"`
	Type         string   `json:"type"`
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Codename     string   `json:"codename,omitempty"`
	Ramal        string   `json:"ramal,omitempty"`
	Begin        string   `json:"begin,omitempty"`
	End          string   `json:"end,omitempty"`
	Campaigns    []string `json:"campaigns,omitempty"`
	SessionToken string   `json:"session_token,omitempty"`
	Status       string   `json:"status,omitempty"`
	GroupID      string   `json:"groupId,omitempty"`
	Photo        string   `json:"photo,omitempty"`
	UpdatedAt    string   `json:"updated_at,omitempty"`
	CreatedAt    string   `json:"created_at,omitempty"`
}

type loginRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DeviceToken string `json:"device_token,omitempty"`
}

// LoginResponse is the payload returned by POST /auth/login.
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// RefreshResponse is the payload returned by GET /auth/refresh.
type RefreshResponse struct {
	Token string `json:"token"`
}

// Login authenticates against POST /auth/login using the credentials the
// Client was created with, and stores the returned bearer token for
// subsequent requests. It is normally unnecessary to call this directly:
// the Client authenticates lazily on the first request and re-authenticates
// automatically when the token expires.
func (c *Client) Login(ctx context.Context) (*LoginResponse, error) {
	var resp LoginResponse
	req := loginRequest{
		Email:       c.email,
		Password:    c.password,
		DeviceToken: c.deviceToken,
	}
	if err := c.post(withNoAuthRetry(ctx), "/auth/login", req, &resp); err != nil {
		return nil, err
	}
	c.setToken(resp.Token)
	return &resp, nil
}

// Me returns the authenticated user's profile via GET /auth/me.
func (c *Client) Me(ctx context.Context) (*User, error) {
	var me User
	if err := c.get(ctx, "/auth/me", nil, &me); err != nil {
		return nil, err
	}
	return &me, nil
}

// Logout closes the current session via GET /auth/logout. scope selects
// which sessions to close: "api" (default when empty), "web", or "all".
func (c *Client) Logout(ctx context.Context, scope string) error {
	if scope == "" {
		scope = "api"
	}
	q := url.Values{"type": []string{scope}}
	if err := c.get(withNoAuthRetry(ctx), "/auth/logout", q, nil); err != nil {
		return err
	}
	c.setToken("")
	return nil
}

// refreshToken renews the bearer token, called automatically by do() on a
// 401 response. If the client has never authenticated it performs a full
// login; otherwise it calls GET /auth/refresh, falling back to a full login
// if the refresh itself is unauthorized.
func (c *Client) refreshToken(ctx context.Context) error {
	if c.getToken() == "" {
		_, err := c.Login(ctx)
		return err
	}

	var resp RefreshResponse
	err := c.get(withNoAuthRetry(ctx), "/auth/refresh", nil, &resp)
	if err == nil {
		c.setToken(resp.Token)
		return nil
	}

	if IsUnauthorized(err) {
		_, loginErr := c.Login(ctx)
		return loginErr
	}
	return err
}
