package sash

import (
	"testing"

	"labvent/internal/hood"
	"labvent/internal/store"
)

// orderingBlob wraps the file store and records the sequence in which blobs are
// saved, so a test can assert that linked airflow is persisted before the
// sash-move displacement record becomes visible.
type orderingBlob struct {
	store.Blob
	saves []string
}

func (b *orderingBlob) Save(kind string, id string, payload any) error {
	b.saves = append(b.saves, kind)
	return b.Blob.Save(kind, id, payload)
}

// TestMovePersistsAirflowBeforeDisplacement is a regression guard for the
// morning hood-3 incident: when the sash is raised, the linked airflow reading
// (used to compute makeup air and hold cabinet negative pressure) must be
// persisted BEFORE the sash-move displacement record. Otherwise a reader
// reacting to the displacement still sees the stale airflow value and computes
// makeup air from the old reading, briefly driving the hood pressure back to
// positive and letting reagent odors escape at the sash opening.
func TestMovePersistsAirflowBeforeDisplacement(t *testing.T) {
	blobs := &orderingBlob{Blob: store.New(t.TempDir())}
	const hoodID = "hood3"

	airflow := hood.NewAirflowService(blobs)
	// Pre-existing airflow reading derived from the old sash height; makeup air
	// computed from this value alone holds the cabinet negative pressure.
	if err := airflow.Persist(hoodID, 1200); err != nil {
		t.Fatal(err)
	}
	// Reset the observed save order so only the move itself is captured below.
	blobs.saves = nil

	sashes := NewSashService(blobs)
	sashItem, err := sashes.Register(hoodID, 50)
	if err != nil {
		t.Fatal(err)
	}
	// Lower the sash from full-open (50) to 30: the new airflow must reflect
	// the new height, not the stale 1200 reading.
	moves := NewMoveService(blobs, airflow.Persist, func(string, float64) error { return nil })
	if err := moves.Move(sashItem.ID, 30); err != nil {
		t.Fatal(err)
	}

	// The sash-move record must not appear before the airflow reading.
	airflowIdx, moveIdx := indexOf(blobs.saves, "airflow"), indexOf(blobs.saves, "sash-move")
	if airflowIdx == -1 {
		t.Fatalf("airflow was never persisted; saves=%v", blobs.saves)
	}
	if moveIdx == -1 {
		t.Fatalf("sash-move was never persisted; saves=%v", blobs.saves)
	}
	if airflowIdx > moveIdx {
		t.Fatalf("airflow (idx %d) must be persisted before sash-move (idx %d); saves=%v", airflowIdx, moveIdx, blobs.saves)
	}

	// The persisted reading must be the new flow (30*40=1200) and, more
	// importantly, must already be on disk by the time the displacement is
	// recorded — so makeup air is computed from the current reading, not stale.
	latest, err := airflow.Latest(hoodID)
	if err != nil {
		t.Fatal(err)
	}
	if latest.Flow != 1200 {
		t.Fatalf("unexpected airflow after move: %d", latest.Flow)
	}
	makeup, err := airflow.Makeup(hoodID)
	if err != nil {
		t.Fatal(err)
	}
	if makeup != 1320 {
		t.Fatalf("unexpected makeup flow: %d", makeup)
	}
}

func indexOf(items []string, want string) int {
	for i, v := range items {
		if v == want {
			return i
		}
	}
	return -1
}
