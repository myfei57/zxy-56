package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"labvent/internal/hood"
	"labvent/internal/sash"
	"labvent/internal/store"
)

func TestLvSashAfterFlowPersist(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	airflow := hood.NewAirflowService(blobs)
	velocity := hood.NewVelocityService(blobs)
	moves := sash.NewMoveService(blobs, airflow.Persist, velocity.SetFaceThreshold)
	sashes := sash.NewSashService(blobs)
	item, err := sashes.Register("h1", 900)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "airflow"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "airflow", "h1.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := moves.Move(item.ID, 300); err == nil {
		t.Fatal("move must fail when the airflow reading cannot persist")
	}
	if blobs.Exists("sash-move", item.ID) {
		t.Fatal("sash move must not be recorded when airflow persistence failed")
	}
}
