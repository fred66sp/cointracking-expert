package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLRUBasicGetSet(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", &Entry{Method: "m", Value: []byte("1"), ExpiresAt: time.Now().Add(time.Minute)})
	e, ok := c.Get("a")
	if !ok || string(e.Value) != "1" {
		t.Fatalf("expected hit with value 1, got ok=%v e=%v", ok, e)
	}
}

func TestLRUEviction(t *testing.T) {
	c := NewLRU(2)
	future := time.Now().Add(time.Minute)
	c.Set("a", &Entry{Method: "m", Value: []byte("1"), ExpiresAt: future})
	c.Set("b", &Entry{Method: "m", Value: []byte("2"), ExpiresAt: future})
	c.Set("c", &Entry{Method: "m", Value: []byte("3"), ExpiresAt: future}) // evicts "a" (LRU)

	if _, ok := c.Get("a"); ok {
		t.Fatal("expected 'a' to be evicted")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatal("expected 'b' to still be present")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatal("expected 'c' to still be present")
	}
}

func TestLRUTTLExpiry(t *testing.T) {
	c := NewLRU(10)
	c.Set("a", &Entry{Method: "m", Value: []byte("1"), ExpiresAt: time.Now().Add(50 * time.Millisecond)})
	if _, ok := c.Get("a"); !ok {
		t.Fatal("expected hit before expiry")
	}
	time.Sleep(80 * time.Millisecond)
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected miss after TTL expiry")
	}
}

func TestLRUInvalidateFunc(t *testing.T) {
	c := NewLRU(10)
	future := time.Now().Add(time.Minute)
	c.Set("k1", &Entry{Method: "getTrades", Value: []byte("1"), ExpiresAt: future})
	c.Set("k2", &Entry{Method: "getTrades", Value: []byte("2"), ExpiresAt: future})
	c.Set("k3", &Entry{Method: "getBalance", Value: []byte("3"), ExpiresAt: future})

	removed := c.InvalidateFunc(func(method, key string) bool { return method == "getTrades" })
	if removed != 2 {
		t.Fatalf("expected 2 removed, got %d", removed)
	}
	if _, ok := c.Get("k3"); !ok {
		t.Fatal("expected getBalance entry to survive invalidation")
	}
}

func TestLRUConcurrency(t *testing.T) {
	c := NewLRU(100)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("k%d", i%10)
			for j := 0; j < 100; j++ {
				c.Set(key, &Entry{Method: "m", Value: []byte("v"), ExpiresAt: time.Now().Add(time.Minute)})
				c.Get(key)
			}
		}(i)
	}
	wg.Wait()
}
