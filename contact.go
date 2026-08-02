package szchat

import (
	"context"
	"fmt"
	"net/url"
)

// ContactPlatform links a Contact to an identifier on a messaging platform.
type ContactPlatform struct {
	Platform   string `json:"platform"`
	PlatformID string `json:"platform_id"`
}

// ContactAuthor identifies who created/annotated a Contact.
type ContactAuthor struct {
	Name    string `json:"name,omitempty"`
	AgentID string `json:"agent_id,omitempty"`
	Date    struct {
		CreatedAt string `json:"created_at,omitempty"`
		UpdatedAt string `json:"updated_at,omitempty"`
	} `json:"date,omitzero"`
}

// Contact is a SZChat contact/lead record.
type Contact struct {
	ID              string            `json:"_id"`
	Name            string            `json:"name"`
	Email           string            `json:"email,omitempty"`
	DDI             string            `json:"ddi,omitempty"`
	DDD             string            `json:"ddd,omitempty"`
	Number          string            `json:"number,omitempty"`
	FinalNumber     string            `json:"final_number,omitempty"`
	Group           []string          `json:"group,omitempty"`
	Platforms       []ContactPlatform `json:"platforms,omitempty"`
	DefaultLanguage string            `json:"default_language,omitempty"`
	OptIn           string            `json:"opt_in,omitempty"`
	Observation     string            `json:"observation,omitempty"`
	Author          *ContactAuthor    `json:"author,omitempty"`
	StatusWhatsapp  int               `json:"statusWhatsapp,omitempty"`
	CreatedAt       string            `json:"created_at,omitempty"`
	UpdatedAt       string            `json:"updated_at,omitempty"`
}

// ContactRequest is the payload for creating/updating a Contact. Extra
// carries tenant-specific channel identifiers not covered by the named
// fields (e.g. a custom "Generic_MyBot" channel) and is merged into the
// top-level JSON object.
type ContactRequest struct {
	Name            string   `json:"name"`
	Email           string   `json:"email,omitempty"`
	Number          string   `json:"number,omitempty"`
	Group           []string `json:"group,omitempty"`
	DDI             string   `json:"ddi,omitempty"`
	DefaultLanguage string   `json:"default_language,omitempty"`
	FinalNumber     string   `json:"final_number,omitempty"`
	OptIn           string   `json:"opt_in,omitempty"`
	OptinGupshup    *bool    `json:"optin_gupshup,omitempty"`
	ContactToMerge  []string `json:"contactToMerge,omitempty"`

	Whatsapp         string `json:"Whatsapp,omitempty"`
	WhatsappBusiness string `json:"WhatsappBusiness,omitempty"`
	Instagram        string `json:"Instagram,omitempty"`
	InstagramDirect  string `json:"InstagramDirect,omitempty"`
	Telegram         string `json:"Telegram,omitempty"`
	Messenger        string `json:"Messenger,omitempty"`
	GoogleChat       string `json:"GoogleChat,omitempty"`
	ChatWeb          string `json:"ChatWeb,omitempty"`
	MercadoLivre     string `json:"MercadoLivre,omitempty"`
	SMS              string `json:"SMS,omitempty"`

	Extra map[string]any `json:"-"`
}

// MarshalJSON merges Extra into the encoded object.
func (r ContactRequest) MarshalJSON() ([]byte, error) {
	type alias ContactRequest
	return marshalWithExtra(alias(r), r.Extra)
}

// ContactListFilter holds the query parameters accepted by List and Search.
type ContactListFilter struct {
	ListOptions
	Name           string
	Email          string
	Platform       string
	PlatformID     string
	StartCreatedAt string
	EndCreatedAt   string
	UpdatedAt      string
}

func (f ContactListFilter) values() url.Values {
	q := f.ListOptions.values()
	setParam(q, "name", f.Name)
	setParam(q, "email", f.Email)
	setParam(q, "platform", f.Platform)
	setParam(q, "platform_id", f.PlatformID)
	setParam(q, "start_created_at", f.StartCreatedAt)
	setParam(q, "end_created_at", f.EndCreatedAt)
	setParam(q, "updated_at", f.UpdatedAt)
	return q
}

// ContactAnnotationRequest is the payload for ContactAPI.SaveAnnotation.
type ContactAnnotationRequest struct {
	ContactID   string `json:"_id"`
	Observation string `json:"observation"`
	AgentID     string `json:"agent_id"`
}

// ContactAnnotation is the stored annotation returned by SaveAnnotation.
type ContactAnnotation struct {
	Observation string         `json:"observation"`
	Author      *ContactAuthor `json:"author,omitempty"`
	CreatedAt   string         `json:"created_at,omitempty"`
	UpdatedAt   string         `json:"updated_at,omitempty"`
}

// ContactAttendanceStats summarizes a contact's attendance history on a
// given platform.
type ContactAttendanceStats struct {
	TotalMonthAttendances int    `json:"total_month_attendances"`
	WeeklyAttendances     int    `json:"weekly_attendances"`
	LastAttendanceDate    string `json:"last_attendance_date"`
}

// ContactAPI groups the /contacts endpoints.
type ContactAPI struct {
	client *Client
}

// List returns a paginated list of contacts via GET /contacts.
func (a *ContactAPI) List(ctx context.Context, filter ContactListFilter) (*PaginatedResponse[Contact], error) {
	var resp PaginatedResponse[Contact]
	if err := a.client.get(ctx, "/contacts", filter.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search returns contacts matching filter via GET /contacts/search.
func (a *ContactAPI) Search(ctx context.Context, filter ContactListFilter) (*PaginatedResponse[Contact], error) {
	var resp PaginatedResponse[Contact]
	if err := a.client.get(ctx, "/contacts/search", filter.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a single contact via GET /contacts/{id}.
func (a *ContactAPI) Get(ctx context.Context, id string) (*Contact, error) {
	var contact Contact
	path := "/contacts/" + url.PathEscape(id)
	if err := a.client.get(ctx, path, nil, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// Create creates a contact via POST /contacts.
func (a *ContactAPI) Create(ctx context.Context, req ContactRequest) (*Contact, error) {
	var contact Contact
	if err := a.client.post(ctx, "/contacts", req, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// Update updates a contact via PUT /contacts/{id}.
func (a *ContactAPI) Update(ctx context.Context, id string, req ContactRequest) (*Contact, error) {
	var contact Contact
	path := "/contacts/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, req, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// UpdateFields updates a contact's custom fields via
// PUT /contacts/update_fields/{id}.
func (a *ContactAPI) UpdateFields(ctx context.Context, id string, fields map[string]any) (map[string]any, error) {
	var resp map[string]any
	path := "/contacts/update_fields/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, fields, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Delete deletes a contact via DELETE /contacts/{id}.
func (a *ContactAPI) Delete(ctx context.Context, id string) error {
	path := "/contacts/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}

// SaveAnnotation creates/updates a contact's annotation via
// POST /contacts/annotation.
func (a *ContactAPI) SaveAnnotation(ctx context.Context, req ContactAnnotationRequest) (*ContactAnnotation, error) {
	var resp ContactAnnotation
	if err := a.client.post(ctx, "/contacts/annotation", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AttendanceStats returns a contact's attendance stats on a platform via
// GET /contacts/{id}/attendances.
func (a *ContactAPI) AttendanceStats(ctx context.Context, contactID, platform string) (*ContactAttendanceStats, error) {
	var resp ContactAttendanceStats
	path := fmt.Sprintf("/contacts/%s/attendances", url.PathEscape(contactID))
	q := url.Values{"platform": []string{platform}}
	if err := a.client.get(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
