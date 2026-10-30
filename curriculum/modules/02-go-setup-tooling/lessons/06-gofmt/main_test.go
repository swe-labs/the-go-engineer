package main

import (
	"strings"
	"testing"
)

func TestFormatSourceCode(t *testing.T) {
	// 1. Valid messy Go code
	messy := []byte("package main\n\nfunc Hello( ) {\nreturn\n}\n")
	formatted, err := FormatSourceCode(messy)
	if err != nil {
		t.Fatalf("FormatSourceCode failed: %v", err)
	}

	formattedStr := string(formatted)
	if strings.Contains(formattedStr, "Hello( )") {
		t.Errorf("expected spaces inside parameter list to be normalized")
	}
	if !strings.Contains(formattedStr, "\treturn") {
		t.Errorf("expected tab indentation for return statement")
	}

	// 2. Syntax error code -> should fail gracefully
	broken := []byte("package main\nfunc Broken() {")
	_, err = FormatSourceCode(broken)
	if err == nil {
		t.Errorf("expected syntax error when formatting incomplete Go code")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

