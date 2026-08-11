package szchat

import (
	"context"
	"net/url"
	"strings"
)

// AgentPermission holds per-agent capability overrides.
type AgentPermission struct {
	HideReviewStars               bool `json:"hideReviewStars,omitempty"`
	HandleSimultaneousAttendances bool `json:"handleSimultaneousAttendances,omitempty"`
	SimultaneousAttendancesLimit  int  `json:"simultaneousAttendancesLimit,omitempty"`
}

// Agent is a SZChat attendant account.
type Agent struct {
	ID                        string           `json:"_id"`
	Type                      string           `json:"type,omitempty"`
	Name                      string           `json:"name"`
	Email                     string           `json:"email"`
	Codename                  string           `json:"codename,omitempty"`
	Ramal                     string           `json:"ramal,omitempty"`
	Begin                     string           `json:"begin,omitempty"`
	End                       string           `json:"end,omitempty"`
	Campaigns                 []string         `json:"campaigns,omitempty"`
	Groups                    []string         `json:"groups,omitempty"`
	Photo                     string           `json:"photo,omitempty"`
	History                   string           `json:"history,omitempty"`
	UsernameCallcenter        string           `json:"username_callcenter,omitempty"`
	EmailForgotPassword       string           `json:"email_forgot_password,omitempty"`
	EnableLoginWithRemoteAuth bool             `json:"enable_login_with_remote_auth,omitempty"`
	AgentPermission           *AgentPermission `json:"agentPermission,omitempty"`
	Status                    string           `json:"status,omitempty"`
	LoggedAt                  string           `json:"logged_at,omitempty"`
	CreatedAt                 string           `json:"created_at,omitempty"`
	UpdatedAt                 string           `json:"updated_at,omitempty"`
}

// AgentRequest is the payload for creating/updating an Agent.
type AgentRequest struct {
	Name                      string           `json:"name"`
	Password                  string           `json:"password,omitempty"`
	Email                     string           `json:"email"`
	Codename                  string           `json:"codename,omitempty"`
	Begin                     string           `json:"begin,omitempty"`
	End                       string           `json:"end,omitempty"`
	Campaigns                 []string         `json:"campaigns,omitempty"`
	Photo                     string           `json:"photo,omitempty"`
	Groups                    []string         `json:"groups,omitempty"`
	History                   string           `json:"history,omitempty"`
	Ramal                     string           `json:"ramal,omitempty"`
	UsernameCallcenter        string           `json:"username_callcenter,omitempty"`
	PasswordCallcenter        string           `json:"password_callcenter,omitempty"`
	EmailForgotPassword       string           `json:"email_forgot_password,omitempty"`
	EnableLoginWithRemoteAuth *bool            `json:"enable_login_with_remote_auth,omitempty"`
	AgentPermission           *AgentPermission `json:"agentPermission,omitempty"`
}

// AgentListFilter holds the query parameters accepted by AgentAPI.List.
type AgentListFilter struct {
	ListOptions
	Name       string
	CampaignID string
}

func (f AgentListFilter) values() url.Values {
	q := f.ListOptions.values()
	setParam(q, "name", f.Name)
	setParam(q, "campaign_id", f.CampaignID)
	return q
}

// AgentOnline is an agent/admin entry as returned by the online-status
// endpoints.
type AgentOnline struct {
	ID                 string        `json:"_id"`
	Type               string        `json:"type,omitempty"`
	Email              string        `json:"email,omitempty"`
	Name               string        `json:"name"`
	Codename           string        `json:"codename,omitempty"`
	Ramal              string        `json:"ramal,omitempty"`
	Begin              string        `json:"begin,omitempty"`
	End                string        `json:"end,omitempty"`
	Status             string        `json:"status,omitempty"`
	LoggedAt           string        `json:"logged_at,omitempty"`
	Campaigns          []string      `json:"campaigns,omitempty"`
	CampaignsOnline    []TeamSummary `json:"campaigns_online,omitempty"`
	Groups             []string      `json:"groups,omitempty"`
	Attendances        int           `json:"attendances,omitempty"`
	Phone              string        `json:"phone,omitempty"`
	UserSince          string        `json:"user_since,omitempty"`
	RecentContacts     []string      `json:"recent_contacts,omitempty"`
	Photo              string        `json:"photo,omitempty"`
	Language           string        `json:"language,omitempty"`
	History            string        `json:"history,omitempty"`
	Pause              any           `json:"pause,omitempty"`
	MsTeams            any           `json:"msTeams,omitempty"`
	SessionToken       string        `json:"session_token,omitempty"`
	UsernameCallcenter string        `json:"username_callcenter,omitempty"`
	PasswordCallcenter string        `json:"password_callcenter,omitempty"`
	CreatedAt          string        `json:"created_at,omitempty"`
	UpdatedAt          string        `json:"updated_at,omitempty"`
}

// AgentTeam is a team/campaign entry as returned by AgentAPI.MyTeams,
// richer than TeamSummary since /user/agents/campaigns embeds most of the
// team's own configuration. Timer is typed any because this endpoint
// returns its days/hours/minutes as strings, unlike Team.Timer's ints.
type AgentTeam struct {
	ID             string `json:"_id"`
	Name           string `json:"name"`
	History        string `json:"history,omitempty"`
	Timer          any    `json:"timer,omitempty"`
	Transhipment   string `json:"transhipment,omitempty"`
	MessageEnd     string `json:"messageEnd,omitempty"`
	MessageAgent   string `json:"messageAgent,omitempty"`
	RuleAttendance string `json:"ruleAttendance,omitempty"`
	Tabulations    string `json:"tabulations,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// AgentTeamsResponse is the payload returned by AgentAPI.MyTeams.
type AgentTeamsResponse struct {
	Success bool `json:"success"`
	Agent   struct {
		ID        string   `json:"_id"`
		Type      string   `json:"type,omitempty"`
		Email     string   `json:"email,omitempty"`
		Name      string   `json:"name,omitempty"`
		Codename  string   `json:"codename,omitempty"`
		Begin     string   `json:"begin,omitempty"`
		End       string   `json:"end,omitempty"`
		Campaigns []string `json:"campaigns"`
		Groups    []string `json:"groups,omitempty"`
		CreatedAt string   `json:"created_at,omitempty"`
		UpdatedAt string   `json:"updated_at,omitempty"`
	} `json:"agent"`
	Campaigns []AgentTeam `json:"campaigns"`
}

// AgentAttendanceContact is a contact currently or awaiting attendance by
// the authenticated agent.
type AgentAttendanceContact struct {
	ID             string   `json:"_id"`
	Status         string   `json:"status,omitempty"`
	Phase          string   `json:"phase,omitempty"`
	Email          string   `json:"email,omitempty"`
	Phone          string   `json:"phone,omitempty"`
	UserSince      string   `json:"user_since,omitempty"`
	RecentContacts []string `json:"recent_contacts,omitempty"`
}

// AgentAttendances is the payload returned by AgentAPI.MyAttendances and
// MyAttendancesPlus.
type AgentAttendances struct {
	Attendance []AgentAttendanceContact `json:"attendance"`
	Wait       []AgentAttendanceContact `json:"wait"`
}

// AgentGrades is the payload returned by AgentAPI.MyGrades.
type AgentGrades struct {
	Average struct {
		Average float64 `json:"average"`
		Message string  `json:"message"`
	} `json:"average"`
}

// AgentAPI groups the /agents and current-agent (/user/agents/*) endpoints.
type AgentAPI struct {
	client *Client
}

// List returns a paginated list of agents via GET /agents.
func (a *AgentAPI) List(ctx context.Context, filter AgentListFilter) (*PaginatedResponse[Agent], error) {
	var resp PaginatedResponse[Agent]
	if err := a.client.get(ctx, "/agents", filter.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates an agent via POST /agents.
func (a *AgentAPI) Create(ctx context.Context, req AgentRequest) (*Agent, error) {
	var agent Agent
	if err := a.client.post(ctx, "/agents", req, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// Get returns a single agent via GET /agents/{id}.
func (a *AgentAPI) Get(ctx context.Context, id string) (*Agent, error) {
	var agent Agent
	path := "/agents/" + url.PathEscape(id)
	if err := a.client.get(ctx, path, nil, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetByEmail returns a single agent via GET /agents/email/{email}.
func (a *AgentAPI) GetByEmail(ctx context.Context, email string) (*Agent, error) {
	var agent Agent
	path := "/agents/email/" + url.PathEscape(email)
	if err := a.client.get(ctx, path, nil, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// Update updates an agent via PUT /agents/{id}.
func (a *AgentAPI) Update(ctx context.Context, id string, req AgentRequest) (*Agent, error) {
	var agent Agent
	path := "/agents/" + url.PathEscape(id)
	if err := a.client.put(ctx, path, req, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// Delete deletes one or more agents via DELETE /agents/{id}, accepting
// several ids for a bulk delete.
func (a *AgentAPI) Delete(ctx context.Context, ids ...string) error {
	path := "/agents/" + url.PathEscape(strings.Join(ids, ","))
	return a.client.delete(ctx, path, nil)
}

// ListOnlineStatus lists agents' online status via GET /online-status,
// optionally filtered by platform and/or status. Pass an empty string to
// leave a filter unset.
func (a *AgentAPI) ListOnlineStatus(ctx context.Context, platform, status string) ([]AgentOnline, error) {
	q := url.Values{}
	setParam(q, "platform", platform)
	setParam(q, "status", status)

	var agents []AgentOnline
	if err := a.client.get(ctx, "/online-status", q, &agents); err != nil {
		return nil, err
	}
	return agents, nil
}

// ListOnlineAgents lists online agents (and admins, if allUsers is true) via
// GET /user/agents/online.
func (a *AgentAPI) ListOnlineAgents(ctx context.Context, allUsers bool, paginate string) (*PaginatedResponse[AgentOnline], error) {
	q := url.Values{}
	if allUsers {
		q.Set("allUsers", "true")
	}
	setParam(q, "paginate", paginate)

	var resp PaginatedResponse[AgentOnline]
	if err := a.client.get(ctx, "/user/agents/online", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MyTeams returns the authenticated agent's teams via
// GET /user/agents/campaigns.
func (a *AgentAPI) MyTeams(ctx context.Context) (*AgentTeamsResponse, error) {
	var resp AgentTeamsResponse
	if err := a.client.get(ctx, "/user/agents/campaigns", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MyAttendances returns the authenticated agent's active/waiting contacts
// via GET /user/agents/attendances.
func (a *AgentAPI) MyAttendances(ctx context.Context) (*AgentAttendances, error) {
	var resp AgentAttendances
	if err := a.client.get(ctx, "/user/agents/attendances", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MyAttendancesPlus is like MyAttendances but with extended contact details
// via GET /user/agents/attendances_plus.
func (a *AgentAPI) MyAttendancesPlus(ctx context.Context) (*AgentAttendances, error) {
	var resp AgentAttendances
	if err := a.client.get(ctx, "/user/agents/attendances_plus", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MyGrades returns the authenticated agent's attendance grades via
// GET /user/agents/grades.
func (a *AgentAPI) MyGrades(ctx context.Context) (*AgentGrades, error) {
	var resp AgentGrades
	if err := a.client.get(ctx, "/user/agents/grades", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ToggleTeam toggles the authenticated agent's active team/campaign via
// POST /user/agents/toggle/campaign.
func (a *AgentAPI) ToggleTeam(ctx context.Context, teamID string) (string, error) {
	var resp struct {
		Message string `json:"message"`
	}
	req := struct {
		CampaignID string `json:"campaign_id"`
	}{CampaignID: teamID}
	if err := a.client.post(ctx, "/user/agents/toggle/campaign", req, &resp); err != nil {
		return "", err
	}
	return resp.Message, nil
}

// UpdateLastInteraction refreshes the authenticated agent's last-interaction
// timestamp via POST /user/agents/updateLastInteraction.
func (a *AgentAPI) UpdateLastInteraction(ctx context.Context) error {
	return a.client.post(ctx, "/user/agents/updateLastInteraction", nil, nil)
}

// AgentSessionListFilter holds the query parameters accepted by
// AgentAPI.SessionAttendances and AgentAPI.SessionWaits.
type AgentSessionListFilter struct {
	Name        string
	Paginate    bool
	Page        int
	ContactPlus bool
}

func (f AgentSessionListFilter) values() url.Values {
	q := url.Values{}
	setParam(q, "name", f.Name)
	if f.Paginate {
		q.Set("paginate", "true")
	}
	setIntParam(q, "page", f.Page)
	if f.ContactPlus {
		q.Set("contact_plus", "true")
	}
	return q
}

// SessionAttendances returns a paginated list of the authenticated agent's
// in-progress attendance sessions via GET /user/agents/sessions/attendances.
func (a *AgentAPI) SessionAttendances(ctx context.Context, filter AgentSessionListFilter) (*PaginatedResponse[Attendance], error) {
	var resp PaginatedResponse[Attendance]
	if err := a.client.get(ctx, "/user/agents/sessions/attendances", filter.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SessionWaits returns a paginated list of the authenticated agent's waiting
// attendance sessions via GET /user/agents/sessions/waits.
func (a *AgentAPI) SessionWaits(ctx context.Context, filter AgentSessionListFilter) (*PaginatedResponse[Attendance], error) {
	var resp PaginatedResponse[Attendance]
	if err := a.client.get(ctx, "/user/agents/sessions/waits", filter.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
