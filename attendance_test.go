package szchat

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAttendanceLifecycle(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v4/session/init":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true, "message": "Atendimento iniciado",
				"response": map[string]string{"session_id": "sess-1"},
			})
		case "/api/v4/attendances/accept":
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "Contato aceito com sucesso"})
		case "/api/v4/attendances/finish":
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "Atendimento finalizado"})
		case "/api/v4/attendances/transfer":
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "Transferência efetuada com sucesso!"})
		case "/api/v4/attendances/conference/invite":
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "Convite enviado com sucesso"})
		case "/api/v4/attendances/conference/accept":
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "Entrou em conferencia com sucesso"})
		case "/api/v4/attendances/conference/finish":
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "Conferência finalizada com sucesso!"})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})
	c.setToken("token")
	ctx := context.Background()

	initResp, err := c.AttendanceAPI.Init(ctx, AttendanceInitRequest{
		ContactID: "c1", Platform: "Whatsapp", ChannelID: "ch1", TeamID: "t1", AgentID: "a1",
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if initResp.Response.SessionID != "sess-1" {
		t.Errorf("SessionID = %q, want sess-1", initResp.Response.SessionID)
	}

	if msg, err := c.AttendanceAPI.Accept(ctx, AttendanceAcceptRequest{SessionID: "sess-1"}); err != nil || msg == "" {
		t.Fatalf("Accept: msg=%q err=%v", msg, err)
	}
	if msg, err := c.AttendanceAPI.Finish(ctx, AttendanceFinishRequest{SessionID: "sess-1"}); err != nil || msg == "" {
		t.Fatalf("Finish: msg=%q err=%v", msg, err)
	}
	if msg, err := c.AttendanceAPI.Transfer(ctx, AttendanceTransferRequest{SessionID: "sess-1", Type: "agent", AgentID: "a2"}); err != nil || msg == "" {
		t.Fatalf("Transfer: msg=%q err=%v", msg, err)
	}
	if msg, err := c.AttendanceAPI.ConferenceInvite(ctx, AttendanceConferenceInviteRequest{SessionID: "sess-1", AgentID: "a2"}); err != nil || msg == "" {
		t.Fatalf("ConferenceInvite: msg=%q err=%v", msg, err)
	}
	if msg, err := c.AttendanceAPI.ConferenceAccept(ctx, AttendanceConferenceAcceptRequest{SessionID: "sess-1", Accept: true}); err != nil || msg == "" {
		t.Fatalf("ConferenceAccept: msg=%q err=%v", msg, err)
	}
	if msg, err := c.AttendanceAPI.ConferenceFinish(ctx, AttendanceConferenceFinishRequest{SessionID: "sess-1"}); err != nil || msg == "" {
		t.Fatalf("ConferenceFinish: msg=%q err=%v", msg, err)
	}
}

func TestAttendanceFindAndHistoric(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v4/attendances":
			if r.URL.Query().Get("status") != "wait" {
				t.Errorf("expected status=wait query param, got %q", r.URL.Query().Get("status"))
			}
			_ = json.NewEncoder(w).Encode(PaginatedResponse[Attendance]{
				CurrentPage: 1,
				Data:        []Attendance{{ID: "s1", Status: "wait"}},
			})
		case r.URL.Path == "/api/v4/attendances/phase/wait":
			_ = json.NewEncoder(w).Encode([]Attendance{{ID: "s1", Phase: "wait"}})
		case r.URL.Path == "/api/v4/attendances/show":
			_ = json.NewEncoder(w).Encode(AttendanceDetail{Attendance: Attendance{ID: "s1"}})
		case r.URL.Path == "/api/v4/attendances/historic":
			_ = json.NewEncoder(w).Encode(AttendanceHistoricResponse{
				Attendances: []AttendanceHistoricEntry{{ID: "h1", Protocol: "123"}},
				Periods:     []string{"2026_08"},
			})
		case r.URL.Path == "/api/v4/attendances/historic/period":
			_ = json.NewEncoder(w).Encode(AttendanceHistoricByPeriodResponse{
				Attendances: []AttendanceHistoricByPeriodEntry{{ID: "h1"}},
			})
		case r.URL.Path == "/api/v4/attendances/historic/interval":
			_ = json.NewEncoder(w).Encode(PaginatedResponse[AttendanceHistoricIntervalEntry]{
				CurrentPage: 1,
				Data:        []AttendanceHistoricIntervalEntry{{ID: "h1"}},
			})
		case r.URL.Path == "/api/v4/attendances/historic/messages":
			_ = json.NewEncoder(w).Encode([]AttendanceHistoricMessage{{MessageID: "m1"}})
		case r.URL.Path == "/api/v4/attendances/historic/protocol":
			_ = json.NewEncoder(w).Encode(AttendanceDetail{Attendance: Attendance{ID: "s1", Protocol: "123"}})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})
	c.setToken("token")
	ctx := context.Background()

	page, err := c.AttendanceAPI.Find(ctx, AttendanceFindFilter{Status: "wait"})
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	if len(page.Data) != 1 || page.Data[0].ID != "s1" {
		t.Fatalf("unexpected Find result: %+v", page)
	}

	byPhase, err := c.AttendanceAPI.FindByPhase(ctx, "wait")
	if err != nil || len(byPhase) != 1 {
		t.Fatalf("FindByPhase: %v %+v", err, byPhase)
	}

	if _, err := c.AttendanceAPI.Show(ctx, AttendanceShowRequest{SessionID: "s1"}); err != nil {
		t.Fatalf("Show: %v", err)
	}
	if _, err := c.AttendanceAPI.Historic(ctx, "contact-1"); err != nil {
		t.Fatalf("Historic: %v", err)
	}
	if _, err := c.AttendanceAPI.HistoricByPeriod(ctx, AttendanceHistoricByPeriodRequest{ContactID: "contact-1", Period: "2026_08"}); err != nil {
		t.Fatalf("HistoricByPeriod: %v", err)
	}
	if _, err := c.AttendanceAPI.HistoricByInterval(ctx, AttendanceHistoricIntervalFilter{InitialDate: "2026-08-01", EndDate: "2026-08-10"}); err != nil {
		t.Fatalf("HistoricByInterval: %v", err)
	}
	if _, err := c.AttendanceAPI.HistoricMessages(ctx, AttendanceHistoricMessagesRequest{SessionID: "s1", CreatedAt: "2026-08-01", FinishedAt: "2026-08-02"}); err != nil {
		t.Fatalf("HistoricMessages: %v", err)
	}
	if detail, err := c.AttendanceAPI.HistoricByProtocol(ctx, "123"); err != nil || detail.Protocol != "123" {
		t.Fatalf("HistoricByProtocol: detail=%+v err=%v", detail, err)
	}
}
