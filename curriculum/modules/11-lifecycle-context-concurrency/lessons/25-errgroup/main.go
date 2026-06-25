package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func fetchURL(ctx context.Context, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(50 * time.Millisecond):
		if url == "https://bad.example" {
			return fmt.Errorf("failed to fetch %s: 500", url)
		}
		return nil
	}
}

func main() {
	urls := []string{
		"https://example.com",
		"https://api.example.com",
		"https://bad.example",
		"https://cdn.example.com",
	}

	g, ctx := errgroup.WithContext(context.Background())

	for _, url := range urls {
		url := url
		g.Go(func() error {
			return fetchURL(ctx, url)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Fetch failed: %v\n", err)
		return
	}
	fmt.Println("All fetches succeeded")
}
