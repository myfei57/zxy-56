package verifycase

import (
	"testing"

	"labvent/internal/audit"
	"labvent/internal/fan"
	"labvent/internal/hood"
	"labvent/internal/store"
)

func TestLvFanVibrationLatchReset(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	pressure := hood.NewPressureService(blobs)
	doors := hood.NewDoorService(blobs)
	auditService := audit.NewService(blobs)
	controller := fan.NewController(blobs, pressure, doors, auditService, nil)
	auditService.SetResetCallback(controller.ClearVibration)
	item, err := controller.Register("排风机", "row1", "h1", "z1", "primary")
	if err != nil {
		t.Fatal(err)
	}
	if err := controller.TripVibration(item.ID); err != nil {
		t.Fatal(err)
	}
	if err := auditService.Reset("fan", item.ID, "vibration reset"); err != nil {
		t.Fatal(err)
	}
	latched, err := controller.VibrationLatched(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if latched {
		t.Fatal("fan vibration latch must release after the reset signal")
	}
}
