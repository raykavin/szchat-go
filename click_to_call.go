package szchat

import "context"

// ClickToCallResponse is the payload returned by ClickToCallAPI.Call.
type ClickToCallResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// ClickToCallAPI groups the agent click-to-call endpoint.
type ClickToCallAPI struct {
	client *Client
}

// Call places a call from the authenticated agent's extension to a contact
// via POST /user/agent/call.
func (a *ClickToCallAPI) Call(ctx context.Context, contactID string) (*ClickToCallResponse, error) {
	var resp ClickToCallResponse
	req := struct {
		ContactID string `json:"contact_id"`
	}{ContactID: contactID}
	if err := a.client.post(ctx, "/user/agent/call", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
