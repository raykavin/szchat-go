package szchat

import (
	"context"
	"net/url"
)

// ContactField is a tenant-defined custom contact field.
type ContactField struct {
	ID               string `json:"_id,omitempty"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Identifier       bool   `json:"identifier"`
	Required         bool   `json:"required,omitempty"`
	Confidential     bool   `json:"confidential,omitempty"`
	Encrypted        bool   `json:"encrypted,omitempty"`
	AllowSearch      bool   `json:"allow_search,omitempty"`
	DescriptionAgent bool   `json:"description_agent,omitempty"`
	Validation       string `json:"validation,omitempty"`
	Conditional      string `json:"conditional,omitempty"`
	TotalCharacters  string `json:"total_characters,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

// ContactFieldDeleteResponse is the payload returned by
// ContactFieldAPI.Delete.
type ContactFieldDeleteResponse struct {
	Success []bool `json:"success"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

// ContactFieldAPI groups the /contacts/fields endpoints, which manage the
// custom field definitions used by ContactAPI.UpdateFields.
type ContactFieldAPI struct {
	client *Client
}

// List returns all custom contact field definitions via
// GET /contacts/fields.
func (a *ContactFieldAPI) List(ctx context.Context) ([]ContactField, error) {
	var fields []ContactField
	if err := a.client.get(ctx, "/contacts/fields", nil, &fields); err != nil {
		return nil, err
	}
	return fields, nil
}

// Create creates a custom contact field via POST /contacts/fields.
func (a *ContactFieldAPI) Create(ctx context.Context, field ContactField) (*ContactField, error) {
	var created ContactField
	if err := a.client.post(ctx, "/contacts/fields", field, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a custom contact field via PUT /contacts/fields/{id}.
func (a *ContactFieldAPI) Update(ctx context.Context, id string, field ContactField) (*ContactField, error) {
	var updated ContactField
	path := "/contacts/fields/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, field, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete deletes a custom contact field via DELETE /contacts/fields/{id}.
func (a *ContactFieldAPI) Delete(ctx context.Context, id string) (*ContactFieldDeleteResponse, error) {
	var resp ContactFieldDeleteResponse
	path := "/contacts/fields/" + url.PathEscape(id)
	if err := a.client.delete(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
