package main

import (
	"errors"
	"fmt"
)

type HTTPError struct {
	StatusCode int
	Body       string
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", h.StatusCode, h.Body)
}

func fetchURL(url string) error {
	switch url {
	case "/missing":
		return fmt.Errorf("fetch %s: %w", url, &HTTPError{StatusCode: 404, Body: "not found"})
	case "/broken":
		return fmt.Errorf("fetch %s: %w", url, &HTTPError{StatusCode: 500, Body: "server error"})
	case "/ok":
		return nil
	default:
		return fmt.Errorf("fetch %s: %w", url, &HTTPError{StatusCode: 400, Body: "bad request"})
	}
}

func main() {
	urls := []string{"/missing", "/broken", "/ok", "/other"}
	for _, url := range urls {
		err := fetchURL(url)
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			fmt.Printf("URL %s -> HTTP %d: %q\n", url, httpErr.StatusCode, httpErr.Body)
		} else if err == nil {
			fmt.Printf("URL %s -> OK\n", url)
		} else {
			fmt.Printf("URL %s -> unknown: %v\n", url, err)
		}
	}
}
