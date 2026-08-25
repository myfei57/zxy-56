package gas

import (
	"testing"

	"labvent/internal/audit"
	"labvent/internal/store"
)

// TestAlarmLatchAutoClearAfterConfirm reproduces the operator scenario:
// a gas alarm latches, the operator confirms it, and once the reading
// drops below threshold the latch must auto-clear. Before the fix the
// confirm callback was never invoked, so Confirmed stayed false and the
// alarm never cleared, leading the operator to misdiagnose a sensor fault.
func TestAlarmLatchAutoClearAfterConfirm(t *testing.T) {
	blobs := store.New(t.TempDir())
	auditService := audit.NewService(blobs)
	gasService := NewGasService(blobs)
	latch := NewLatch(blobs)
	alarm := NewAlarmService(gasService, latch, auditService, nil)
	auditService.SetConfirmCallback(latch.MarkConfirmed)

	sensor, err := gasService.Register("A 区传感器", "room1", 10)
	if err != nil {
		t.Fatal(err)
	}

	// Concentration exceeds threshold -> alarm latches.
	if err := alarm.Sample(sensor.ID, 15); err != nil {
		t.Fatal(err)
	}
	state, err := alarm.State(sensor.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !state.Latched {
		t.Fatalf("expected alarm latched after over-threshold sample")
	}
	if state.Confirmed {
		t.Fatalf("alarm must not be confirmed before operator action")
	}

	// Operator confirms the alarm.
	if err := alarm.Confirm(sensor.ID); err != nil {
		t.Fatal(err)
	}
	state, _ = alarm.State(sensor.ID)
	if !state.Latched {
		t.Fatalf("latch must remain raised after confirm until reading recovers")
	}
	if !state.Confirmed {
		t.Fatalf("alarm must be marked confirmed after operator confirm")
	}

	// Reading recovers below threshold -> latch auto-clears.
	if err := alarm.Sample(sensor.ID, 3); err != nil {
		t.Fatal(err)
	}
	state, _ = alarm.State(sensor.ID)
	if state.Latched {
		t.Fatalf("latch should auto-clear once a confirmed alarm recovers below threshold")
	}
}
