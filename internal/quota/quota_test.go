package quota

import (
	"testing"

	"labvent/internal/store"
)

func TestQuotaReserveRelease(t *testing.T) {
	service := NewQuotaService(store.New(t.TempDir()))
	if _, err := service.Register("lab1", 10); err != nil {
		t.Fatal(err)
	}
	if err := service.Reserve("lab1", 4); err != nil {
		t.Fatal(err)
	}
	used, err := service.Usage("lab1")
	if err != nil {
		t.Fatal(err)
	}
	if used != 4 {
		t.Fatalf("expected usage 4, got %f", used)
	}
	if err := service.Release("lab1", 2); err != nil {
		t.Fatal(err)
	}
	used, err = service.Usage("lab1")
	if err != nil {
		t.Fatal(err)
	}
	if used != 2 {
		t.Fatalf("expected usage 2, got %f", used)
	}
}

func TestQuotaExceeded(t *testing.T) {
	service := NewQuotaService(store.New(t.TempDir()))
	if _, err := service.Register("lab2", 3); err != nil {
		t.Fatal(err)
	}
	if err := service.Reserve("lab2", 4); err == nil {
		t.Fatal("reserve beyond quota must fail")
	}
}
