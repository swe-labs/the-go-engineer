package main

import (
	"strings"
	"testing"
)

func TestHandleLSPRequest_Hover(t *testing.T) {
	req := LSPRequest{
		JSONRPC: "2.0",
		ID:      42,
		Method:  "textDocument/hover",
	}

	resp := HandleLSPRequest(req)
	if resp.ID != 42 || resp.JSONRPC != "2.0" {
		t.Errorf("invalid LSP response header: %+v", resp)
	}

	hover, ok := resp.Result.(HoverResult)
	if !ok {
		t.Fatalf("expected HoverResult type in result")
	}
	if !strings.Contains(hover.Contents, "func Println") {
		t.Errorf("hover contents missing function signature: %s", hover.Contents)
	}
}

func TestHandleLSPRequest_Definition(t *testing.T) {
	req := LSPRequest{
		JSONRPC: "2.0",
		ID:      43,
		Method:  "textDocument/definition",
	}

	resp := HandleLSPRequest(req)
	loc, ok := resp.Result.(LocationResult)
	if !ok {
		t.Fatalf("expected LocationResult type in result")
	}
	if !strings.HasSuffix(loc.URI, "print.go") || loc.Line <= 0 {
		t.Errorf("invalid definition location: %+v", loc)
	}
}

func TestHandleLSPRequest_Unsupported(t *testing.T) {
	req := LSPRequest{
		JSONRPC: "2.0",
		ID:      44,
		Method:  "unsupported/customMethod",
	}

	resp := HandleLSPRequest(req)
	if resp.Error == "" {
		t.Errorf("expected error for unsupported method")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
