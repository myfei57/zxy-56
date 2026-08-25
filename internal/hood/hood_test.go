package hood

import (
	"testing"

	"labvent/internal/store"
)

func TestAirflowPersistAndMakeup(t *testing.T) {
	blobs := store.New(t.TempDir())
	airflow := NewAirflowService(blobs)
	if err := airflow.Persist("h1", 1800); err != nil {
		t.Fatal(err)
	}
	reading, err := airflow.Latest("h1")
	if err != nil {
		t.Fatal(err)
	}
	if reading.Flow != 1800 {
		t.Fatalf("unexpected flow: %d", reading.Flow)
	}
	makeup, err := airflow.Makeup("h1")
	if err != nil {
		t.Fatal(err)
	}
	if makeup != 1920 {
		t.Fatalf("unexpected makeup flow: %d", makeup)
	}
}

func TestPressureAndBaseline(t *testing.T) {
	blobs := store.New(t.TempDir())
	pressure := NewPressureService(blobs)
	if err := pressure.Observe("h2", -12); err != nil {
		t.Fatal(err)
	}
	healthy, err := pressure.Healthy("h2")
	if err != nil {
		t.Fatal(err)
	}
	if !healthy {
		t.Fatal("pressure must be healthy at -12")
	}
	baseline := NewBaselineService(blobs)
	if err := baseline.Set("r1", 15); err != nil {
		t.Fatal(err)
	}
	current, err := baseline.Current("r1")
	if err != nil {
		t.Fatal(err)
	}
	if current != 15 {
		t.Fatalf("unexpected baseline: %f", current)
	}
}

func TestDoorOpenClose(t *testing.T) {
	blobs := store.New(t.TempDir())
	doors := NewDoorService(blobs)
	if err := doors.Open("row1"); err != nil {
		t.Fatal(err)
	}
	open, err := doors.IsOpen("row1")
	if err != nil {
		t.Fatal(err)
	}
	if !open {
		t.Fatal("door must be open")
	}
	if err := doors.Close("row1"); err != nil {
		t.Fatal(err)
	}
	open, err = doors.IsOpen("row1")
	if err != nil {
		t.Fatal(err)
	}
	if open {
		t.Fatal("door must be closed")
	}
}

func TestVelocityThresholdVerdict(t *testing.T) {
	blobs := store.New(t.TempDir())
	velocity := NewVelocityService(blobs)
	if err := velocity.SetFaceThreshold("h3", 0.5); err != nil {
		t.Fatal(err)
	}
	if err := velocity.Verdict("h3", 0.45); err != nil {
		t.Fatal(err)
	}
	if err := velocity.Verdict("h3", 0.6); err == nil {
		t.Fatal("velocity above threshold must fail verdict")
	}
}
