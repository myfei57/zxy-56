package audit

import (
	"testing"

	"labvent/internal/store"
)

func TestAuditRecordAndList(t *testing.T) {
	svc := NewService(store.New(t.TempDir()))
	if err := svc.Record("fan", "f1", "start", "morning"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Record("gas", "g1", "sample", "12ppm"); err != nil {
		t.Fatal(err)
	}
	entries, err := svc.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	fanEntries, err := svc.BySource("fan")
	if err != nil {
		t.Fatal(err)
	}
	if len(fanEntries) != 1 || fanEntries[0].Action != "start" {
		t.Fatalf("unexpected fan entries: %v", fanEntries)
	}
	count, err := svc.Count("g1")
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}
}

func TestAuditRecentLimit(t *testing.T) {
	svc := NewService(store.New(t.TempDir()))
	for i := 0; i < 5; i++ {
		if err := svc.Record("hood", "h1", "airflow", "ok"); err != nil {
			t.Fatal(err)
		}
	}
	entries, err := svc.Recent(3)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 recent entries, got %d", len(entries))
	}
}
