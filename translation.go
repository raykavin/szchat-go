package szchat

import "context"

// TranslateRequest is the payload for TranslationAPI.Translate.
type TranslateRequest struct {
	Language  string `json:"language"`
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
	MessageID string `json:"message_id,omitempty"`
}

// TranslateResponse is the payload returned by TranslationAPI.Translate.
type TranslateResponse struct {
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
	TranslatedText         string `json:"translatedText"`
}

// DetectLanguageRequest is the payload for TranslationAPI.DetectLanguage.
type DetectLanguageRequest struct {
	Message   string `json:"message"`
	ContactID string `json:"contact_id"`
}

// DetectLanguageResponse is the payload returned by
// TranslationAPI.DetectLanguage.
type DetectLanguageResponse struct {
	Confidence float64 `json:"confidence"`
	Language   string  `json:"language"`
	Message    string  `json:"message"`
}

// ToggleAutoTranslateResponse is the payload returned by
// TranslationAPI.ToggleAutoTranslate.
type ToggleAutoTranslateResponse struct {
	Success       bool `json:"success"`
	AutoTranslate bool `json:"auto_translate"`
}

// TranslationAPI groups the agent simultaneous-translation endpoints
// (/user/agent/stt/translate*).
type TranslationAPI struct {
	client *Client
}

// Translate translates a message to the target language via
// POST /user/agent/stt/translate.
func (a *TranslationAPI) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error) {
	var resp TranslateResponse
	if err := a.client.post(ctx, "/user/agent/stt/translate", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DetectLanguage detects the source language of a message via
// POST /user/agent/stt/translate/detect.
func (a *TranslationAPI) DetectLanguage(ctx context.Context, req DetectLanguageRequest) (*DetectLanguageResponse, error) {
	var resp DetectLanguageResponse
	if err := a.client.post(ctx, "/user/agent/stt/translate/detect", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ToggleAutoTranslate toggles automatic translation for a session via
// POST /user/agent/stt/translate/activeAutoTranslate.
func (a *TranslationAPI) ToggleAutoTranslate(ctx context.Context, sessionID string) (*ToggleAutoTranslateResponse, error) {
	var resp ToggleAutoTranslateResponse
	req := struct {
		SessionID string `json:"session_id"`
	}{SessionID: sessionID}
	if err := a.client.post(ctx, "/user/agent/stt/translate/activeAutoTranslate", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
