package fan

import (
	"testing"

	"labvent/internal/audit"
	"labvent/internal/hood"
	"labvent/internal/store"
)

func TestFanStartStop(t *testing.T) {
	blobs := store.New(t.TempDir())
	auditService := audit.NewService(blobs)
	controller := NewController(blobs, hood.NewPressureService(blobs), hood.NewDoorService(blobs), auditService, nil)
	item, err := controller.Register("1号排风机", "row1", "h1", "z1", "end")
	if err != nil {
		t.Fatal(err)
	}
	if err := controller.Start(item.ID); err != nil {
		t.Fatal(err)
	}
	state, err := controller.StateOf(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if state != StateRunning {
		t.Fatalf("expected running, got %s", state)
	}
	if err := controller.Stop(item.ID); err != nil {
		t.Fatal(err)
	}
	state, err = controller.StateOf(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if state != StateStopped {
		t.Fatalf("expected stopped, got %s", state)
	}
}

func TestFanListByRow(t *testing.T) {
	blobs := store.New(t.TempDir())
	auditService := audit.NewService(blobs)
	controller := NewController(blobs, hood.NewPressureService(blobs), hood.NewDoorService(blobs), auditService, nil)
	if _, err := controller.Register("A", "rowA", "h1", "z1", "end"); err != nil {
		t.Fatal(err)
	}
	if _, err := controller.Register("B", "rowB", "h2", "z2", "end"); err != nil {
		t.Fatal(err)
	}
	items, err := controller.List("rowA")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Name != "A" {
		t.Fatalf("unexpected row fans: %v", items)
	}
}
