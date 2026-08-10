package szchat

import (
	"encoding/json"
	"testing"
)

func TestContactUnmarshalJSONNormalizesNumericDDIAndDDD(t *testing.T) {
	cases := []struct {
		name    string
		body    string
		wantDDI string
		wantDDD string
	}{
		{
			name:    "string ddi and ddd",
			body:    `{"_id":"1","name":"a","ddi":"55","ddd":"11"}`,
			wantDDI: "55",
			wantDDD: "11",
		},
		{
			name:    "numeric ddi and ddd",
			body:    `{"_id":"2","name":"b","ddi":55,"ddd":11}`,
			wantDDI: "55",
			wantDDD: "11",
		},
		{
			name:    "mixed numeric ddi, string ddd",
			body:    `{"_id":"3","name":"c","ddi":351,"ddd":"21"}`,
			wantDDI: "351",
			wantDDD: "21",
		},
		{
			name:    "missing ddi and ddd",
			body:    `{"_id":"4","name":"d"}`,
			wantDDI: "",
			wantDDD: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var c Contact
			if err := json.Unmarshal([]byte(tc.body), &c); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			if c.DDI != tc.wantDDI {
				t.Errorf("DDI = %q, want %q", c.DDI, tc.wantDDI)
			}
			if c.DDD != tc.wantDDD {
				t.Errorf("DDD = %q, want %q", c.DDD, tc.wantDDD)
			}
		})
	}
}

func TestContactUnmarshalJSONStillCollectsExtraFields(t *testing.T) {
	var c Contact
	body := `{"_id":"1","name":"a","ddi":55,"Generic_MyBot":"custom-id"}`
	if err := json.Unmarshal([]byte(body), &c); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if c.DDI != "55" {
		t.Errorf("DDI = %q, want %q", c.DDI, "55")
	}
	if got, ok := c.Extra["Generic_MyBot"]; !ok || got != "custom-id" {
		t.Errorf("Extra[Generic_MyBot] = %v, ok=%v, want %q", got, ok, "custom-id")
	}
}

func TestContactRoundTripsThroughPaginatedResponse(t *testing.T) {
	body := `{"current_page":1,"per_page":50,"total":2,"data":[{"_id":"1","name":"a","ddi":55},{"_id":"2","name":"b","ddi":"351"}]}`

	var resp PaginatedResponse[Contact]
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("Unmarshal PaginatedResponse[Contact]: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 contacts, got %d", len(resp.Data))
	}
	if resp.Data[0].DDI != "55" {
		t.Errorf("Data[0].DDI = %q, want %q", resp.Data[0].DDI, "55")
	}
	if resp.Data[1].DDI != "351" {
		t.Errorf("Data[1].DDI = %q, want %q", resp.Data[1].DDI, "351")
	}
}
