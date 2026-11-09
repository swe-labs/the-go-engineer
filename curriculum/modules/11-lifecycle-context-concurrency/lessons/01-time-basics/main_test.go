package main

import (
	"testing"
	"time"
)

func TestParseRFC3339(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"valid date", "2026-06-01T12:00:00Z", "2026-06-01T12:00:00Z", false},
		{"valid with tz", "2026-01-15T08:30:00+05:30", "2026-01-15T08:30:00+05:30", false},
		{"invalid format", "01-06-2026", "", true},
		{"empty string", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := time.Parse(time.RFC3339, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr = %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Format(time.RFC3339) != tt.want {
					t.Errorf("Parse(%q) = %q, want %q", tt.input, got.Format(time.RFC3339), tt.want)
				}
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tm := time.Date(2026, 6, 1, 14, 30, 0, 0, time.UTC)
	tests := []struct {
		name   string
		layout string
		want   string
	}{
		{"RFC3339", time.RFC3339, "2026-06-01T14:30:00Z"},
		{"Date only", "2006-01-02", "2026-06-01"},
		{"US style", "01/02/2006", "06/01/2026"},
		{"Time only", "15:04:05", "14:30:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.Format(tt.layout)
			if got != tt.want {
				t.Errorf("Format(%q) = %q, want %q", tt.layout, got, tt.want)
			}
		})
	}
}

func TestSinceUntil(t *testing.T) {
	now := time.Now()
	past := now.Add(-2 * time.Hour)
	future := now.Add(30 * time.Minute)

	since := time.Since(past)
	if since < 2*time.Hour {
		t.Errorf("Since(past) = %v, want >= 2h", since)
	}

	until := time.Until(future)
	if until > 31*time.Minute || until < 29*time.Minute {
		t.Errorf("Until(future) = %v, want ~30m", until)
	}
}
