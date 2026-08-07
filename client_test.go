package szchat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c, err := NewClient(server.URL, "agent@example.com", "secret",
		WithRetry(3, time.Millisecond, 5*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c, server
}

func TestClientGetDecodesSuccessResponse(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"name": "ok"})
	})
	c.setToken("valid-token")

	var out map[string]string
	if err := c.get(context.Background(), "/whatever", nil, &out); err != nil {
		t.Fatalf("get: %v", err)
	}
	if out["name"] != "ok" {
		t.Fatalf("unexpected body: %v", out)
	}
}

func TestClientRefreshesTokenOnUnauthorized(t *testing.T) {
	var calls atomic.Int32

	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v4/auth/login":
			_ = json.NewEncoder(w).Encode(LoginResponse{Token: "new-token"})
		case "/api/v4/resource":
			n := calls.Add(1)
			if n == 1 {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
				return
			}
			if r.Header.Get("Authorization") != "Bearer new-token" {
				t.Errorf("expected refreshed token on retry, got %q", r.Header.Get("Authorization"))
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})

	var out map[string]string
	if err := c.get(context.Background(), "/resource", nil, &out); err != nil {
		t.Fatalf("get: %v", err)
	}
	if out["status"] != "ok" {
		t.Fatalf("unexpected body: %v", out)
	}
	if calls.Load() != 2 {
		t.Fatalf("expected 2 calls to /resource, got %d", calls.Load())
	}
}

func TestClientRetriesOnRateLimitThenSucceeds(t *testing.T) {
	var calls atomic.Int32

	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Too Many Requests"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	c.setToken("token")

	var out map[string]string
	if err := c.get(context.Background(), "/resource", nil, &out); err != nil {
		t.Fatalf("get: %v", err)
	}
	if calls.Load() != 3 {
		t.Fatalf("expected 3 attempts, got %d", calls.Load())
	}
}

func TestClientReturnsAPIErrorAfterExhaustingRetries(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Service Unavailable"})
	})
	c.setToken("token")

	err := c.get(context.Background(), "/resource", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsRateLimited(err) && !statusIs(err, http.StatusServiceUnavailable) {
		t.Fatalf("expected a 503 APIError, got %v", err)
	}
}

func TestParseAPIErrorFormats(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		wantMessage string
		wantFields  bool
	}{
		{"plain error", `{"error":"Unauthorized"}`, "Unauthorized", false},
		{"validation errors", `{"status":false,"errors":{"email":["required"]}}`, "validation failed", true},
		{"status fail response", `{"status":"fail","response":"existing name"}`, "existing name", false},
		{"success message", `{"success":false,"message":"not found"}`, "not found", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := parseAPIError(http.StatusBadRequest, []byte(tc.body))
			if err.Message != tc.wantMessage {
				t.Errorf("Message = %q, want %q", err.Message, tc.wantMessage)
			}
			if tc.wantFields && len(err.FieldErrors) == 0 {
				t.Errorf("expected FieldErrors to be populated")
			}
		})
	}
}

func TestErrorHelpers(t *testing.T) {
	err := &APIError{StatusCode: http.StatusNotFound}
	if !IsNotFound(err) {
		t.Error("IsNotFound should be true")
	}
	if IsConflict(err) {
		t.Error("IsConflict should be false")
	}
}
