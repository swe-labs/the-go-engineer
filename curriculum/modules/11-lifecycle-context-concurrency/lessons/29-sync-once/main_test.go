package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestConfigOnce(t *testing.T) {
	tests := []struct {
		name     string
		workers  int
		wantAddr string
		wantPort int
	}{
		{"single goroutine", 1, "0.0.0.0", 8080},
		{"ten goroutines", 10, "0.0.0.0", 8080},
		{"hundred goroutines", 100, "0.0.0.0", 8080},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configOnce = sync.Once{}
			config = nil
			var wg sync.WaitGroup
			wg.Add(tt.workers)
			for i := 0; i < tt.workers; i++ {
				go func() {
					defer wg.Done()
					cfg := LoadConfig()
					if cfg.Addr != tt.wantAddr {
						t.Errorf("Addr = %q, want %q", cfg.Addr, tt.wantAddr)
					}
					if cfg.Port != tt.wantPort {
						t.Errorf("Port = %d, want %d", cfg.Port, tt.wantPort)
					}
				}()
			}
			wg.Wait()
		})
	}
}

func TestOnceExecutesExactlyOnce(t *testing.T) {
	var count int32
	var once sync.Once
	var wg sync.WaitGroup
	n := 100
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			once.Do(func() {
				atomic.AddInt32(&count, 1)
			})
		}()
	}
	wg.Wait()
	if c := atomic.LoadInt32(&count); c != 1 {
		t.Errorf("Do executed %d times, want 1", c)
	}
}
