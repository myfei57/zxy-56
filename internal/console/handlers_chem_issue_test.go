package console

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestIssueRejectsExpiredReagentAfterRevision reproduces the reported incident:
// a reagent whose expiry was revised into the past must be rejected on issue.
// Before the fix the in-memory expiry cache shadowed the revised record, so
// IssueGate read the stale (future) date and let the expired solvent through.
func TestIssueRejectsExpiredReagentAfterRevision(t *testing.T) {
	server := NewServer(buildDeps(t))

	// Register a reagent with a future expiry so registration and an initial
	// issue both succeed.
	registerBody := bytes.NewBufferString(`{"name":"丙酮","class":"ordinary","expiry_at":"2099-01-01T00:00:00Z","quantity":5}`)
	req := httptest.NewRequest(http.MethodPost, "/chemicals", registerBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("register status: %d body=%s", rec.Code, rec.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	// Revise the expiry into the past.
	past := time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)
	reviseBody := bytes.NewBufferString(`{"effective_until":"` + past + `"}`)
	req = httptest.NewRequest(http.MethodPost, "/chemicals/"+created.ID+"/revise", reviseBody)
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("revise status: %d body=%s", rec.Code, rec.Body.String())
	}

	// Issuing the now-expired reagent must be rejected.
	issueBody := bytes.NewBufferString(`{"quantity":1}`)
	req = httptest.NewRequest(http.MethodPost, "/chemicals/"+created.ID+"/issue", issueBody)
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("issue of expired reagent must be rejected (409), got %d body=%s", rec.Code, rec.Body.String())
	}
}
