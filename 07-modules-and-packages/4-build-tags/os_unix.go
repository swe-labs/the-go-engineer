//go:build linux || darwin

package main

// GetSystemDetails is compiled ONLY on Unix-like systems (Linux or macOS).
func GetSystemDetails() string {
	return "[Unix] Using POSIX system calls under the hood..."
}
