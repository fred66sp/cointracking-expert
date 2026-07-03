package cache

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStoreSetGetPurge(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "cointracking-mcp.db")
	store, err := OpenStore(dbPath)
	if err != nil {
		t.Fatalf("OpenStore: %v", err)
	}
	defer store.Close()

	done := make(chan struct{})
	entry := &Entry{Method: "getBalance", Value: []byte(`{"ok":true}`), ExpiresAt: time.Now().Add(time.Minute)}
	store.SetAsync("k1", entry, func(err error) { t.Errorf("unexpected write error: %v", err) })
	// SetAsync is fire-and-forget; poll briefly for the write to land.
	waitFor(t, func() bool {
		_, ok := store.Get("k1")
		return ok
	})
	close(done)

	got, ok := store.Get("k1")
	if !ok {
		t.Fatal("expected hit after async write completed")
	}
	if string(got.Value) != `{"ok":true}` {
		t.Fatalf("unexpected value: %s", got.Value)
	}

	expired := &Entry{Method: "getBalance", Value: []byte("x"), ExpiresAt: time.Now().Add(-time.Second)}
	store.SetAsync("k2", expired, nil)
	waitFor(t, func() bool {
		n, _ := store.Count()
		return n >= 2
	})
	if _, ok := store.Get("k2"); ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestStoreLoadAll(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "cointracking-mcp.db")
	store, err := OpenStore(dbPath)
	if err != nil {
		t.Fatalf("OpenStore: %v", err)
	}
	defer store.Close()

	future := time.Now().Add(time.Minute)
	store.SetAsync("k1", &Entry{Method: "getTrades", Value: []byte("1"), ExpiresAt: future}, nil)
	store.SetAsync("k2", &Entry{Method: "getBalance", Value: []byte("2"), ExpiresAt: future}, nil)
	waitFor(t, func() bool { n, _ := store.Count(); return n == 2 })

	all, err := store.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}
