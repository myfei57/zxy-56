package verifycase

import (
	"testing"
	"time"

	"labvent/internal/chem"
	"labvent/internal/store"
)

func TestLvChemExpiryGate(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	expiry := chem.NewExpiryService(blobs)
	classify := chem.NewClassifyService(blobs)
	chems := chem.NewChemService(blobs, expiry, classify)
	issue := store.NewIssueGate(expiry.EffectiveDate)
	base := time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC)
	item, err := chems.Register("溶剂", "ordinary", base.Add(24*time.Hour), 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := expiry.Revise(item.ID, base.Add(1*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := issue.Check(item.ID, base.Add(2*time.Hour)); err == nil {
		t.Fatal("expired reagent must be blocked after the effective date is revised")
	}
}
