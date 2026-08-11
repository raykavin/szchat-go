package szchat

// This file models the "API Receptiva" outbound webhook payloads: Chat
// Center POSTs these to a host configured in Integrações > API whenever a
// session event occurs. There is no endpoint to call for these — the SDK's
// role is limited to giving callers typed shapes to decode the payloads
// their own HTTP server receives. Use WebhookEnvelope.Event (via a first
// decode into WebhookEnvelope) to pick which concrete *Data type to decode
// data into next.

// WebhookInfo identifies which webhook configuration delivered an event.
type WebhookInfo struct {
	App   string `json:"app"`
	Host  string `json:"host"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

// WebhookEnvelope is the outer shape of every receptive-API webhook
// delivery. Decode into this first, then decode the raw Data into the
// concrete type matching Webhook.Key (see the Webhook* Data types below).
type WebhookEnvelope struct {
	Data    map[string]any `json:"data"`
	Webhook WebhookInfo    `json:"webhook"`
}

// WebhookMessageData is the "data.content" shape for the "client_message"
// webhook (event "message").
type WebhookMessageData struct {
	MessageID  string `json:"message_id"`
	Type       string `json:"type"`
	StorageID  string `json:"storage_id,omitempty"`
	Legend     string `json:"legend,omitempty"`
	PlatformID string `json:"platform_id"`
	Message    string `json:"message,omitempty"`
	Origin     string `json:"origin"`
	CreatedAt  string `json:"created_at"`
}

// WebhookAgentMessageData is the "data.content" shape for the
// "agent_message" webhook (event "messageAgent").
type WebhookAgentMessageData struct {
	Message   string `json:"message,omitempty"`
	AgentFrom string `json:"agent_from"`
	AgentTo   string `json:"agent_to"`
	Type      string `json:"type"`
	Legend    string `json:"legend,omitempty"`
	StorageID string `json:"storage_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

// WebhookSessionData is the "data" shape shared by the session lifecycle
// webhooks: "enter_queue" (waitStart), "accept_attendance" (humanStart),
// "attendance_finish" (humanFinish), and "attendance_transfer"
// (humanTransferAgent). Not every field is populated by every event; see
// the SZChat "API Receptiva" documentation for the exact set per event.
type WebhookSessionData struct {
	ID                        string `json:"_id"`
	Name                      string `json:"name"`
	CampaignID                string `json:"campaign_id"`
	ContactID                 string `json:"contact_id"`
	ChannelID                 string `json:"channel_id"`
	Platform                  string `json:"platform"`
	PlatformID                string `json:"platform_id"`
	Status                    string `json:"status"`
	IsAttendance              string `json:"isAttendance"`
	CreatedAt                 string `json:"created_at"`
	LastInteraction           string `json:"lastInteraction,omitempty"`
	CountMessagesNotification int    `json:"count_messages_notification,omitempty"`
	Protocol                  string `json:"protocol,omitempty"`
	AgentID                   string `json:"agent_id,omitempty"`
	TabulationID              string `json:"tabulation_id,omitempty"`
	TabulationName            string `json:"tabulation_name,omitempty"`
}

// WebhookConferenceData is the "data" shape shared by the conference
// webhooks: "conference_invite" (humanInviteConfer), "conference_accept"
// (humanStartConfer), and "conference_finish" (humanFinishConfer).
type WebhookConferenceData struct {
	WebhookSessionData
	AgentInvited        string   `json:"agent_invited,omitempty"`
	AgentsConference    []string `json:"agents_conference,omitempty"`
	AgentAcceptedConfer string   `json:"agent_accepted_confer,omitempty"`
	AgentExitConfer     string   `json:"agent_exit_confer,omitempty"`
}

// WebhookPause is the pause snapshot embedded in the "agent_pause" and
// "agent_resume" webhooks.
type WebhookPause struct {
	Pause
	StartedAt string `json:"started_at,omitempty"`
}

// WebhookAgentPauseData is the "data" shape for the "agent_pause" webhook
// (event "adminPauseAgent") and "agent_resume" webhook (event
// "agentResume").
type WebhookAgentPauseData struct {
	Pause  WebhookPause `json:"pause"`
	UserID string       `json:"user_id"`
}

// WebhookSessionAgent is the agent snapshot embedded in the
// "agent_login"/"agent_logoff" webhooks.
type WebhookSessionAgent struct {
	ID        string `json:"_id"`
	Codename  string `json:"codename,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Status    string `json:"status,omitempty"`
	OriginWeb bool   `json:"origin_web,omitempty"`
	OriginApp bool   `json:"origin_app,omitempty"`
	IP        string `json:"ip,omitempty"`
}

// WebhookAgentSessionData is the "data" shape for the "agent_login" webhook
// (event "agentSignIn") and "agent_logoff" webhook (event "agentSignOut").
type WebhookAgentSessionData struct {
	Session struct {
		ID        string              `json:"_id,omitempty"`
		Agent     WebhookSessionAgent `json:"agent"`
		SessionID string              `json:"session_id,omitempty"`
	} `json:"session"`
	SessionID string `json:"session_id"`
}

// GenericChannelOutboundMessage is the payload Chat Center POSTs to a
// generic channel's configured host to deliver an outbound message. Like
// the receptive-API webhooks above, this is not something the SDK calls —
// it documents the shape a generic-channel integration's own HTTP server
// receives, for callers that need to decode it.
type GenericChannelOutboundMessage struct {
	To        string                  `json:"to"`
	SessionID string                  `json:"session_id"`
	ContactID string                  `json:"contact_id"`
	ChannelID string                  `json:"channel_id"`
	Type      string                  `json:"type"`
	Text      *GenericChannelText     `json:"text,omitempty"`
	Image     *GenericChannelMedia    `json:"image,omitempty"`
	Document  *GenericChannelMedia    `json:"document,omitempty"`
	Audio     *GenericChannelMedia    `json:"audio,omitempty"`
	Video     *GenericChannelMedia    `json:"video,omitempty"`
	Location  *GenericChannelLocation `json:"location,omitempty"`
	Contact   *struct {
		VCard string `json:"vcard"`
	} `json:"contact,omitempty"`
}
