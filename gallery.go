package szchat

import (
	"context"
	"net/url"
)

// GalleryMedia is a single media item as returned by GalleryAPI.Medias.
type GalleryMedia struct {
	MessageID string `json:"message_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Origin    string `json:"origin,omitempty"`
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Blocked   bool   `json:"blocked,omitempty"`
	StorageID string `json:"storage_id,omitempty"`
	Filename  string `json:"filename,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`
	Legend    string `json:"legend,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	User      struct {
		Name  string `json:"name,omitempty"`
		Photo string `json:"photo,omitempty"`
	} `json:"user,omitempty"`
	SessionID        string `json:"session_id,omitempty"`
	SessionCreatedAt string `json:"session_created_at,omitempty"`
	FinishedAt       string `json:"finished_at,omitempty"`
	Status           string `json:"status,omitempty"`
	AgentName        string `json:"agent_name,omitempty"`
	Protocol         string `json:"protocol,omitempty"`
}

// GalleryMediasResponse is the payload returned by GalleryAPI.Medias.
type GalleryMediasResponse struct {
	CurrentSession   []GalleryMedia                  `json:"current_session"`
	HistoricSessions PaginatedResponse[GalleryMedia] `json:"historic_sessions"`
}

// GalleryAPI groups the contact media gallery endpoint.
type GalleryAPI struct {
	client *Client
}

// Medias returns the media exchanged with a contact, across the current and
// historic sessions, via GET /agent/historic/medias. mediaType filters by
// kind (e.g. "images", "videos", "sounds", "files"); pass "" for all types.
// limit caps historic_sessions entries (1-50, defaults to 5 server-side).
func (a *GalleryAPI) Medias(ctx context.Context, contactID, mediaType string, limit int) (*GalleryMediasResponse, error) {
	q := url.Values{"contact_id": []string{contactID}}
	setParam(q, "type", mediaType)
	setIntParam(q, "limit", limit)

	var resp GalleryMediasResponse
	if err := a.client.get(ctx, "/agent/historic/medias", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
