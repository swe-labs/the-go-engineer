// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package config

import (
	"strconv"
	"time"
)

type LookupFunc func(string) (string, bool)

func stringFromEnv(lookup LookupFunc, key, fallback string) string {
	if value, ok := lookup(key); ok && value != "" {
		return value
	}

	return fallback
}

func durationFromEnv(lookup LookupFunc, key string, fallback time.Duration) (time.Duration, error) {
	raw, ok := lookup(key)
	if !ok || raw == "" {
		return fallback, nil
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func intFromEnv(lookup LookupFunc, key string, fallback int) (int, error) {
	raw, ok := lookup(key)
	if !ok || raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}

	return value, nil
}
