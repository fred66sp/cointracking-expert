package cache

import (
	"path/filepath"
	"testing"
	"time"
)

func TestManagerCascade(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "cointracking-mcp.db")
	store, err := OpenStore(dbPath)
	if err != nil {
		t.Fatalf("OpenStore: %v", err)
	}
	defer store.Close()

	l1 := NewLRU(10)
	mgr := NewManager(l1, store)

	if _, _, ok := mgr.Lookup("missing"); ok {
		t.Fatal("expected miss on empty cache")
	}

	mgr.Store("k1", "getBalance", []byte("v1"), time.Minute, func(err error) { t.Errorf("disk write error: %v", err) })

	e, src, ok := mgr.Lookup("k1")
	if !ok || src != FromMemory || string(e.Value) != "v1" {
		t.Fatalf("expected memory hit, got ok=%v src=%v e=%v", ok, src, e)
	}

	// Simulate a cold L1 (e.g. after restart): value should still be
	// found via L2 and promoted back into L1.
	l1.Clear()
	waitFor(t, func() bool { n, _ := store.Count(); return n >= 1 })

	e, src, ok = mgr.Lookup("k1")
	if !ok || src != FromDisk {
		t.Fatalf("expected disk hit after L1 clear, got ok=%v src=%v", ok, src)
	}
	if _, ok := l1.Get("k1"); !ok {
		t.Fatal("expected disk hit to promote entry back into L1")
	}
}

func TestManagerInvalidate(t *testing.T) {
	l1 := NewLRU(10)
	mgr := NewManager(l1, nil)

	mgr.Store("k1", "getTrades", []byte("1"), time.Minute, nil)
	mgr.Store("k2", "getBalance", []byte("2"), time.Minute, nil)

	removed, err := mgr.Invalidate(MatchPattern("getTrades*"))
	if err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	if _, _, ok := mgr.Lookup("k1"); ok {
		t.Fatal("expected getTrades entry invalidated")
	}
	if _, _, ok := mgr.Lookup("k2"); !ok {
		t.Fatal("expected getBalance entry to survive")
	}
}
