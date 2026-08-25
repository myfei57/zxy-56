package hood

import (
	"testing"

	"labvent/internal/store"
)

func TestSetFaceThresholdPersistsArgument(t *testing.T) {
	blobs := store.New(t.TempDir())
	velocity := NewVelocityService(blobs)

	// A threshold above the full-open value must be persisted verbatim, not
	// silently replaced with the constant full-open threshold.
	if err := velocity.SetFaceThreshold("h4", 0.8); err != nil {
		t.Fatal(err)
	}
	got, err := velocity.Threshold("h4")
	if err != nil {
		t.Fatal(err)
	}
	if got != 0.8 {
		t.Fatalf("threshold not persisted verbatim: got %.4f, want 0.8", got)
	}
}

func TestSetFaceThresholdRejectsNonPositive(t *testing.T) {
	blobs := store.New(t.TempDir())
	velocity := NewVelocityService(blobs)
	if err := velocity.SetFaceThreshold("h5", 0); err == nil {
		t.Fatal("zero threshold must be rejected")
	}
	if err := velocity.SetFaceThreshold("h5", -0.5); err == nil {
		t.Fatal("negative threshold must be rejected")
	}
}

func TestVelocityVerdictUsesPersistedThreshold(t *testing.T) {
	blobs := store.New(t.TempDir())
	velocity := NewVelocityService(blobs)

	// Persist a low-sash threshold of 0.8 (the whole point of the fix: the
	// limit rises when the sash is lowered). A measured 0.6 must then pass,
	// whereas it would trip the old hardcoded 0.5 limit.
	if err := velocity.SetFaceThreshold("h6", 0.8); err != nil {
		t.Fatal(err)
	}
	if err := velocity.Verdict("h6", 0.6); err != nil {
		t.Fatalf("measured 0.6 must pass under raised threshold 0.8: %v", err)
	}
	if err := velocity.Verdict("h6", 0.85); err == nil {
		t.Fatal("measured 0.85 must trip threshold 0.8")
	}
}
