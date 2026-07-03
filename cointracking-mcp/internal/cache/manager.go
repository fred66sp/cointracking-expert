package cache

import "time"

// Source identifies where a value was served from, for logging/telemetry.
type Source string

const (
	FromMemory Source = "FROM_MEMORY"
	FromDisk   Source = "FROM_DISK"
	FromAPI    Source = "FROM_API"
)

// Manager implements the L1 (memory) -> L2 (disk) -> caller cascade from
// SPEC/03-cache-strategy.md ("¿En memoria? -> ¿En disco? -> Llamar API").
// It does not call the API itself; Lookup reports a miss and the caller
// (a tool handler) is responsible for fetching and calling Store.
type Manager struct {
	L1 *LRU
	L2 *Store // nil if disk persistence disabled
}

func NewManager(l1 *LRU, l2 *Store) *Manager {
	return &Manager{L1: l1, L2: l2}
}

// Lookup returns the cached value and its source, or ok=false on a full miss.
func (m *Manager) Lookup(key string) (*Entry, Source, bool) {
	if e, ok := m.L1.Get(key); ok {
		return e, FromMemory, true
	}
	if m.L2 != nil {
		if e, ok := m.L2.Get(key); ok {
			m.L1.Set(key, e) // promote to L1
			return e, FromDisk, true
		}
	}
	return nil, "", false
}

// Store writes value to L1 immediately and to L2 asynchronously.
func (m *Manager) Store(key, method string, value []byte, ttl time.Duration, onDiskErr func(error)) {
	entry := &Entry{Method: method, Value: value, ExpiresAt: time.Now().Add(ttl)}
	m.L1.Set(key, entry)
	if m.L2 != nil {
		m.L2.SetAsync(key, entry, onDiskErr)
	}
}

// Invalidate removes matching entries from both L1 and L2, returning the
// total count removed (L1 + L2; an entry present in both counts twice,
// matching "entradas eliminadas" as a raw operation count).
func (m *Manager) Invalidate(match func(method, key string) bool) (int, error) {
	removed := m.L1.InvalidateFunc(match)
	if m.L2 != nil {
		n, err := m.L2.InvalidateFunc(match)
		removed += n
		if err != nil {
			return removed, err
		}
	}
	return removed, nil
}
