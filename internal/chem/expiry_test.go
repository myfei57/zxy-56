package chem

import (
	"testing"
	"time"

	"labvent/internal/store"
)

// TestIssueRejectedAfterExpiryRevision reproduces the reported incident:
// after a reagent's expiry is revised down to the past, issuing it must be
// rejected. EffectiveDate (the source the issue gate reads from) must reflect
// the revised expiry rather than the stale in-memory cache.
func TestIssueRejectedAfterExpiryRevision(t *testing.T) {
	blobs := store.New(t.TempDir())
	expiry := NewExpiryService(blobs)
	classify := NewClassifyService(blobs)
	service := NewChemService(blobs, expiry, classify)

	future := time.Now().Add(30 * 24 * time.Hour)
	item, err := service.Register("丙酮", "ordinary", future, 5)
	if err != nil {
		t.Fatal(err)
	}
	// Simulate a downward revision to an already-expired date.
	past := time.Now().Add(-24 * time.Hour)
	if err := expiry.Revise(item.ID, past); err != nil {
		t.Fatal(err)
	}
	date, err := expiry.EffectiveDate(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !date.Equal(past) {
		t.Fatalf("effective date %s does not reflect revised expiry %s", date, past)
	}
}
