// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

// ============================================================================
// Tests for: JSON Encoder (Streaming)
// ============================================================================
//
// These tests verify that json.NewEncoder produces valid, parseable JSON
// output with the correct field names from struct tags.
//
// RUN: go test -v ./05-packages-io/02-io-and-cli/encoding/3-encoder
// ============================================================================

func TestDeviceLogEncoding(t *testing.T) {
	log := DeviceLog{
		DeviceID:  "server-01",
		Timestamp: 1684501234,
		CPUUsage:  82.5,
		MemUsage:  65.0,
		Status:    "WARNING",
	}

	// Encode into a buffer (simulates writing to a file or HTTP response)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(&log)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode back to verify the JSON is valid and fields are correct
	var decoded DeviceLog
	err = json.Unmarshal(buf.Bytes(), &decoded)
	if err != nil {
		t.Fatalf("Unmarshal of encoded output failed: %v", err)
	}

	if decoded.DeviceID != log.DeviceID {
		t.Errorf("DeviceID = %q, want %q", decoded.DeviceID, log.DeviceID)
	}
	if decoded.Status != log.Status {
		t.Errorf("Status = %q, want %q", decoded.Status, log.Status)
	}
	if decoded.CPUUsage != log.CPUUsage {
		t.Errorf("CPUUsage = %.1f, want %.1f", decoded.CPUUsage, log.CPUUsage)
	}
}

func TestDeviceLogJSONFieldNames(t *testing.T) {
	log := DeviceLog{
		DeviceID:  "test-device",
		Timestamp: 100,
		CPUUsage:  50.0,
		MemUsage:  30.0,
		Status:    "OK",
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(&log)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Parse into a generic map to verify struct tags produce correct JSON keys
	var raw map[string]any
	err = json.Unmarshal(buf.Bytes(), &raw)
	if err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	// Verify that struct tags (e.g., `json:"device_id"`) are respected
	expectedKeys := []string{"device_id", "timestamp", "cpu_usage", "mem_usage", "status"}
	for _, key := range expectedKeys {
		if _, exists := raw[key]; !exists {
			t.Errorf("Expected JSON key %q not found in output. Got keys: %v", key, raw)
		}
	}
}
