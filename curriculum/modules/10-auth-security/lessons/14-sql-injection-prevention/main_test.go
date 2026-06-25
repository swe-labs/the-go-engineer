package main

import (
	"testing"
)

func TestUnsafeQuery(t *testing.T) {
	baseQuery := "SELECT email FROM users WHERE username = '{input}'"

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal input unchanged structure",
			input: "alice",
			want:  "SELECT email FROM users WHERE username = 'alice'",
		},
		{
			name:  "injection alters structure",
			input: "alice' OR '1'='1",
			want:  "SELECT email FROM users WHERE username = 'alice' OR '1'='1'",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := unsafeQuery(baseQuery, tc.input)
			if got != tc.want {
				t.Errorf("unsafeQuery(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestSafeQuery(t *testing.T) {
	baseQuery := "SELECT email FROM users WHERE username = ?"

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal input",
			input: "alice",
			want:  "SELECT email FROM users WHERE username = 'alice'",
		},
		{
			name:  "injection input escaped",
			input: "alice' OR '1'='1",
			want:  "SELECT email FROM users WHERE username = 'alice'' OR ''1''=''1'",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := safeQuery(baseQuery, tc.input)
			if got != tc.want {
				t.Errorf("safeQuery(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestEscapeSQLString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "alice", want: "alice"},
		{input: "alice' OR '1'='1", want: "alice'' OR ''1''=''1"},
		{input: "test; DROP TABLE", want: "test; DROP TABLE"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := escapeSQLString(tc.input)
			if got != tc.want {
				t.Errorf("escapeSQLString(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsInjection(t *testing.T) {
	tests := []struct {
		query string
		want  bool
	}{
		{query: "SELECT * FROM users WHERE name = 'alice'", want: false},
		{query: "SELECT * FROM users WHERE name = 'alice' OR '1'='1'", want: true},
		{query: "SELECT * FROM users; DROP TABLE users;", want: true},
		{query: "SELECT * FROM users UNION SELECT * FROM admins", want: true},
		{query: "SELECT * FROM users WHERE name LIKE '%input%'", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.query[:20], func(t *testing.T) {
			got := isInjection(tc.query)
			if got != tc.want {
				t.Errorf("isInjection(%q) = %v, want %v", tc.query, got, tc.want)
			}
		})
	}
}

func TestSimulateDBLookup(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{query: "SELECT * FROM users WHERE name = 'alice'", want: "alice@example.com"},
		{query: "SELECT * FROM users WHERE name = 'alice' OR '1'='1'", want: "BREACHED: All users returned!"},
		{query: "SELECT * FROM users WHERE name = 'nonexistent'", want: "no results"},
		{query: "SELECT * FROM users WHERE name = 'bob'", want: "bob@example.com"},
		{query: "SELECT * FROM users; DROP TABLE users;", want: "BREACHED: All users returned!"},
	}

	for _, tc := range tests {
		t.Run(tc.query[:15], func(t *testing.T) {
			got := simulateDBLookup(tc.query)
			if got != tc.want {
				t.Errorf("simulateDBLookup(%q) = %q, want %q", tc.query, got, tc.want)
			}
		})
	}
}
