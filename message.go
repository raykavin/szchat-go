package szchat

import (
	"context"
	"net/url"
)

// MessageEmailAttachment is a base64-encoded email attachment accepted by
// SendMessageRequest.Attachments.
type MessageEmailAttachment struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	MimeType string `json:"mime_type,omitempty"`
}

// SendMessageRequest is the payload for MessageAPI.Send and (embedded)
// MessageAPI.SendPlus.
type SendMessageRequest struct {
	PlatformID      string                   `json:"platform_id"`
	ChannelID       string                   `json:"channel_id"`
	Type            string                   `json:"type"`
	Message         string                   `json:"message,omitempty"`
	File            string                   `json:"file,omitempty"`
	ContactName     string                   `json:"contact_name,omitempty"`
	Agent           string                   `json:"agent,omitempty"`
	AttendanceID    string                   `json:"attendance_id,omitempty"`
	CloseSession    int                      `json:"close_session,omitempty"`
	IsHSM           bool                     `json:"is_hsm,omitempty"`
	HSMTemplateName string                   `json:"hsm_template_name,omitempty"`
	Attachments     []MessageEmailAttachment `json:"attachments,omitempty"`
}

// SendMessageContactVariables creates/updates the target contact as part of
// a send_plus call. Fields carries any additional custom-field keys the
// tenant has configured.
type SendMessageContactVariables struct {
	Name     string            `json:"name,omitempty"`
	Email    string            `json:"email,omitempty"`
	Groups   []string          `json:"groups,omitempty"`
	Channels map[string]string `json:"channels,omitempty"`
	Fields   map[string]any    `json:"-"`
}

// MarshalJSON merges Fields into the encoded object.
func (v SendMessageContactVariables) MarshalJSON() ([]byte, error) {
	type alias SendMessageContactVariables
	return marshalWithExtra(alias(v), v.Fields)
}

// SendMessagePlusRequest is the payload for MessageAPI.SendPlus.
type SendMessagePlusRequest struct {
	SendMessageRequest
	ContactVariables *SendMessageContactVariables `json:"contact_variables,omitempty"`
}

// SendMessageUser identifies the agent shown as the sender of a message.
type SendMessageUser struct {
	Name string `json:"name"`
}

// SentMessage is the message envelope returned after a successful send.
type SentMessage struct {
	MessageID string          `json:"message_id"`
	Type      string          `json:"type"`
	CreatedAt string          `json:"created_at"`
	Message   string          `json:"message"`
	User      SendMessageUser `json:"user"`
}

// SendMessageResponse is the payload returned by MessageAPI.Send and
// MessageAPI.SendPlus.
type SendMessageResponse struct {
	Messages SentMessage `json:"messages"`
	Message  string      `json:"message"`
}

// MessageAPI groups the /message endpoints.
type MessageAPI struct {
	client *Client
}

// Send sends a message via POST /message/send.
func (a *MessageAPI) Send(ctx context.Context, req SendMessageRequest) (*SendMessageResponse, error) {
	var resp SendMessageResponse
	if err := a.client.post(ctx, "/message/send", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendPlus sends a message and creates/updates the target contact via
// POST /message/send_plus.
func (a *MessageAPI) SendPlus(ctx context.Context, req SendMessagePlusRequest) (*SendMessageResponse, error) {
	var resp SendMessageResponse
	if err := a.client.post(ctx, "/message/send_plus", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Read retrieves a session's messages via POST /message/read. The SZChat
// documentation does not publish the exact request/response shape for this
// endpoint, so params/the result are passed through as free-form JSON.
func (a *MessageAPI) Read(ctx context.Context, params map[string]any) (map[string]any, error) {
	var resp map[string]any
	if err := a.client.post(ctx, "/message/read", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Pending returns the count of messages awaiting delivery via
// POST /message/pending. See the note on Read regarding the undocumented
// payload shape.
func (a *MessageAPI) Pending(ctx context.Context, params map[string]any) (map[string]any, error) {
	var resp map[string]any
	if err := a.client.post(ctx, "/message/pending", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ReplaceVars substitutes template placeholders in a message via
// POST /message/replace_vars. See the note on Read regarding the
// undocumented payload shape.
func (a *MessageAPI) ReplaceVars(ctx context.Context, params map[string]any) (map[string]any, error) {
	var resp map[string]any
	if err := a.client.post(ctx, "/message/replace_vars", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AddAnnotation adds a note to a message via POST /message/{id}/annotation.
func (a *MessageAPI) AddAnnotation(ctx context.Context, messageID string, note map[string]any) (map[string]any, error) {
	var resp map[string]any
	path := "/message/" + url.PathEscape(messageID) + "/annotation"
	if err := a.client.post(ctx, path, note, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// RemoveAnnotation removes a message's note via
// DELETE /message/{id}/annotation.
func (a *MessageAPI) RemoveAnnotation(ctx context.Context, messageID string) error {
	path := "/message/" + url.PathEscape(messageID) + "/annotation"
	return a.client.delete(ctx, path, nil)
}
