package storage

import (
	"path/filepath"
	"testing"
	"time"
)

func openTemp(t *testing.T) *Store {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestStore_GetExpiration_NotFound(t *testing.T) {
	store := openTemp(t)

	date, err := store.GetExpiration("example.kz")
	if err != nil {
		t.Fatalf("GetExpiration: %v", err)
	}
	if date != nil {
		t.Fatalf("expected nil for unknown domain, got %v", date)
	}
}

func TestStore_SetAndGetExpiration(t *testing.T) {
	store := openTemp(t)

	want := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := store.SetExpiration("example.kz", want); err != nil {
		t.Fatalf("SetExpiration: %v", err)
	}

	got, err := store.GetExpiration("example.kz")
	if err != nil {
		t.Fatalf("GetExpiration: %v", err)
	}
	if got == nil || !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestStore_SetExpiration_Overwrites(t *testing.T) {
	store := openTemp(t)

	old := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := store.SetExpiration("example.kz", old); err != nil {
		t.Fatalf("SetExpiration old: %v", err)
	}
	if err := store.SetExpiration("example.kz", updated); err != nil {
		t.Fatalf("SetExpiration updated: %v", err)
	}

	got, err := store.GetExpiration("example.kz")
	if err != nil {
		t.Fatalf("GetExpiration: %v", err)
	}
	if got == nil || !got.Equal(updated) {
		t.Fatalf("expected %v, got %v", updated, got)
	}
}
