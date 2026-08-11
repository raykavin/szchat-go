package szchat

import (
	"context"
	"io"
)

// AgentProfileUpdateRequest is the payload for AgentAPI.UpdateProfile.
type AgentProfileUpdateRequest struct {
	Name            string `json:"name,omitempty"`
	Ramal           int    `json:"ramal,omitempty"`
	Password        string `json:"password,omitempty"`
	NewPassword     string `json:"new_password,omitempty"`
	DefaultLanguage string `json:"default_language,omitempty"`
}

// AgentProfileUpdateResponse is the payload returned by
// AgentAPI.UpdateProfile.
type AgentProfileUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AgentPhotoUploadResponse is the payload returned by AgentAPI.UploadPhoto.
type AgentPhotoUploadResponse struct {
	Success bool   `json:"success"`
	Date    string `json:"date"`
}

// UploadPhoto uploads the authenticated agent's profile photo via
// POST /agents/photo. filename should include the image extension
// (jpeg/jpg/png).
func (a *AgentAPI) UploadPhoto(ctx context.Context, filename string, photo io.Reader) (*AgentPhotoUploadResponse, error) {
	var resp AgentPhotoUploadResponse
	if err := a.client.postMultipart(ctx, "/agents/photo", "photo", filename, photo, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateProfile edits the authenticated agent's own name, extension,
// password, or default language via PUT /user/agents/update.
func (a *AgentAPI) UpdateProfile(ctx context.Context, req AgentProfileUpdateRequest) (*AgentProfileUpdateResponse, error) {
	var resp AgentProfileUpdateResponse
	if err := a.client.put(ctx, "/user/agents/update", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
