package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var (
	ErrTransient = errors.New("transient error")
	ErrPermanent = errors.New("permanent error")
)

func build(service string) error {
	if service == "" {
		return fmt.Errorf("build: %w", ErrPermanent)
	}
	return nil
}

func test(service string) error {
	if rand.Intn(2) == 0 {
		return fmt.Errorf("test: %w", ErrTransient)
	}
	return nil
}

func deployToProd(service string) error {
	return nil
}

func isTransient(err error) bool {
	return errors.Is(err, ErrTransient)
}

func retry(attempts int, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if !isTransient(err) {
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("all %d attempts failed: %w", attempts, err)
}

func deploy(service string) error {
	if err := build(service); err != nil {
		return fmt.Errorf("deploy: %w", err)
	}
	if err := retry(3, func() error { return test(service) }); err != nil {
		return fmt.Errorf("deploy: %w", err)
	}
	if err := deployToProd(service); err != nil {
		return fmt.Errorf("deploy: %w", err)
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for _, svc := range []string{"", "my-api"} {
		err := deploy(svc)
		if err == nil {
			fmt.Printf("deploy(%q): success\n", svc)
		} else {
			fmt.Printf("deploy(%q): %v\n", svc, err)
		}
	}
}
