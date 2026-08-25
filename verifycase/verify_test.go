package verifycase

import (
	"testing"

	"labvent/internal/hood"
	"labvent/internal/sash"
	"labvent/internal/store"
)

func TestLvSashFaceVelocity(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	airflow := hood.NewAirflowService(blobs)
	velocity := hood.NewVelocityService(blobs)
	moves := sash.NewMoveService(blobs, airflow.Persist, velocity.SetFaceThreshold)
	sashes := sash.NewSashService(blobs)
	hoods := hood.NewHoodService(blobs)
	hoodItem, err := hoods.Register("lab1", "room1", "H1", 1.5)
	if err != nil {
		t.Fatal(err)
	}
	item, err := sashes.Register(hoodItem.ID, 900)
	if err != nil {
		t.Fatal(err)
	}
	if err := moves.Move(item.ID, 300); err != nil {
		t.Fatal(err)
	}
	if err := velocity.Verdict(hoodItem.ID, 0.7); err != nil {
		t.Fatalf("low-sash face velocity 0.7 must stay within the raised threshold: %v", err)
	}
}
