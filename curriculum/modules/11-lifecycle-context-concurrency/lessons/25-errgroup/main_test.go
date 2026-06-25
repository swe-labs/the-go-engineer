package main

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestErrgroup(t *testing.T) {
	tests := []struct {
		name    string
		urls    []string
		wantErr bool
	}{
		{
			name:    "all succeed",
			urls:    []string{"https://example.com", "https://api.example.com"},
			wantErr: false,
		},
		{
			name:    "one fails",
			urls:    []string{"https://example.com", "https://bad.example"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, ctx := errgroup.WithContext(context.Background())
			for _, url := range tt.urls {
				url := url
				g.Go(func() error {
					return fetchURL(ctx, url)
				})
			}
			err := g.Wait()
			if (err != nil) != tt.wantErr {
				t.Errorf("g.Wait() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestErrgroupContextCancellation(t *testing.T) {
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	g.Go(func() error {
		return errors.New("first failure")
	})
	err := g.Wait()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
