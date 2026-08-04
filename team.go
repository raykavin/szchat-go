package szchat

import (
	"context"
	"net/url"
)

// TeamTimer configures how long an attendance can remain idle before it is
// recycled.
type TeamTimer struct {
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

// TeamAgentRef is a lightweight agent reference used in a Team's sequence.
type TeamAgentRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TeamPermissions holds the per-team agent capability flags.
type TeamPermissions struct {
	ContactEdit              bool `json:"contact_edit,omitempty"`
	FinishAttendance         bool `json:"finish_attendance,omitempty"`
	ShowNumber               bool `json:"show_number,omitempty"`
	OfflineAttendance        bool `json:"offline_attendance,omitempty"`
	AgentChat                bool `json:"agent_chat,omitempty"`
	CampaignsChat            bool `json:"campaigns_chat,omitempty"`
	SendEmojis               bool `json:"send_emojis,omitempty"`
	OptinGupshup             bool `json:"optin_gupshup,omitempty"`
	ShowPreviewMessages      bool `json:"show_preview_messages,omitempty"`
	AutoAttendance           bool `json:"auto_attendance,omitempty"`
	HideContactOnHold        bool `json:"hideContactOnHold,omitempty"`
	HideContactOnAgentScreen bool `json:"hideContactOnAgentScreen,omitempty"`
	VoiceRecordEnabled       bool `json:"voice_record_enabled,omitempty"`
}

// TeamRestriction toggles a list of allowed channels for a team.
type TeamRestriction struct {
	Enabled  bool     `json:"enabled"`
	Channels []string `json:"channels,omitempty"`
}

// TeamRestrictionForTransfer restricts which teams an attendance can be
// transferred to.
type TeamRestrictionForTransfer struct {
	Enabled bool     `json:"enabled"`
	Teams   []string `json:"teams,omitempty"`
}

// TeamGroupRestriction toggles a list of allowed contact/template groups.
type TeamGroupRestriction struct {
	Enabled bool     `json:"enabled"`
	Groups  []string `json:"groups,omitempty"`
}

// TeamEndingFlow triggers a bot flow after an attendance ends.
type TeamEndingFlow struct {
	Enabled bool   `json:"enabled"`
	FlowID  string `json:"flow_id,omitempty"`
}

// TeamCopilot configures the AI copilot integration for a team.
type TeamCopilot struct {
	RestID            StringSlice `json:"rest_id,omitempty"`
	ShowSourceContent bool        `json:"show_source_content,omitempty"`
}

// TeamAIFeature is the shared shape for the summary_attendance and
// sentiment_analysis toggles.
type TeamAIFeature struct {
	Enabled bool   `json:"enabled"`
	RestID  string `json:"rest_id,omitempty"`
}

// TeamAfterAttendance configures post-attendance AI processing.
type TeamAfterAttendance struct {
	Summary                bool   `json:"summary,omitempty"`
	SentimentAnalysis      bool   `json:"sentiment_analysis,omitempty"`
	SentimentAnalysisScore bool   `json:"sentiment_analysis_score,omitempty"`
	RestID                 string `json:"rest_id,omitempty"`
}

// TeamPendingResponseTimer configures the pending-response notification
// delay.
type TeamPendingResponseTimer struct {
	Enabled bool   `json:"enabled"`
	Limit   int    `json:"limit,omitempty"`
	Type    string `json:"type,omitempty"`
}

// TeamNotifyPendingResponse configures pending-response notifications for
// contacts and agents.
type TeamNotifyPendingResponse struct {
	Contacts struct {
		Message string                   `json:"message,omitempty"`
		Timer   TeamPendingResponseTimer `json:"timer"`
	} `json:"contacts"`
	Agents struct {
		Timer TeamPendingResponseTimer `json:"timer"`
	} `json:"agents"`
}

// TeamTimerWait redirects an attendance to another team after waiting too
// long in queue.
type TeamTimerWait struct {
	Enabled  bool   `json:"enabled"`
	Limit    int    `json:"limit,omitempty"`
	Type     string `json:"type,omitempty"`
	Redirect string `json:"redirect,omitempty"`
}

// TeamWaitPositionNotice is a single queue-position notification message.
type TeamWaitPositionNotice struct {
	Enabled bool   `json:"enabled"`
	Message string `json:"message,omitempty"`
}

// TeamNotifyWaitingPosition configures queue-position notifications.
type TeamNotifyWaitingPosition struct {
	InitialWaitPosition TeamWaitPositionNotice `json:"initialWaitPosition"`
	UpdateWaitPosition  TeamWaitPositionNotice `json:"updateWaitPosition"`
}

// Team ("equipe") groups agents under shared distribution rules,
// permissions, and automation settings. The same struct is used both to
// send create/update requests and to decode list/detail responses.
type Team struct {
	ID                        string                      `json:"_id,omitempty"`
	CanBeDisabled             bool                        `json:"can_be_disabled,omitempty"`
	Name                      string                      `json:"name"`
	History                   string                      `json:"history"`
	Timer                     TeamTimer                   `json:"timer"`
	MessageEnd                string                      `json:"messageEnd,omitempty"`
	MessageAgent              string                      `json:"messageAgent,omitempty"`
	Predefined                []string                    `json:"predefined,omitempty"`
	RuleAttendance            string                      `json:"ruleAttendance"`
	Sequence                  []TeamAgentRef              `json:"sequence,omitempty"`
	Transhipment              string                      `json:"transhipment,omitempty"`
	TranshipmentName          string                      `json:"transhipmentName,omitempty"`
	Tabulations               []string                    `json:"tabulations,omitempty"`
	Permissions               TeamPermissions             `json:"permissions"`
	OptinStart                bool                        `json:"optin_start,omitempty"`
	Agents                    []string                    `json:"agents,omitempty"`
	Tags                      []string                    `json:"tags,omitempty"`
	RestrictedTeam            *TeamRestriction            `json:"restricted_team,omitempty"`
	RestrictedTeamForTransfer *TeamRestrictionForTransfer `json:"restricted_team_for_transfer,omitempty"`
	ContactActiveByGroup      *TeamGroupRestriction       `json:"contact_active_by_group,omitempty"`
	RestrictMessageTemplates  *TeamGroupRestriction       `json:"restrict_message_templates,omitempty"`
	EndingFlow                *TeamEndingFlow             `json:"ending_flow,omitempty"`
	Copilot                   *TeamCopilot                `json:"copilot,omitempty"`
	ValidatedAgents           bool                        `json:"validatedAgents,omitempty"`
	SummaryAttendance         *TeamAIFeature              `json:"summary_attendance,omitempty"`
	SentimentAnalysis         *TeamAIFeature              `json:"sentiment_analysis,omitempty"`
	AfterAttendance           *TeamAfterAttendance        `json:"after_attendance,omitempty"`
	NotifyPendingResponse     *TeamNotifyPendingResponse  `json:"notify_pending_response,omitempty"`
	NotifyWaitingPosition     *TeamNotifyWaitingPosition  `json:"notifyWaitingPosition,omitempty"`
	TimerWait                 *TeamTimerWait              `json:"timer_wait,omitempty"`
	TransferSession           bool                        `json:"transferSession,omitempty"`
	TransferSessionToHold     []string                    `json:"transferSessionToHold,omitempty"`
	SelectAgentsType          string                      `json:"selectAgentsType,omitempty"`
	HasOnlineAgents           bool                        `json:"hasOnlineAgents,omitempty"`
	HasOfflineAttendance      bool                        `json:"hasOfflineAttendance,omitempty"`
	CreatedAt                 string                      `json:"created_at,omitempty"`
	UpdatedAt                 string                      `json:"updated_at,omitempty"`
}

// TeamSummary is the lightweight {_id, name} shape returned by
// TeamAPI.Resume and AgentAPI.MyTeams.
type TeamSummary struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// TeamUpdateResponse is the ack returned by TeamAPI.Update.
type TeamUpdateResponse struct {
	Success bool   `json:"success"`
	Date    string `json:"date"`
}

// TeamListFilter holds the query parameters accepted by TeamAPI.List.
type TeamListFilter struct {
	Limit    int
	Name     string
	Paginate string
}

func (f TeamListFilter) values() url.Values {
	q := url.Values{}
	setIntParam(q, "limit", f.Limit)
	setParam(q, "name", f.Name)
	setParam(q, "paginate", f.Paginate)
	return q
}

// TeamFilterByIDsRequest is the payload for TeamAPI.FilterByIDs. CampaignIDs
// maps to the wire field "campaings", a misspelling in the SZChat API that
// is preserved here intentionally.
type TeamFilterByIDsRequest struct {
	CampaignIDs []string `json:"campaings"`
	Paginate    bool     `json:"paginate,omitempty"`
	Limit       int      `json:"limit,omitempty"`
	Page        int      `json:"page,omitempty"`
}

// TeamAPI groups the /campaigns ("equipes"/teams) endpoints.
type TeamAPI struct {
	client *Client
}

// List returns a paginated list of teams via GET /campaigns.
func (a *TeamAPI) List(ctx context.Context, filter TeamListFilter) (*PaginatedResponse[Team], error) {
	var resp PaginatedResponse[Team]
	if err := a.client.get(ctx, "/campaigns", filter.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Resume returns a lightweight {_id, name} summary of every team via
// GET /campaigns/resume/{paginate}.
func (a *TeamAPI) Resume(ctx context.Context, paginate bool) ([]TeamSummary, error) {
	flag := "0"
	if paginate {
		flag = "1"
	}
	var summaries []TeamSummary
	path := "/campaigns/resume/" + url.PathEscape(flag)
	if err := a.client.get(ctx, path, nil, &summaries); err != nil {
		return nil, err
	}
	return summaries, nil
}

// FilterByIDs returns teams matching the given ids via
// POST /campaigns/filterByIds.
func (a *TeamAPI) FilterByIDs(ctx context.Context, req TeamFilterByIDsRequest) (*PaginatedResponse[Team], error) {
	var resp PaginatedResponse[Team]
	if err := a.client.post(ctx, "/campaigns/filterByIds", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a team via POST /campaigns.
func (a *TeamAPI) Create(ctx context.Context, team Team) (*Team, error) {
	var created Team
	if err := a.client.post(ctx, "/campaigns", team, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates a team via PUT /campaigns/{id}.
func (a *TeamAPI) Update(ctx context.Context, id string, team Team) (*TeamUpdateResponse, error) {
	var resp TeamUpdateResponse
	path := "/campaigns/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, team, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a team via DELETE /campaigns/{id}.
func (a *TeamAPI) Delete(ctx context.Context, id string) error {
	path := "/campaigns/" + url.PathEscape(id)
	return a.client.delete(ctx, path, nil)
}
