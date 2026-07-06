package cache

import (
	"sync"
	"testing"
	"time"
)

// TestSweepEvictsExpired confirms Sweep drops entries past their TTL.
func TestSweepEvictsExpired(t *testing.T) {
	c := New[string, int](10, 20*time.Millisecond)
	c.Set("a", 1)
	c.Set("b", 2)
	if got, ok := c.Get("a"); !ok || got != 1 {
		t.Fatalf("Get(a) = %v, %v before TTL", got, ok)
	}

	time.Sleep(30 * time.Millisecond)
	c.Sweep()

	if c.Len() != 0 {
		t.Fatalf("after Sweep, Len = %d, want 0", c.Len())
	}
	if _, ok := c.Get("a"); ok {
		t.Error("Get(a) should miss after Sweep evicted it")
	}
}

// TestSweepKeepsFresh confirms Sweep leaves non-expired entries alone.
func TestSweepKeepsFresh(t *testing.T) {
	c := New[string, int](10, time.Hour)
	c.Set("alive", 99)
	c.Sweep()
	if c.Len() != 1 {
		t.Fatalf("fresh entry evicted: Len = %d, want 1", c.Len())
	}
	if got, ok := c.Get("alive"); !ok || got != 99 {
		t.Errorf("Get(alive) = %v, %v after Sweep", got, ok)
	}
}

// TestGetEvictsExpiredLazy verifies a Get on an expired entry returns a miss
// and removes it (the upgrade-to-write-lock path in Get).
func TestGetEvictsExpiredLazy(t *testing.T) {
	c := New[string, int](10, 15*time.Millisecond)
	c.Set("x", 7)
	time.Sleep(25 * time.Millisecond)

	if _, ok := c.Get("x"); ok {
		t.Fatal("Get on expired entry should miss")
	}
	if c.Len() != 0 {
		t.Fatalf("expired entry not lazily removed: Len = %d", c.Len())
	}
}

// TestConcurrentGetSet hammers Get and Set from multiple goroutines to
// surface any race in the RLock/MoveToFront-skipped path. Run with -race.
func TestConcurrentGetSet(t *testing.T) {
	c := New[int, int](50, time.Hour)
	const workers = 8
	const iters = 200

	var wg sync.WaitGroup
	wg.Add(workers * 2)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				c.Set(i%50, i)
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				if _, ok := c.Get(i % 50); ok {
					// hit; sweep occasionally to exercise the lock
				}
			}
		}()
	}
	wg.Wait()

	if c.Len() > 50 {
		t.Fatalf("Len = %d, exceeded maxSize 50", c.Len())
	}
}

// TestLRUEvictionOnOverflow confirms Set beyond maxSize evicts the oldest
// (insertion-order LRU after Item C's read-path change).
func TestLRUEvictionOnOverflow(t *testing.T) {
	c := New[int, int](2, time.Hour)
	c.Set(1, 1)
	c.Set(2, 2)
	c.Set(3, 3) // should evict 1 (oldest by insertion)

	if _, ok := c.Get(1); ok {
		t.Error("entry 1 should have been evicted")
	}
	if _, ok := c.Get(3); !ok {
		t.Error("entry 3 should still be present")
	}
}