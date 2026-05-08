package events

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBusConcurrentPublishAndClose(t *testing.T) {
	t.Parallel()
	bus := NewBus(64)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = bus.Publish(ctx, Event{Type: TypeOrderCreated, TenantID: 1})
		}()
	}

	time.Sleep(10 * time.Millisecond)
	bus.Close()
	wg.Wait()
}
