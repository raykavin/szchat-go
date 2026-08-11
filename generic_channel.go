package szchat

import "context"

// GenericChannelText is a generic-channel text message body.
type GenericChannelText struct {
	Body string `json:"body"`
}

// GenericChannelMedia is a generic-channel media message body. Caption only
// applies to image/video messages; Filename only applies to document
// messages.
type GenericChannelMedia struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// GenericChannelLocation is a generic-channel location message body.
type GenericChannelLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GenericChannelMessage is a single inbound message forwarded to Chat
// Center by a generic-channel integration. Exactly one of Text, Image,
// Document, Audio, Video, or Location should be set, matching Type.
type GenericChannelMessage struct {
	From      string                  `json:"from"`
	ID        string                  `json:"id"`
	Timestamp string                  `json:"timestamp"`
	Type      string                  `json:"type"`
	Text      *GenericChannelText     `json:"text,omitempty"`
	Image     *GenericChannelMedia    `json:"image,omitempty"`
	Document  *GenericChannelMedia    `json:"document,omitempty"`
	Audio     *GenericChannelMedia    `json:"audio,omitempty"`
	Video     *GenericChannelMedia    `json:"video,omitempty"`
	Location  *GenericChannelLocation `json:"location,omitempty"`
}

// GenericChannelContactProfile is a generic-channel contact's display
// profile.
type GenericChannelContactProfile struct {
	Name  string `json:"name"`
	Photo string `json:"photo,omitempty"`
}

// GenericChannelContact identifies the contact a GenericChannelMessage is
// from. VCard is only used when the accompanying message's Type is
// "contact".
type GenericChannelContact struct {
	Profile    GenericChannelContactProfile `json:"profile"`
	PlatformID string                       `json:"platform_id"`
	VCard      string                       `json:"vcard,omitempty"`
}

// GenericChannelSendRequest is the payload for
// GenericChannelAPI.SendMessage.
type GenericChannelSendRequest struct {
	Contacts []GenericChannelContact `json:"contacts"`
	Messages []GenericChannelMessage `json:"messages"`
}

// GenericChannelDeviceStatus is the payload for
// GenericChannelAPI.SendDeviceNotification.
type GenericChannelDeviceStatus struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// GenericChannelDeviceNotificationRequest is the payload for
// GenericChannelAPI.SendDeviceNotification.
type GenericChannelDeviceNotificationRequest struct {
	Type   string                     `json:"type"`
	Device GenericChannelDeviceStatus `json:"device"`
}

// GenericChannelSendResponse is the payload returned by
// GenericChannelAPI.SendMessage and GenericChannelAPI.SendDeviceNotification.
type GenericChannelSendResponse struct {
	Status bool `json:"status"`
}

// GenericChannelAPI groups the /generic/messages/send endpoint used by
// generic-channel integrations to forward inbound messages and device
// status changes into Chat Center. Authentication is a per-channel API key
// (the "API-KEY" header), not the client's agent bearer token.
type GenericChannelAPI struct {
	client *Client
}

// SendMessage forwards an inbound message (or set of messages) from a
// generic-channel integration via POST /generic/messages/send. apiKey is
// the channel's configured API key.
func (a *GenericChannelAPI) SendMessage(ctx context.Context, apiKey string, req GenericChannelSendRequest) (*GenericChannelSendResponse, error) {
	var resp GenericChannelSendResponse
	if err := a.client.postWithKey(ctx, "/generic/messages/send", "API-KEY", apiKey, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendDeviceNotification reports a device/connection status change (e.g.
// QR code pending, low battery, connected/disconnected) via
// POST /generic/messages/send. apiKey is the channel's configured API key.
func (a *GenericChannelAPI) SendDeviceNotification(ctx context.Context, apiKey string, req GenericChannelDeviceNotificationRequest) (*GenericChannelSendResponse, error) {
	var resp GenericChannelSendResponse
	if err := a.client.postWithKey(ctx, "/generic/messages/send", "API-KEY", apiKey, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
