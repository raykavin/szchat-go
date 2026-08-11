package szchat

import (
	"context"
	"net/url"
)

// AttendanceInitRequest is the payload for AttendanceAPI.Init.
type AttendanceInitRequest struct {
	ContactID          string   `json:"contact_id"`
	Platform           string   `json:"platform"`
	ChannelID          string   `json:"channel_id"`
	TeamID             string   `json:"team_id"`
	AgentID            string   `json:"agent_id"`
	HSMID              string   `json:"hsm_id,omitempty"`
	TagMessage         string   `json:"tagMessage,omitempty"`
	TagType            string   `json:"tagType,omitempty"`
	PlaceholdersParams []string `json:"placeholders_params,omitempty"`
}

// AttendanceInitResponse is the payload returned by AttendanceAPI.Init.
type AttendanceInitResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Response struct {
		SessionID string `json:"session_id"`
	} `json:"response"`
}

// AttendanceAcceptRequest is the payload for AttendanceAPI.Accept.
type AttendanceAcceptRequest struct {
	SessionID    string `json:"session_id"`
	Agent        string `json:"agent,omitempty"`
	AttendanceID string `json:"attendance_id,omitempty"`
}

// AttendanceFinishRequest is the payload for AttendanceAPI.Finish.
type AttendanceFinishRequest struct {
	SessionID    string `json:"session_id"`
	TabulationID string `json:"tabulation_id,omitempty"`
}

// AttendanceConferenceInviteRequest is the payload for
// AttendanceAPI.ConferenceInvite.
type AttendanceConferenceInviteRequest struct {
	SessionID string `json:"session_id"`
	AgentID   string `json:"agent_id"`
}

// AttendanceConferenceAcceptRequest is the payload for
// AttendanceAPI.ConferenceAccept.
type AttendanceConferenceAcceptRequest struct {
	SessionID string `json:"session_id"`
	Accept    bool   `json:"accept"`
}

// AttendanceConferenceFinishRequest is the payload for
// AttendanceAPI.ConferenceFinish.
type AttendanceConferenceFinishRequest struct {
	SessionID string `json:"session_id"`
}

// AttendanceTransferRequest is the payload for AttendanceAPI.Transfer.
type AttendanceTransferRequest struct {
	SessionID    string `json:"session_id"`
	Type         string `json:"type"`
	AttendanceID string `json:"attendance_id,omitempty"`
	AgentID      string `json:"agent_id,omitempty"`
	TransferWait bool   `json:"transfer_wait,omitempty"`
}

// AttendanceShowRequest is the payload for AttendanceAPI.Show.
type AttendanceShowRequest struct {
	SessionID string `json:"session_id,omitempty"`
	ContactID string `json:"contact_id,omitempty"`
}

// AttendanceMessage is a message reply.
type AttendanceMessage struct {
	Message   string `json:"message"`
	AgentTo   string `json:"agent_to,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// AttendanceEvent is a lifecycle event recorded on an attendance session
// (e.g. "waitStart", "humanStart").
type AttendanceEvent struct {
	Event      string `json:"event"`
	CreatedAt  string `json:"created_at,omitempty"`
	AgentFrom  string `json:"agent_from,omitempty"`
	AgentTo    string `json:"agent_to,omitempty"`
	CampaignTo string `json:"campaign_to,omitempty"`
}

// AttendanceTag is a tag applied to an attendance session.
type AttendanceTag struct {
	Question  string `json:"question,omitempty"`
	Answer    string `json:"answer,omitempty"`
	Tag       string `json:"tag,omitempty"`
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// AttendanceTalk is a message exchanged within an attendance session, as
// returned by history/detail endpoints.
type AttendanceTalk struct {
	MessageID  string `json:"message_id,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
	Origin     string `json:"origin,omitempty"`
	Type       string `json:"type,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	Blocked    bool   `json:"blocked,omitempty"`
	Message    string `json:"message,omitempty"`
	MessageRef string `json:"message_ref,omitempty"`
}

// Attendance is a session/attendance record, as returned by AttendanceAPI's
// Find, FindByPhase, and Show operations. Fields not present on a given
// endpoint's response are simply left zero-valued.
type Attendance struct {
	ID                    string            `json:"_id,omitempty"`
	Name                  string            `json:"name,omitempty"`
	ContactID             string            `json:"contact_id,omitempty"`
	ChannelID             string            `json:"channel_id,omitempty"`
	AgentID               string            `json:"agent_id,omitempty"`
	CampaignID            string            `json:"campaign_id,omitempty"`
	PlatformID            string            `json:"platform_id,omitempty"`
	Platform              string            `json:"platform,omitempty"`
	Status                string            `json:"status,omitempty"`
	Protocol              string            `json:"protocol,omitempty"`
	Phase                 string            `json:"phase,omitempty"`
	Group                 string            `json:"group,omitempty"`
	Position              int               `json:"position,omitempty"`
	LastInteraction       string            `json:"lastInteraction,omitempty"`
	IsAttendance          bool              `json:"isAttendance,omitempty"`
	ContactAlreadyStarted bool              `json:"contactAlreadyStarted,omitempty"`
	CreatedAt             string            `json:"created_at,omitempty"`
	TimerOnWait           string            `json:"timerOnWait,omitempty"`
	TimerAccept           string            `json:"timerAccept,omitempty"`
	WaitingByCodename     bool              `json:"waitingByCodename,omitempty"`
	Wait                  string            `json:"wait,omitempty"`
	ContinueFlow          bool              `json:"continueFlow,omitempty"`
	CampaignFinishMessage string            `json:"campaignFinishMessage,omitempty"`
	MenuAttempts          int               `json:"menuAttempts,omitempty"`
	AgentsConference      []string          `json:"agents_conference,omitempty"`
	Tags                  []AttendanceTag   `json:"tags,omitempty"`
	Events                []AttendanceEvent `json:"events,omitempty"`
	Talks                 []AttendanceTalk  `json:"talks,omitempty"`
	FinishedAt            string            `json:"finished_at,omitempty"`
}

// AttendanceHistoricEntry is a single item in AttendanceAPI.Historic's
// response.
type AttendanceHistoricEntry struct {
	ID         string `json:"_id"`
	Protocol   string `json:"protocol,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	FinishedAt string `json:"finished_at,omitempty"`
}

// AttendanceHistoricResponse is the payload returned by AttendanceAPI.Historic.
type AttendanceHistoricResponse struct {
	Attendances []AttendanceHistoricEntry `json:"attendances"`
	Periods     []string                  `json:"periods"`
}

// AttendanceHistoricByPeriodRequest is the payload for
// AttendanceAPI.HistoricByPeriod.
type AttendanceHistoricByPeriodRequest struct {
	ContactID        string `json:"contact_id"`
	Period           string `json:"period"`
	Platform         string `json:"platform,omitempty"`
	StatusAttendance string `json:"statusAttendance,omitempty"`
}

// AttendanceHistoricByPeriodEntry is a single item in
// AttendanceAPI.HistoricByPeriod's response.
type AttendanceHistoricByPeriodEntry struct {
	ID               string   `json:"_id"`
	SessionID        string   `json:"session_id,omitempty"`
	PlatformID       string   `json:"platform_id,omitempty"`
	Protocol         string   `json:"protocol,omitempty"`
	ChannelPlatform  string   `json:"channel_platform,omitempty"`
	ChannelID        string   `json:"channel_id,omitempty"`
	TME              float64  `json:"tme,omitempty"`
	CreatedAt        string   `json:"createdAt,omitempty"`
	AgentName        []string `json:"agentName,omitempty"`
	CampaignName     []string `json:"campaign_name,omitempty"`
	FinishedAt       string   `json:"finishedAt,omitempty"`
	StatusAttendance string   `json:"statusAttendance,omitempty"`
}

// AttendanceHistoricByPeriodResponse is the payload returned by
// AttendanceAPI.HistoricByPeriod.
type AttendanceHistoricByPeriodResponse struct {
	Attendances []AttendanceHistoricByPeriodEntry `json:"attendances"`
}

// AttendanceHistoricIntervalFilter holds the query parameters accepted by
// AttendanceAPI.HistoricByInterval.
type AttendanceHistoricIntervalFilter struct {
	InitialDate      string
	EndDate          string
	ContactID        string
	PlatformID       string
	Email            string
	StatusAttendance string
	PerPage          int
}

func (f AttendanceHistoricIntervalFilter) values() url.Values {
	q := url.Values{}
	setParam(q, "initial_date", f.InitialDate)
	setParam(q, "end_date", f.EndDate)
	setParam(q, "contact_id", f.ContactID)
	setParam(q, "platform_id", f.PlatformID)
	setParam(q, "email", f.Email)
	setParam(q, "statusAttendance", f.StatusAttendance)
	setIntParam(q, "per_page", f.PerPage)
	return q
}

// AttendanceHistoricIntervalEntry is a single item in
// AttendanceAPI.HistoricByInterval's response.
type AttendanceHistoricIntervalEntry struct {
	ID               string `json:"_id"`
	Protocol         string `json:"protocol,omitempty"`
	PlatformID       string `json:"platform_id,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	FinishedAt       string `json:"finished_at,omitempty"`
	Email            string `json:"email,omitempty"`
	StatusAttendance string `json:"statusAttendance,omitempty"`
}

// AttendanceHistoricMessagesRequest is the payload for
// AttendanceAPI.HistoricMessages.
type AttendanceHistoricMessagesRequest struct {
	SessionID  string `json:"session_id"`
	CreatedAt  string `json:"created_at"`
	FinishedAt string `json:"finished_at"`
}

// AttendanceHistoricMessage is a single item in
// AttendanceAPI.HistoricMessages's response.
type AttendanceHistoricMessage struct {
	MessageID  string `json:"message_id"`
	RequestID  string `json:"request_id,omitempty"`
	Origin     string `json:"origin,omitempty"`
	Type       string `json:"type,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	Blocked    bool   `json:"blocked,omitempty"`
	Message    string `json:"message,omitempty"`
	Privacy    string `json:"privacy,omitempty"`
	MessageRef string `json:"message_ref,omitempty"`
}

// AttendanceHistoricProtocolRequest is the payload for
// AttendanceAPI.HistoricByProtocol.
type AttendanceHistoricProtocolRequest struct {
	ProtocolID string `json:"protocol_id"`
}

// AttendanceDetail is the full session record returned by
// AttendanceAPI.HistoricByProtocol and AttendanceAPI.Show, including the
// denormalized contact/campaign/channel snapshots the API embeds.
type AttendanceDetail struct {
	Attendance
	Contact         *Contact `json:"contact,omitempty"`
	Campaign        *Team    `json:"campaign,omitempty"`
	Channel         *Channel `json:"channel,omitempty"`
	Tabulation      any      `json:"tabulation,omitempty"`
	AttendanceTimer struct {
		TTA  float64 `json:"tta,omitempty"`
		TMA  float64 `json:"tma,omitempty"`
		TME  float64 `json:"tme,omitempty"`
		TTA2 float64 `json:"tta2,omitempty"`
	} `json:"attendance_timer,omitempty"`
	ReportAt string `json:"report_at,omitempty"`
	Lite     bool   `json:"lite,omitempty"`
}

// AttendanceAPI groups the /attendances, /session/init and
// /attendances/conference endpoints covering the full attendance/session
// lifecycle: creation, acceptance, transfer, conferencing, search, and
// history.
type AttendanceAPI struct {
	client *Client
}

// Init starts a new attendance session via POST /session/init.
func (a *AttendanceAPI) Init(ctx context.Context, req AttendanceInitRequest) (*AttendanceInitResponse, error) {
	var resp AttendanceInitResponse
	if err := a.client.post(ctx, "/session/init", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Accept accepts a waiting session into attendance via
// POST /attendances/accept.
func (a *AttendanceAPI) Accept(ctx context.Context, req AttendanceAcceptRequest) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	if err := a.client.post(ctx, "/attendances/accept", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}

// Finish closes an attendance session via POST /attendances/finish.
func (a *AttendanceAPI) Finish(ctx context.Context, req AttendanceFinishRequest) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	if err := a.client.post(ctx, "/attendances/finish", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}

// Transfer transfers a session to another team or agent via
// POST /attendances/transfer.
func (a *AttendanceAPI) Transfer(ctx context.Context, req AttendanceTransferRequest) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	if err := a.client.post(ctx, "/attendances/transfer", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}

// ConferenceInvite invites an online agent to join a session's conference
// via POST /attendances/conference/invite.
func (a *AttendanceAPI) ConferenceInvite(ctx context.Context, req AttendanceConferenceInviteRequest) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	if err := a.client.post(ctx, "/attendances/conference/invite", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}

// ConferenceAccept accepts or declines a conference invite via
// POST /attendances/conference/accept.
func (a *AttendanceAPI) ConferenceAccept(ctx context.Context, req AttendanceConferenceAcceptRequest) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	if err := a.client.post(ctx, "/attendances/conference/accept", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}

// ConferenceFinish ends the authenticated agent's participation in a
// session's conference via POST /attendances/conference/finish.
func (a *AttendanceAPI) ConferenceFinish(ctx context.Context, req AttendanceConferenceFinishRequest) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	if err := a.client.post(ctx, "/attendances/conference/finish", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}

// AttendanceFindFilter holds the query parameters accepted by
// AttendanceAPI.Find.
type AttendanceFindFilter struct {
	Name       string
	ContactID  string
	ChannelID  string
	AgentID    string
	CampaignID string
	PlatformID string
	Platform   string
	Status     string
	Protocol   string
}

func (f AttendanceFindFilter) values() url.Values {
	q := url.Values{}
	setParam(q, "name", f.Name)
	setParam(q, "contact_id", f.ContactID)
	setParam(q, "channel_id", f.ChannelID)
	setParam(q, "agent_id", f.AgentID)
	setParam(q, "campaign_id", f.CampaignID)
	setParam(q, "platform_id", f.PlatformID)
	setParam(q, "platform", f.Platform)
	setParam(q, "status", f.Status)
	setParam(q, "protocol", f.Protocol)
	return q
}

// Find searches attendances via GET /attendances.
func (a *AttendanceAPI) Find(ctx context.Context, filter AttendanceFindFilter) (*PaginatedResponse[Attendance], error) {
	var resp PaginatedResponse[Attendance]
	if err := a.client.get(ctx, "/attendances", filter.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FindByPhase lists attendances in a given phase ("auto", "wait", or
// "human") via GET /attendances/phase/{phase}.
func (a *AttendanceAPI) FindByPhase(ctx context.Context, phase string) ([]Attendance, error) {
	var attendances []Attendance
	path := "/attendances/phase/" + url.PathEscape(phase)
	if err := a.client.get(ctx, path, nil, &attendances); err != nil {
		return nil, err
	}
	return attendances, nil
}

// Show returns a single session by session_id or contact_id via
// POST /attendances/show.
func (a *AttendanceAPI) Show(ctx context.Context, req AttendanceShowRequest) (*AttendanceDetail, error) {
	var detail AttendanceDetail
	if err := a.client.post(ctx, "/attendances/show", req, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// Historic returns a contact's most recent attendances and the list of
// periods with history available via POST /attendances/historic.
func (a *AttendanceAPI) Historic(ctx context.Context, contactID string) (*AttendanceHistoricResponse, error) {
	var resp AttendanceHistoricResponse
	req := struct {
		ContactID string `json:"contact_id"`
	}{ContactID: contactID}
	if err := a.client.post(ctx, "/attendances/historic", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// HistoricByPeriod returns a contact's attendances within a given period
// (format "YYYY_MM" or "YYYY_MM_DD") via POST /attendances/historic/period.
func (a *AttendanceAPI) HistoricByPeriod(ctx context.Context, req AttendanceHistoricByPeriodRequest) (*AttendanceHistoricByPeriodResponse, error) {
	var resp AttendanceHistoricByPeriodResponse
	if err := a.client.post(ctx, "/attendances/historic/period", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// HistoricByInterval returns a paginated list of attendances within a date
// interval via GET /attendances/historic/interval.
func (a *AttendanceAPI) HistoricByInterval(ctx context.Context, filter AttendanceHistoricIntervalFilter) (*PaginatedResponse[AttendanceHistoricIntervalEntry], error) {
	var resp PaginatedResponse[AttendanceHistoricIntervalEntry]
	if err := a.client.get(ctx, "/attendances/historic/interval", filter.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// HistoricMessages returns the messages exchanged in a finished session via
// POST /attendances/historic/messages.
func (a *AttendanceAPI) HistoricMessages(ctx context.Context, req AttendanceHistoricMessagesRequest) ([]AttendanceHistoricMessage, error) {
	var messages []AttendanceHistoricMessage
	if err := a.client.post(ctx, "/attendances/historic/messages", req, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

// HistoricByProtocol returns the full session record for a given protocol
// number via POST /attendances/historic/protocol.
func (a *AttendanceAPI) HistoricByProtocol(ctx context.Context, protocolID string) (*AttendanceDetail, error) {
	var detail AttendanceDetail
	req := AttendanceHistoricProtocolRequest{ProtocolID: protocolID}
	if err := a.client.post(ctx, "/attendances/historic/protocol", req, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}
