package szchat

import "context"

// WhatsAppMessageContent is a single WhatsApp message's content, as
// forwarded via WhatsAppAPI.CreateAttendance. Fields used depend on Type
// ("TEXT", "IMAGE", "VIDEO", "AUDIO", "DOCUMENT").
type WhatsAppMessageContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	MediaURL string `json:"mediaUrl,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	Caption  string `json:"caption,omitempty"`
	FileName string `json:"fileName,omitempty"`
}

// WhatsAppMessage is a single message within a WhatsAppAPI.CreateAttendance
// request.
type WhatsAppMessage struct {
	MessageID  string                 `json:"messageId"`
	SenderType string                 `json:"senderType"`
	EventAt    int64                  `json:"eventAt"`
	Message    WhatsAppMessageContent `json:"message"`
}

// WhatsAppCreateAttendanceRequest is the payload for
// WhatsAppAPI.CreateAttendance.
type WhatsAppCreateAttendanceRequest struct {
	Messages      []WhatsAppMessage `json:"messages"`
	ChatUserID    string            `json:"chatUserId"`
	AgentHandoff  bool              `json:"agentHandoff"`
	AgentsGroupID string            `json:"agentsGroupId,omitempty"`
	ChatUserPhoto string            `json:"chatUserPhoto,omitempty"`
	ChatUserName  string            `json:"chatUserName,omitempty"`
	CustomerID    string            `json:"customerId,omitempty"`
}

// WhatsAppCreateAttendanceResponse is the payload returned by
// WhatsAppAPI.CreateAttendance.
type WhatsAppCreateAttendanceResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// WhatsAppAPI groups the WhatsApp receptive API endpoint, used by a BOT
// integration to hand a WhatsApp conversation off into Chat Center.
// Authentication is a per-channel API key (the "apiKey" header, holding the
// channel id), not the client's agent bearer token.
type WhatsAppAPI struct {
	client *Client
}

// CreateAttendance forwards a WhatsApp conversation (and optionally hands
// it off to a human agent) via POST /whatsapp/attendances. apiKey is the
// target channel's id.
func (a *WhatsAppAPI) CreateAttendance(ctx context.Context, apiKey string, req WhatsAppCreateAttendanceRequest) (*WhatsAppCreateAttendanceResponse, error) {
	var resp WhatsAppCreateAttendanceResponse
	if err := a.client.postWithKey(ctx, "/whatsapp/attendances", "apiKey", apiKey, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
