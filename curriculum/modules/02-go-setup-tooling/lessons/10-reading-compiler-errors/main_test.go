package main

import (
	"testing"
)

func TestParseCompilerError(t *testing.T) {
	tests := []struct {
		raw          string
		wantFile     string
		wantLine     int
		wantCol      int
		wantCategory string
		wantErr      bool
	}{
		{
			raw:          `server.go:24:9: "fmt" imported and not used`,
			wantFile:     "server.go",
			wantLine:     24,
			wantCol:      9,
			wantCategory: "Unused Import",
			wantErr:      false,
		},
		{
			raw:          `math.go:10:5: count declared and not used`,
			wantFile:     "math.go",
			wantLine:     10,
			wantCol:      5,
			wantCategory: "Unused Variable",
			wantErr:      false,
		},
		{
			raw:          `app.go:50:12: cannot use 100 as string`,
			wantFile:     "app.go",
			wantLine:     50,
			wantCol:      12,
			wantCategory: "Type Mismatch",
			wantErr:      false,
		},
		{
			raw:     `malformed string`,
			wantErr: true,
		},
		{
			raw:     `main.go:notanumber:5: bad line`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		got, err := ParseCompilerError(tt.raw)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseCompilerError(%q) error = %v, wantErr %v", tt.raw, err, tt.wantErr)
			continue
		}
		if !tt.wantErr {
			if got.File != tt.wantFile || got.Line != tt.wantLine || got.Col != tt.wantCol || got.Category != tt.wantCategory {
				t.Errorf("ParseCompilerError(%q) = %+v, want file=%s line=%d col=%d cat=%s",
					tt.raw, got, tt.wantFile, tt.wantLine, tt.wantCol, tt.wantCategory)
			}
		}
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

