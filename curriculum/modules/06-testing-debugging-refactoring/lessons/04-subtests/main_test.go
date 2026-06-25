package main

import (
	"testing"
	"time"
)

func TestFetchData(t *testing.T) {
	sources := []struct {
		name string
		src  string
	}{
		{name: "cache", src: "cache"},
		{name: "database", src: "db"},
		{name: "api", src: "api"},
	}
	for _, s := range sources {
		s := s
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			got := FetchData(s.src)
			elapsed := time.Since(start)
			if elapsed > 100*time.Millisecond {
				t.Errorf("FetchData(%q) took %v; want <100ms", s.src, elapsed)
			}
			if got == "" {
				t.Errorf("FetchData(%q) returned empty", s.src)
			}
		})
	}
}

func TestProcess(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		got, err := Process(5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 10 {
			t.Errorf("Process(5) = %d; want 10", got)
		}
	})
	t.Run("zero", func(t *testing.T) {
		_, err := Process(0)
		if err == nil {
			t.Error("expected error for id=0")
		}
	})
	t.Run("negative", func(t *testing.T) {
		_, err := Process(-3)
		if err == nil {
			t.Error("expected error for id=-3")
		}
	})
	t.Run("large group", func(t *testing.T) {
		ids := []struct {
			name string
			id   int
		}{
			{name: "large_100", id: 100},
			{name: "large_200", id: 200},
			{name: "large_300", id: 300},
		}
		for _, tc := range ids {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				time.Sleep(10 * time.Millisecond)
				got, err := Process(tc.id)
				if err != nil {
					t.Errorf("unexpected error for id=%d: %v", tc.id, err)
				}
				if got != tc.id*2 {
					t.Errorf("Process(%d) = %d; want %d", tc.id, got, tc.id*2)
				}
			})
		}
	})
}
