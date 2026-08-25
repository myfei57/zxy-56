package verifycase

import (
	"testing"
	"time"

	"labvent/internal/chem"
	"labvent/internal/store"
)

func TestLvChemReclassifyFresh(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	expiry := chem.NewExpiryService(blobs)
	classify := chem.NewClassifyService(blobs)
	chems := chem.NewChemService(blobs, expiry, classify)
	cabinets := store.NewCabinetService(classify.ClassOf)
	item, err := chems.Register("管制试剂", "ordinary", time.Now().Add(72*time.Hour), 5)
	if err != nil {
		t.Fatal(err)
	}
	rules := []chem.Rule{{From: "ordinary", To: "controlled", MinQuantity: 1}}
	if err := classify.SetRules(rules); err != nil {
		t.Fatal(err)
	}
	if err := classify.Refresh([]string{item.ID}); err != nil {
		t.Fatal(err)
	}
	if err := cabinets.Verify(item.ID, "ordinary"); err == nil {
		t.Fatal("controlled reagent must not stay in an ordinary cabinet after rule refresh")
	}
}
