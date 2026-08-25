package fan

import (
	"fmt"
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

func TestFanFailoverMakeBeforeBreak(t *testing.T) {
	blobs := store.New(t.TempDir())
	auditService := audit.NewService(blobs)
	controller := NewController(blobs, hood.NewPressureService(blobs), hood.NewDoorService(blobs), auditService, nil)
	primary, err := controller.Register("1号排风机", "row1", "h1", "z1", "end")
	if err != nil {
		t.Fatal(err)
	}
	standby, err := controller.Register("2号备机", "row1", "h1", "z1", "standby")
	if err != nil {
		t.Fatal(err)
	}
	if err := controller.Start(primary.ID); err != nil {
		t.Fatal(err)
	}
	if err := controller.Failover(primary.ID); err != nil {
		t.Fatal(err)
	}
	// Make-before-break: standby must be running before the primary is stopped,
	// so at no point are both fans down and cabinet negative pressure lost.
	primaryNow, err := controller.Get(primary.ID)
	if err != nil {
		t.Fatal(err)
	}
	standbyNow, err := controller.Get(standby.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !standbyNow.Running() {
		t.Fatalf("standby %s should be running after failover, state=%s", standby.ID, standbyNow.State)
	}
	if primaryNow.State != StateStopped {
		t.Fatalf("primary %s should be stopped after failover, state=%s", primary.ID, primaryNow.State)
	}
}

// failingStandbyStore wraps a store.Blob and rejects any save that would mark
// the standby fan as running, simulating a hardware no-start so that Start()
// cannot confirm the standby up. Everything else passes through to the real
// store.
type failingStandbyStore struct {
	store.Blob
	standbyID string
}

func (f failingStandbyStore) Save(kind string, id string, payload any) error {
	if kind == "fan" && id == f.standbyID {
		return errStandbyNoStart
	}
	return f.Blob.Save(kind, id, payload)
}

var errStandbyNoStart = fmt.Errorf("standby hardware did not start")

func TestFanFailoverKeepsPrimaryWhenStandbyFails(t *testing.T) {
	blobs := store.New(t.TempDir())
	auditService := audit.NewService(blobs)
	controller := NewController(blobs, hood.NewPressureService(blobs), hood.NewDoorService(blobs), auditService, nil)
	primary, err := controller.Register("1号排风机", "row1", "h1", "z1", "end")
	if err != nil {
		t.Fatal(err)
	}
	standby, err := controller.Register("2号备机", "row1", "h1", "z1", "standby")
	if err != nil {
		t.Fatal(err)
	}
	if err := controller.Start(primary.ID); err != nil {
		t.Fatal(err)
	}

	// Switch to a store that refuses to persist the standby as running, so the
	// standby never reaches a running state and Failover must abort.
	controller.blobs = failingStandbyStore{Blob: blobs, standbyID: standby.ID}

	err = controller.Failover(primary.ID)
	if err == nil {
		t.Fatal("expected failover to abort when standby fails to start")
	}
	// Primary must still be running: existing airflow preserved, no
	// break-before-make gap that would drop cabinet negative pressure.
	primaryNow, err := controller.Get(primary.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !primaryNow.Running() {
		t.Fatalf("primary must remain running when standby fails, state=%s", primaryNow.State)
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
