package szchat

import (
	"context"
	"net/url"
)

// ReportAttendanceEntry is a single row of ReportAPI.Attendances's response.
// AgentName and AgentLogin are only populated for the analytic ("a") report
// type.
type ReportAttendanceEntry struct {
	ChannelDescription   string `json:"channel_description,omitempty"`
	CampaignName         string `json:"campaign_name,omitempty"`
	Attendances          int    `json:"attendances,omitempty"`
	TotalMessages        int    `json:"total_messages,omitempty"`
	TotalSessionMessages int    `json:"total_session_messages,omitempty"`
	TotalHSM             int    `json:"total_hsm,omitempty"`
	AgentName            string `json:"agent_name,omitempty"`
	AgentLogin           string `json:"agent_login,omitempty"`
}

// ReportAttendancesResponse is the payload returned by ReportAPI.Attendances.
type ReportAttendancesResponse struct {
	Results  []ReportAttendanceEntry `json:"results"`
	Total    int                     `json:"total"`
	Page     int                     `json:"page"`
	LastPage int                     `json:"last_page"`
}

// ReportAttendancesFilter holds the query parameters accepted by
// ReportAPI.Attendances.
type ReportAttendancesFilter struct {
	InitialDate string
	EndDate     string
	Campaigns   string
	// Type selects the report shape: "a" (analytic, per-agent) or "s"
	// (synthetic, aggregated). Defaults to "a" server-side when empty.
	Type      string
	ChannelID string
	Page      int
}

func (f ReportAttendancesFilter) values() url.Values {
	q := url.Values{}
	setParam(q, "initial_date", f.InitialDate)
	setParam(q, "end_date", f.EndDate)
	setParam(q, "campaigns", f.Campaigns)
	setParam(q, "type", f.Type)
	setParam(q, "channel_id", f.ChannelID)
	setIntParam(q, "page", f.Page)
	return q
}

// ReportAPI groups the /reports endpoints. All operations require an
// administrator account.
type ReportAPI struct {
	client *Client
}

// Attendances returns the attendances report for a date range via
// GET /reports/attendances. initial_date and end_date must fall within the
// same month.
func (a *ReportAPI) Attendances(ctx context.Context, filter ReportAttendancesFilter) (*ReportAttendancesResponse, error) {
	var resp ReportAttendancesResponse
	if err := a.client.get(ctx, "/reports/attendances", filter.values(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
