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

func TestLvBatchExhaustOrder(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	pressure := hood.NewPressureService(blobs)
	doors := hood.NewDoorService(blobs)
	auditService := audit.NewService(blobs)
	controller := fan.NewController(blobs, pressure, doors, auditService, nil)
	endFan, err := controller.Register("末端风机", "row1", "h1", "z1", "end")
	if err != nil {
		t.Fatal(err)
	}
	if err := controller.IssueStart(endFan.ID); err != nil {
		t.Fatal(err)
	}
	block := filepath.Join(dir, "fan-run", endFan.ID+".json")
	if err := os.MkdirAll(filepath.Dir(block), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(block, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := controller.StartRow("row1"); err == nil {
		t.Fatal("row start must fail when an end fan cannot confirm running")
	}
	open, err := doors.IsOpen("row1")
	if err == nil && open {
		t.Fatal("hood doors must not open before the end fans are confirmed running")
	}
}
