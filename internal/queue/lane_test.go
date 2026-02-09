package queue

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestLaneSerialExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mgr := NewManager(ctx, 64, time.Minute)

	var mu sync.Mutex
	var order []int

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		mgr.Enqueue(Task{
			SessionID: "test-session",
			Fn: func(ctx context.Context) error {
				// Simulate work.
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				order = append(order, i)
				mu.Unlock()
				wg.Done()
				return nil
			},
		})
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	// Tasks should execute in order since they go through the same lane.
	for i, v := range order {
		if v != i {
			t.Errorf("expected order[%d] = %d, got %d", i, i, v)
		}
	}
}

func TestDifferentSessionsParallel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mgr := NewManager(ctx, 64, time.Minute)

	var wg sync.WaitGroup
	start := time.Now()

	for s := 0; s < 3; s++ {
		wg.Add(1)
		sid := string(rune('A' + s))
		mgr.Enqueue(Task{
			SessionID: sid,
			Fn: func(ctx context.Context) error {
				time.Sleep(50 * time.Millisecond)
				wg.Done()
				return nil
			},
		})
	}

	wg.Wait()
	elapsed := time.Since(start)

	// 3 sessions running in parallel should take ~50ms, not 150ms.
	if elapsed > 120*time.Millisecond {
		t.Errorf("expected parallel execution, took %v", elapsed)
	}
}
