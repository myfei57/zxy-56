package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"labvent/internal/audit"
	"labvent/internal/fan"
	"labvent/internal/hood"
	"labvent/internal/store"
)

func TestLvFanFailoverOrder(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	pressure := hood.NewPressureService(blobs)
	doors := hood.NewDoorService(blobs)
	auditService := audit.NewService(blobs)
	controller := fan.NewController(blobs, pressure, doors, auditService, nil)
	primary, err := controller.Register("主机", "row1", "h1", "z1", "primary")
	if err != nil {
		t.Fatal(err)
	}
	standby, err := controller.Register("备机", "row1", "h1", "z1", "standby")
	if err != nil {
		t.Fatal(err)
	}
	if err := controller.Start(primary.ID); err != nil {
		t.Fatal(err)
	}
	block := filepath.Join(dir, "fan-run", standby.ID+".json")
	if err := os.MkdirAll(filepath.Dir(block), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(block, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := controller.Failover(primary.ID); err == nil {
		t.Fatal("failover must fail when the standby cannot start")
	}
	state, err := controller.StateOf(primary.ID)
	if err != nil {
		t.Fatal(err)
	}
	if state != fan.StateRunning {
		t.Fatalf("primary must stay running when the standby fails to start, got %s", state)
	}
}
