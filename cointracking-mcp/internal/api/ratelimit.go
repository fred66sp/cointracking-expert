package api

import (
	"sync"
	"time"
)

// RateTracker tracks API calls in a rolling 1-hour window so the server can
// report consumption in stats and warn before hitting CoinTracking's limit.
// It does not block calls itself (SPEC 02: "sin reintento automático").
type RateTracker struct {
	mu      sync.Mutex
	limit   int
	calls   []time.Time
	byMeth  map[string]*methodStats
	lastAPI time.Time
}

type methodStats struct {
	Count    int
	LastCall time.Time
}

func NewRateTracker(limit int) *RateTracker {
	return &RateTracker{limit: limit, byMeth: make(map[string]*methodStats)}
}

func (r *RateTracker) RecordCall(method string) {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, now)
	r.lastAPI = now
	st, ok := r.byMeth[method]
	if !ok {
		st = &methodStats{}
		r.byMeth[method] = st
	}
	st.Count++
	st.LastCall = now
}

// UsedInWindow returns how many calls landed in the trailing 1-hour window.
func (r *RateTracker) UsedInWindow() int {
	cutoff := time.Now().Add(-time.Hour)
	r.mu.Lock()
	defer r.mu.Unlock()
	kept := r.calls[:0]
	count := 0
	for _, t := range r.calls {
		if t.After(cutoff) {
			kept = append(kept, t)
			count++
		}
	}
	r.calls = kept
	return count
}

func (r *RateTracker) Limit() int { return r.limit }

type MethodCallStat struct {
	Method   string
	Count    int
	LastCall time.Time
}

func (r *RateTracker) MethodStats() []MethodCallStat {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]MethodCallStat, 0, len(r.byMeth))
	for m, st := range r.byMeth {
		out = append(out, MethodCallStat{Method: m, Count: st.Count, LastCall: st.LastCall})
	}
	return out
}
