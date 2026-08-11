package szchat

import (
	"context"
	"encoding/json"
	"net/url"
)

// CopilotAssistant describes an available copilot assistant.
type CopilotAssistant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CopilotExecuteRequest is the payload for CopilotAPI.Execute.
type CopilotExecuteRequest struct {
	RestID            string `json:"restId"`
	Content           string `json:"content"`
	ShowSourceContent bool   `json:"showSourceContent,omitempty"`
}

// CopilotExecuteResponse is the payload returned by CopilotAPI.Execute. The
// SZChat API documents two response shapes for this endpoint (a flat
// status/result/sourceContent object, or a map of numbered result items),
// so both are represented here; callers should check which one is
// populated (Result != "" vs. len(Items) > 0).
type CopilotExecuteResponse struct {
	Status        string            `json:"status,omitempty"`
	Result        string            `json:"result,omitempty"`
	SourceContent string            `json:"sourceContent,omitempty"`
	Items         map[string]string `json:"-"`
}

// UnmarshalJSON decodes either of CopilotAPI.Execute's two documented
// response shapes into CopilotExecuteResponse.
func (r *CopilotExecuteResponse) UnmarshalJSON(data []byte) error {
	type flatShape struct {
		Status        string `json:"status,omitempty"`
		Result        string `json:"result,omitempty"`
		SourceContent string `json:"sourceContent,omitempty"`
	}
	var flat flatShape
	if err := json.Unmarshal(data, &flat); err == nil && (flat.Status != "" || flat.Result != "") {
		r.Status = flat.Status
		r.Result = flat.Result
		r.SourceContent = flat.SourceContent
		return nil
	}

	var items map[string]string
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	r.Items = items
	return nil
}

// CopilotAPI groups the agent copilot endpoints.
type CopilotAPI struct {
	client *Client
}

// List returns the copilot assistants available to the authenticated agent
// via GET /user/agent/copilot/list.
func (a *CopilotAPI) List(ctx context.Context, search string) ([]CopilotAssistant, error) {
	q := url.Values{}
	setParam(q, "search", search)
	var assistants []CopilotAssistant
	if err := a.client.get(ctx, "/user/agent/copilot/list", q, &assistants); err != nil {
		return nil, err
	}
	return assistants, nil
}

// Execute runs a copilot assistant against content via
// POST /user/agent/copilot/execute.
func (a *CopilotAPI) Execute(ctx context.Context, req CopilotExecuteRequest) (*CopilotExecuteResponse, error) {
	var resp CopilotExecuteResponse
	if err := a.client.post(ctx, "/user/agent/copilot/execute", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
