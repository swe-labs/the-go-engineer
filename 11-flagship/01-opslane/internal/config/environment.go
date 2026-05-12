// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package config

import (
	"net/netip"
	"strconv"
	"strings"
	"time"
)

// LookupFunc (Type): abstracts environment variable lookup to enable testing.
// It matches the signature of os.LookupEnv.
type LookupFunc func(string) (string, bool)

// stringFromEnv (Function): reads a string environment variable through the lookup
// abstraction, returning the fallback when the variable is unset or empty.
func stringFromEnv(lookup LookupFunc, key, fallback string) string {
	if value, ok := lookup(key); ok && value != "" {
		return value
	}

	return fallback
}

// durationFromEnv (Function): reads a duration environment variable with a default
// fallback, returning an error if the value is present but not a valid Go duration.
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

// intFromEnv (Function): reads an integer environment variable with a default
// fallback, returning an error when parsing fails.
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

// floatFromEnv (Function): reads a float64 environment variable with a default
// fallback, returning an error when parsing fails.
func floatFromEnv(lookup LookupFunc, key string, fallback float64) (float64, error) {
	raw, ok := lookup(key)
	if !ok || raw == "" {
		return fallback, nil
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// cidrPrefixesFromEnv (Function): reads a comma-separated CIDR prefix list from an
// environment variable, returning an error if any prefix fails to parse.
func cidrPrefixesFromEnv(lookup LookupFunc, key string) ([]netip.Prefix, error) {
	raw, ok := lookup(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	prefixes := make([]netip.Prefix, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}

		prefix, err := netip.ParsePrefix(value)
		if err != nil {
			return nil, err
		}

		prefixes = append(prefixes, prefix)
	}

	return prefixes, nil
}
