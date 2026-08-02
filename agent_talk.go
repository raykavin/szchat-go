package szchat

import (
	"context"
	"net/url"
)

// AgentTalkMessageRequest is the payload for AgentTalkAPI.SendMessage.
type AgentTalkMessageRequest struct {
	AgentTo string `json:"agent_to"`
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
	File    string `json:"file,omitempty"`
	Legend  string `json:"legend,omitempty"`
}

// AgentTalkMessage is an internal agent-to-agent chat message. The API is
// inconsistent about the mime-type field name across endpoints, so both
// MimeType ("mime_type", used on send) and Mimetype ("mimetype", used on
// read) are populated depending on which endpoint returned it.
type AgentTalkMessage struct {
	ID        string `json:"_id,omitempty"`
	Message   string `json:"message,omitempty"`
	AgentFrom string `json:"agent_from,omitempty"`
	AgentTo   string `json:"agent_to,omitempty"`
	Type      string `json:"type,omitempty"`
	Filename  string `json:"filename,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`
	Mimetype  string `json:"mimetype,omitempty"`
	Legend    string `json:"legend,omitempty"`
	StorageID string `json:"storage_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// AgentTalkSendResponse is the payload returned by AgentTalkAPI.SendMessage.
type AgentTalkSendResponse struct {
	Event   string           `json:"event"`
	Content AgentTalkMessage `json:"content"`
}

// AgentConversation groups the messages exchanged with a single peer agent.
type AgentConversation struct {
	AgentID string             `json:"agent_id"`
	Talks   []AgentTalkMessage `json:"talks"`
}

// AgentTalkAPI groups the internal agent-to-agent messaging endpoints
// (/user/agents/messages*).
type AgentTalkAPI struct {
	client *Client
}

// SendMessage sends an internal message to another agent via
// POST /user/agents/messages/send.
func (a *AgentTalkAPI) SendMessage(ctx context.Context, req AgentTalkMessageRequest) (*AgentTalkSendResponse, error) {
	var resp AgentTalkSendResponse
	if err := a.client.post(ctx, "/user/agents/messages/send", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListConversations lists the authenticated agent's internal conversations
// via GET /user/agents/messages.
func (a *AgentTalkAPI) ListConversations(ctx context.Context) ([]AgentConversation, error) {
	var conversations []AgentConversation
	if err := a.client.get(ctx, "/user/agents/messages", nil, &conversations); err != nil {
		return nil, err
	}
	return conversations, nil
}

// SearchConversation returns the paginated message history with a peer
// agent via GET /user/agents/messages/read/{agent_id}.
func (a *AgentTalkAPI) SearchConversation(ctx context.Context, agentID string, opts ListOptions) (*PaginatedResponse[AgentTalkMessage], error) {
	var resp PaginatedResponse[AgentTalkMessage]
	path := "/user/agents/messages/read/" + url.PathEscape(agentID)
	if err := a.client.get(ctx, path, opts.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
