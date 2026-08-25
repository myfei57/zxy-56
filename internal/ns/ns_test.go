package ns

import (
	"testing"

	"labvent/internal/store"
)

func TestNamespaceZoneSummary(t *testing.T) {
	blobs := store.New(t.TempDir())
	namespaces := NewNamespaceService(blobs)
	zones := NewZoneService(blobs)
	item, err := namespaces.Register("A 区", "A")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := zones.Register(item.ID, "A-1", 120); err != nil {
		t.Fatal(err)
	}
	if _, err := zones.Register(item.ID, "A-2", 80); err != nil {
		t.Fatal(err)
	}
	summary, err := namespaces.Summary(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if summary.ZoneCount != 2 {
		t.Fatalf("expected 2 zones, got %d", summary.ZoneCount)
	}
	if summary.TotalArea != 200 {
		t.Fatalf("expected area 200, got %f", summary.TotalArea)
	}
	area, err := zones.Area(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if area != 200 {
		t.Fatalf("expected zone area 200, got %f", area)
	}
}
