//go:build windows

package main

// GetSystemDetails is compiled ONLY on Windows.
func GetSystemDetails() string {
	return "[Windows] Using Win32 APIs under the hood..."
}
