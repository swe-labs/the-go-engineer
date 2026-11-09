package main

import (
	"testing"
)

func TestUnsafeRenderer(t *testing.T) {
	r := &UnsafeRenderer{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal name is output verbatim",
			input: "Alice",
			want:  "<html><body><h1>Hello, Alice!</h1></body></html>",
		},
		{
			name:  "script tag is not escaped",
			input: "<script>alert(1)</script>",
			want:  "<html><body><h1>Hello, <script>alert(1)</script>!</h1></body></html>",
		},
		{
			name:  "event handler is not escaped",
			input: "<img src=x onerror=alert(1)>",
			want:  "<html><body><h1>Hello, <img src=x onerror=alert(1)>!</h1></body></html>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := r.Render(tc.input)
			if got != tc.want {
				t.Errorf("Render(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestSafeRenderer(t *testing.T) {
	r := &SafeRenderer{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal name is output",
			input: "Alice",
			want:  "<html><body><h1>Hello, Alice!</h1></body></html>",
		},
		{
			name:  "script tag is escaped",
			input: "<script>alert(1)</script>",
			want:  "<html><body><h1>Hello, &lt;script&gt;alert(1)&lt;/script&gt;!</h1></body></html>",
		},
		{
			name:  "event handler is escaped",
			input: "<img src=x onerror=alert(1)>",
			want:  "<html><body><h1>Hello, &lt;img src=x onerror=alert(1)&gt;!</h1></body></html>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := r.Render(tc.input)
			if got != tc.want {
				t.Errorf("Render(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestDetectXSS(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{input: "Hello, World!", want: false},
		{input: "<script>alert(1)</script>", want: true},
		{input: "<img src=x onerror=alert(1)>", want: true},
		{input: "javascript:alert(1)", want: true},
		{input: "<p onclick=\"evil()\">click</p>", want: true},
		{input: "<b>bold</b>", want: false},
		{input: "safe input", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := detectXSS(tc.input)
			if got != tc.want {
				t.Errorf("detectXSS(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
