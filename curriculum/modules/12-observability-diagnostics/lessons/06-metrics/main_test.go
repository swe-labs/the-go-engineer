package main

import (
	"testing"
)

func TestMetricsStore_RecordRequest(t *testing.T) {
	tests := []struct {
		name       string
		operations []struct {
			durationMs int64
			isError    bool
		}
		wantRequests int64
		wantErrors   int64
		wantAvgMin   int64
		wantAvgMax   int64
	}{
		{
			name: "all successful",
			operations: []struct {
				durationMs int64
				isError    bool
			}{
				{100, false},
				{200, false},
				{300, false},
			},
			wantRequests: 3,
			wantErrors:   0,
			wantAvgMin:   200,
			wantAvgMax:   200,
		},
		{
			name: "some errors",
			operations: []struct {
				durationMs int64
				isError    bool
			}{
				{50, false},
				{5000, true},
				{100, false},
			},
			wantRequests: 3,
			wantErrors:   1,
			wantAvgMin:   1716,
			wantAvgMax:   1717,
		},
		{
			name: "all errors",
			operations: []struct {
				durationMs int64
				isError    bool
			}{
				{1000, true},
				{2000, true},
			},
			wantRequests: 2,
			wantErrors:   2,
			wantAvgMin:   1500,
			wantAvgMax:   1500,
		},
		{
			name: "single request",
			operations: []struct {
				durationMs int64
				isError    bool
			}{
				{42, false},
			},
			wantRequests: 1,
			wantErrors:   0,
			wantAvgMin:   42,
			wantAvgMax:   42,
		},
		{
			name: "no operations",
			operations: []struct {
				durationMs int64
				isError    bool
			}{},
			wantRequests: 0,
			wantErrors:   0,
			wantAvgMin:   0,
			wantAvgMax:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMetricsStore()
			for _, op := range tt.operations {
				store.RecordRequest(op.durationMs, op.isError)
			}
			reqs, errs, avg := store.Snapshot()
			if reqs != tt.wantRequests {
				t.Errorf("requests = %d, want %d", reqs, tt.wantRequests)
			}
			if errs != tt.wantErrors {
				t.Errorf("errors = %d, want %d", errs, tt.wantErrors)
			}
			if avg < tt.wantAvgMin || avg > tt.wantAvgMax {
				t.Errorf("avg = %d, want between %d and %d", avg, tt.wantAvgMin, tt.wantAvgMax)
			}
		})
	}
}
