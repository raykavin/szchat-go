package szchat

import (
	"encoding/json"
	"testing"
)

func TestCopilotExecuteResponseUnmarshalsFlatShape(t *testing.T) {
	var resp CopilotExecuteResponse
	body := `{"status":"success","result":"olá","sourceContent":"fonte"}`
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if resp.Result != "olá" || resp.Status != "success" || resp.SourceContent != "fonte" {
		t.Fatalf("unexpected flat decode: %+v", resp)
	}
	if len(resp.Items) != 0 {
		t.Fatalf("expected no Items for flat shape, got %v", resp.Items)
	}
}

func TestCopilotExecuteResponseUnmarshalsItemsShape(t *testing.T) {
	var resp CopilotExecuteResponse
	body := `{"1":"primeiro item","2":"segundo item"}`
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if resp.Result != "" || resp.Status != "" {
		t.Fatalf("expected empty flat fields for items shape, got %+v", resp)
	}
	if resp.Items["1"] != "primeiro item" || resp.Items["2"] != "segundo item" {
		t.Fatalf("unexpected Items decode: %v", resp.Items)
	}
}
