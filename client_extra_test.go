package szchat

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestPostWithKeyUsesCustomHeaderInsteadOfBearer(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("API-KEY"); got != "channel-secret" {
			t.Errorf("API-KEY header = %q, want %q", got, "channel-secret")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"status": true})
	})
	// No bearer token set: this endpoint authenticates purely via API-KEY.

	var resp map[string]bool
	err := c.postWithKey(context.Background(), "/generic/messages/send", "API-KEY", "channel-secret",
		map[string]string{"hello": "world"}, &resp)
	if err != nil {
		t.Fatalf("postWithKey: %v", err)
	}
	if !resp["status"] {
		t.Fatalf("expected status=true, got %+v", resp)
	}
}

func TestPostMultipartUploadsFile(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Fatalf("expected multipart Content-Type, got %q", r.Header.Get("Content-Type"))
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm: %v", err)
		}
		file, header, err := r.FormFile("photo")
		if err != nil {
			t.Fatalf("FormFile: %v", err)
		}
		defer file.Close()
		if header.Filename != "avatar.png" {
			t.Errorf("filename = %q, want avatar.png", header.Filename)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
	})
	c.setToken("token")

	var resp map[string]bool
	err := c.postMultipart(context.Background(), "/agents/photo", "photo", "avatar.png",
		strings.NewReader("fake-image-bytes"), nil, &resp)
	if err != nil {
		t.Fatalf("postMultipart: %v", err)
	}
	if !resp["success"] {
		t.Fatalf("expected success=true, got %+v", resp)
	}
}

func TestGetRawReturnsBodyAndContentType(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte("binary-data"))
	})

	data, contentType, err := c.getRaw(context.Background(), "/config/storage/view/storage-1")
	if err != nil {
		t.Fatalf("getRaw: %v", err)
	}
	if string(data) != "binary-data" {
		t.Errorf("data = %q, want %q", data, "binary-data")
	}
	if contentType != "image/png" {
		t.Errorf("contentType = %q, want image/png", contentType)
	}
}

func TestGetRawReturnsAPIErrorOnFailure(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
	})

	_, _, err := c.getRaw(context.Background(), "/config/storage/view/missing")
	if !IsNotFound(err) {
		t.Fatalf("expected 404 APIError, got %v", err)
	}
}
