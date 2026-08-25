package gas

import (
	"testing"

	"labvent/internal/store"
)

func TestSensorSample(t *testing.T) {
	blobs := store.New(t.TempDir())
	gasService := NewGasService(blobs)
	samples := NewSampleService(blobs)
	item, err := gasService.Register("A 区传感器", "room1", 10)
	if err != nil {
		t.Fatal(err)
	}
	if err := samples.Record(item.ID, 8.5); err != nil {
		t.Fatal(err)
	}
	reading, err := samples.Latest(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if reading.PPM != 8.5 {
		t.Fatalf("unexpected ppm: %f", reading.PPM)
	}
	items, err := gasService.List("room1")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 sensor, got %d", len(items))
	}
}
