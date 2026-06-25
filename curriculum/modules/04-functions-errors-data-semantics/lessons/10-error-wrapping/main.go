package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("resource not found")

func fetchData(source string) error {
	switch {
	case source == "":
		return fmt.Errorf("fetch %s: %w", source, ErrNotFound)
	case source == "badhost":
		return fmt.Errorf("fetch %s: %w", source, errors.New("connection refused"))
	case source == "valid":
		return nil
	default:
		return fmt.Errorf("fetch %s: %w", source, errors.New("unknown source"))
	}
}

func main() {
	sources := []string{"", "badhost", "valid"}
	for _, src := range sources {
		err := fetchData(src)
		if err == nil {
			fmt.Printf("fetchData(%q): OK\n", src)
			continue
		}
		fmt.Printf("fetchData(%q): %v\n", src, err)
		fmt.Printf("  errors.Is(err, ErrNotFound) = %v\n", errors.Is(err, ErrNotFound))
		fmt.Printf("  errors.Unwrap(err)          = %v\n", errors.Unwrap(err))
	}
}
