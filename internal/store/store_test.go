package store

import (
	"testing"
)

func TestFileStoreRoundtrip(t *testing.T) {
	blobs := New(t.TempDir())
	payload := map[string]string{"name": "demo"}
	if err := blobs.Save("demo", "d1", payload); err != nil {
		t.Fatal(err)
	}
	if !blobs.Exists("demo", "d1") {
		t.Fatal("record must exist after save")
	}
	var loaded map[string]string
	if err := blobs.Load("demo", "d1", &loaded); err != nil {
		t.Fatal(err)
	}
	if loaded["name"] != "demo" {
		t.Fatalf("unexpected payload: %v", loaded)
	}
	ids, err := blobs.List("demo")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != "d1" {
		t.Fatalf("unexpected ids: %v", ids)
	}
}
