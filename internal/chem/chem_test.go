package chem

import (
	"testing"
	"time"

	"labvent/internal/store"
)

func TestChemRegisterGet(t *testing.T) {
	blobs := store.New(t.TempDir())
	expiry := NewExpiryService(blobs)
	classify := NewClassifyService(blobs)
	service := NewChemService(blobs, expiry, classify)
	expiryAt := time.Now().Add(30 * 24 * time.Hour)
	item, err := service.Register("丙酮", "ordinary", expiryAt, 3)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := service.Get(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Name != "丙酮" || loaded.Class != "ordinary" {
		t.Fatalf("unexpected chemical: %v", loaded)
	}
	updated, err := service.UpdateQuantity(item.ID, 5)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Quantity != 5 {
		t.Fatalf("unexpected quantity: %f", updated.Quantity)
	}
}
