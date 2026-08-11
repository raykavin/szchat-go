package szchat

import (
	"context"
	"net/url"
)

// TagCategory is a tag category available to agents/sessions.
type TagCategory struct {
	ID        string   `json:"_id"`
	Name      string   `json:"name"`
	Type      string   `json:"type,omitempty"`
	Campaigns []string `json:"campaigns,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
}

// TagAPI groups the session tag endpoints
// (/user/agents/list/tagsCategory, /user/agents/session/*TagCategory).
type TagAPI struct {
	client *Client
}

// ListCategories lists the available tag categories via
// GET /user/agents/list/tagsCategory.
func (a *TagAPI) ListCategories(ctx context.Context, name string, paginate bool) ([]TagCategory, error) {
	q := url.Values{}
	setParam(q, "name", name)
	if paginate {
		q.Set("paginate", "true")
	}
	var categories []TagCategory
	if err := a.client.get(ctx, "/user/agents/list/tagsCategory", q, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

// SetSessionTag assigns a tag to a session via
// POST /user/agents/session/setTagCategory.
func (a *TagAPI) SetSessionTag(ctx context.Context, sessionID, tagID string) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	req := struct {
		SessionID string `json:"session_id"`
		TagID     string `json:"tag_id"`
	}{SessionID: sessionID, TagID: tagID}
	if err := a.client.post(ctx, "/user/agents/session/setTagCategory", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}

// RemoveSessionTag removes a session's assigned tag via
// POST /user/agents/session/deleteTagCategory.
func (a *TagAPI) RemoveSessionTag(ctx context.Context, sessionID, tagID string) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	req := struct {
		SessionID string `json:"session_id"`
		TagID     string `json:"tag_id"`
	}{SessionID: sessionID, TagID: tagID}
	if err := a.client.post(ctx, "/user/agents/session/deleteTagCategory", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}
