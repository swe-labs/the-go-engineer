package main

import "testing"

func TestParseBool(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{name: "true lowercase", input: "true", want: true, wantErr: false},
		{name: "false lowercase", input: "false", want: false, wantErr: false},
		{name: "TRUE uppercase", input: "TRUE", want: true, wantErr: false},
		{name: "FALSE uppercase", input: "FALSE", want: false, wantErr: false},
		{name: "mixed case True", input: "True", want: true, wantErr: false},
		{name: "numeric 1", input: "1", want: true, wantErr: false},
		{name: "numeric 0", input: "0", want: false, wantErr: false},
		{name: "empty string", input: "", want: false, wantErr: true},
		{name: "invalid word", input: "yes", want: false, wantErr: true},
		{name: "garbage", input: "abc123", want: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBool(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseBool(%q) expected error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseBool(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParseBool(%q) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}
