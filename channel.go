package szchat

import "context"

// Channel is a configured messaging channel (a WhatsApp number, a Telegram
// bot, a web chat widget, etc).
type Channel struct {
	ID             string `json:"_id"`
	Platform       string `json:"platform"`
	FlowID         string `json:"flow_id,omitempty"`
	Receptive      bool   `json:"receptive,omitempty"`
	Description    string `json:"description,omitempty"`
	Number         string `json:"number,omitempty"`
	BotToken       string `json:"bot_token,omitempty"`
	MessengerToken string `json:"messenger_token,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// ChannelPlatform describes a messaging platform SZChat can create channels
// for.
type ChannelPlatform struct {
	Platform   string `json:"platform"`
	Name       string `json:"name"`
	Color      string `json:"color,omitempty"`
	PlatformID string `json:"platform_id,omitempty"`
	SaveRedis  bool   `json:"save_redis,omitempty"`
	Icon       string `json:"icon,omitempty"`
}

// ChannelAPI groups the /channels endpoints.
type ChannelAPI struct {
	client *Client
}

// List returns all configured channels via GET /channels.
func (a *ChannelAPI) List(ctx context.Context) ([]Channel, error) {
	var channels []Channel
	if err := a.client.get(ctx, "/channels", nil, &channels); err != nil {
		return nil, err
	}
	return channels, nil
}

// ListPlatforms returns all messaging platforms available for new channels
// via GET /channels/platforms.
func (a *ChannelAPI) ListPlatforms(ctx context.Context) (map[string]ChannelPlatform, error) {
	platforms := map[string]ChannelPlatform{}
	if err := a.client.get(ctx, "/channels/platforms", nil, &platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}

// ListActivePlatforms returns only the platforms with at least one active
// channel via GET /channels/platforms/active.
func (a *ChannelAPI) ListActivePlatforms(ctx context.Context) (map[string]ChannelPlatform, error) {
	platforms := map[string]ChannelPlatform{}
	if err := a.client.get(ctx, "/channels/platforms/active", nil, &platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}
