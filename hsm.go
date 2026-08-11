package szchat

import "context"

// HSMTag is a placeholder/value pair usable within an HSM template.
type HSMTag struct {
	Placeholder string `json:"placeholder"`
	TagsValue   string `json:"tags_value"`
}

// HSMBroker describes an HSM template's broker approval state.
type HSMBroker struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

// HSMTemplate is a WhatsApp message template ("HSM").
type HSMTemplate struct {
	Name       string   `json:"name"`
	Tags       []HSMTag `json:"tags,omitempty"`
	Message    []string `json:"message"`
	Category   string   `json:"category,omitempty"`
	Broker     any      `json:"broker,omitempty"`
	Visibility string   `json:"visibility,omitempty"`
}

// HSMListRequest is the payload for HSMAPI.ListAll.
type HSMListRequest struct {
	AttendanceID string `json:"attendance_id,omitempty"`
	ChannelID    string `json:"channel_id,omitempty"`
}

// HSMAPI groups the WhatsApp message template listing endpoint.
type HSMAPI struct {
	client *Client
}

// ListAll lists the HSM templates available to the authenticated agent via
// POST /hsm/listAll.
func (a *HSMAPI) ListAll(ctx context.Context, req HSMListRequest) ([]HSMTemplate, error) {
	var templates []HSMTemplate
	if err := a.client.post(ctx, "/hsm/listAll", req, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}
