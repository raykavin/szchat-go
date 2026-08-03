package szchat

import (
	"context"
	"net/url"
)

// ContactGroup is a contact segmentation group.
type ContactGroup struct {
	ID        string `json:"_id"`
	NameGroup string `json:"nameGroup"`
	OptIn     any    `json:"opt_in,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ContactGroupSummary is the lightweight shape returned when listing the
// groups a specific contact belongs to.
type ContactGroupSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type contactGroupRequest struct {
	NameGroup string `json:"nameGroup"`
}

// ContactGroupAPI groups the /contacts/groups endpoints.
type ContactGroupAPI struct {
	client *Client
}

// List returns a paginated list of contact groups via GET /contacts/groups.
func (a *ContactGroupAPI) List(ctx context.Context, opts ListOptions) (*PaginatedResponse[ContactGroup], error) {
	var resp PaginatedResponse[ContactGroup]
	if err := a.client.get(ctx, "/contacts/groups", opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a contact group via POST /contacts/groups.
func (a *ContactGroupAPI) Create(ctx context.Context, name string) (*ContactGroup, error) {
	var group ContactGroup
	if err := a.client.post(ctx, "/contacts/groups", contactGroupRequest{NameGroup: name}, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

// Update renames a contact group via PUT /contacts/groups/{id}.
func (a *ContactGroupAPI) Update(ctx context.Context, id, name string) (*ContactGroup, error) {
	var group ContactGroup
	path := "/contacts/groups/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, contactGroupRequest{NameGroup: name}, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

// Delete deletes a contact group via DELETE /contacts/groups/{id}.
func (a *ContactGroupAPI) Delete(ctx context.Context, id string) error {
	path := "/contacts/groups/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}

// ListByContact returns the groups a contact belongs to via
// GET /contacts/groups/contact/{contact_id}.
func (a *ContactGroupAPI) ListByContact(ctx context.Context, contactID string) ([]ContactGroupSummary, error) {
	var groups []ContactGroupSummary
	path := "/contacts/groups/contact/" + url.PathEscape(contactID)
	if err := a.client.get(ctx, path, nil, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}
