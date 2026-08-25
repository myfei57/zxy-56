package verifycase

import (
	"testing"

	"labvent/internal/audit"
	"labvent/internal/gas"
	"labvent/internal/lab"
	"labvent/internal/store"
)

func TestLvGasAlarmLatchClear(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	auditService := audit.NewService(blobs)
	gasService := gas.NewGasService(blobs)
	latch := gas.NewLatch(blobs)
	mapper := lab.NewMappingService(blobs)
	router := gas.NewZoneRouter(mapper)
	auditService.SetConfirmCallback(latch.MarkConfirmed)
	alarms := gas.NewAlarmService(gasService, latch, auditService, router)
	sensor, err := gasService.Register("A 区传感器", "room1", 10)
	if err != nil {
		t.Fatal(err)
	}
	if err := alarms.Sample(sensor.ID, 12); err != nil {
		t.Fatal(err)
	}
	if err := alarms.Confirm(sensor.ID); err != nil {
		t.Fatal(err)
	}
	if err := alarms.Sample(sensor.ID, 5); err != nil {
		t.Fatal(err)
	}
	state, err := alarms.State(sensor.ID)
	if err != nil {
		t.Fatal(err)
	}
	if state.Latched {
		t.Fatal("gas alarm latch must clear after confirmation when the reading returns to normal")
	}
}
