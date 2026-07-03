package cache

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// Store is the L2 (on-disk, per-project) persistent cache backed by SQLite.
// Writes are async (goroutine) so they never block a tool response, per
// SPEC/03-cache-strategy.md. pendingWrites tracks in-flight writes so Flush
// can guarantee a graceful shutdown (SPEC: "Al shutdown: volcar caché a
// disco") doesn't drop data that was still in flight.
type Store struct {
	db            *sql.DB
	pendingWrites sync.WaitGroup
}

// OpenStore opens (creating if needed) the SQLite cache DB at
// dbPath and ensures its schema exists.
func OpenStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening cache db: %w", err)
	}
	db.SetMaxOpenConns(1) // modernc.org/sqlite: single-writer safety

	const schema = `
CREATE TABLE IF NOT EXISTS cache (
  key TEXT PRIMARY KEY,
  method TEXT NOT NULL,
  value BLOB NOT NULL,
  cached_at INTEGER NOT NULL,
  expires_at INTEGER NOT NULL
);`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("creating cache schema: %w", err)
	}
	return &Store{db: db}, nil
}

// Flush blocks until all in-flight async writes started via SetAsync have
// completed. Call before Close (or at close_project time) to avoid losing
// writes that were still in flight.
func (s *Store) Flush() {
	s.pendingWrites.Wait()
}

func (s *Store) Close() error {
	s.Flush()
	return s.db.Close()
}

// Get returns the entry for key if present and not expired. An expired
// row is deleted and treated as a miss.
func (s *Store) Get(key string) (*Entry, bool) {
	row := s.db.QueryRow(`SELECT method, value, expires_at FROM cache WHERE key = ?`, key)
	var method string
	var value []byte
	var expiresAtUnix int64
	if err := row.Scan(&method, &value, &expiresAtUnix); err != nil {
		return nil, false
	}
	expiresAt := time.Unix(expiresAtUnix, 0)
	if time.Now().After(expiresAt) {
		_, _ = s.db.Exec(`DELETE FROM cache WHERE key = ?`, key)
		return nil, false
	}
	return &Entry{Method: method, Value: value, ExpiresAt: expiresAt}, true
}

// SetAsync persists key/entry to disk without blocking the caller. Errors
// are reported via errFn (typically a logger call), if provided.
func (s *Store) SetAsync(key string, entry *Entry, errFn func(error)) {
	s.pendingWrites.Add(1)
	go func() {
		defer s.pendingWrites.Done()
		_, err := s.db.Exec(
			`INSERT INTO cache (key, method, value, cached_at, expires_at) VALUES (?, ?, ?, ?, ?)
			 ON CONFLICT(key) DO UPDATE SET method=excluded.method, value=excluded.value,
			   cached_at=excluded.cached_at, expires_at=excluded.expires_at`,
			key, entry.Method, entry.Value, time.Now().Unix(), entry.ExpiresAt.Unix(),
		)
		if err != nil && errFn != nil {
			errFn(err)
		}
	}()
}

// PurgeExpired deletes all expired rows, returning the count removed.
// Called at startup and periodically per SPEC 03.
func (s *Store) PurgeExpired() (int, error) {
	res, err := s.db.Exec(`DELETE FROM cache WHERE expires_at < ?`, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

// InvalidateFunc deletes rows whose method matches, using match for
// per-row filtering (key is available too, e.g. for "*" patterns).
func (s *Store) InvalidateFunc(match func(method, key string) bool) (int, error) {
	rows, err := s.db.Query(`SELECT key, method FROM cache`)
	if err != nil {
		return 0, err
	}
	var toDelete []string
	for rows.Next() {
		var key, method string
		if err := rows.Scan(&key, &method); err != nil {
			rows.Close()
			return 0, err
		}
		if match(method, key) {
			toDelete = append(toDelete, key)
		}
	}
	rows.Close()

	for _, k := range toDelete {
		if _, err := s.db.Exec(`DELETE FROM cache WHERE key = ?`, k); err != nil {
			return len(toDelete), err
		}
	}
	return len(toDelete), nil
}

// LoadAll returns every non-expired row, for warming the L1 LRU at startup.
func (s *Store) LoadAll() (map[string]*Entry, error) {
	rows, err := s.db.Query(`SELECT key, method, value, expires_at FROM cache WHERE expires_at >= ?`, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]*Entry)
	for rows.Next() {
		var key, method string
		var value []byte
		var expiresAtUnix int64
		if err := rows.Scan(&key, &method, &value, &expiresAtUnix); err != nil {
			return nil, err
		}
		out[key] = &Entry{Method: method, Value: value, ExpiresAt: time.Unix(expiresAtUnix, 0)}
	}
	return out, nil
}

// Count returns the number of rows currently persisted.
func (s *Store) Count() (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM cache`).Scan(&n)
	return n, err
}
