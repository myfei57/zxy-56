package verifycase

import (
	"testing"

	"labvent/internal/hood"
	"labvent/internal/lab"
	"labvent/internal/store"
)

func TestLvPressureBaselineFresh(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	baseline := hood.NewBaselineService(blobs)
	interlock := lab.NewInterlock(baseline)
	if err := baseline.Set("room1", 15); err != nil {
		t.Fatal(err)
	}
	if err := interlock.Check("room1", 16); err != nil {
		t.Fatal(err)
	}
	if err := baseline.Set("room1", 12); err != nil {
		t.Fatal(err)
	}
	if err := interlock.Check("room1", 13); err != nil {
		t.Fatalf("pressure verdict must follow the retrofit baseline: %v", err)
	}
}
