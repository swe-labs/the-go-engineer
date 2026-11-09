package main

import (
	"encoding/json"
	"fmt"
)

// LSPRequest models a JSON-RPC request sent by an editor to the gopls language server.
type LSPRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		URI      string `json:"uri"`
		Position struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"position"`
	} `json:"params"`
}

// LSPResponse models a JSON-RPC response returned by gopls to the editor.
type LSPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HoverResult represents hover documentation sent back to the editor tooltip.
type HoverResult struct {
	Contents string `json:"contents"`
}

// LocationResult represents the target file and coordinates for 'Go to Definition' (F12).
type LocationResult struct {
	URI  string `json:"uri"`
	Line int    `json:"line"`
}

// HandleLSPRequest simulates how gopls answers editor queries over JSON-RPC.
func HandleLSPRequest(req LSPRequest) LSPResponse {
	resp := LSPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "textDocument/hover":
		resp.Result = HoverResult{
			Contents: "func Println(a ...any) (n int, err error)\nPackage fmt: Println formats using default formats.",
		}
	case "textDocument/definition":
		resp.Result = LocationResult{
			URI:  "file:///usr/local/go/src/fmt/print.go",
			Line: 285,
		}
	default:
		resp.Error = fmt.Sprintf("method not supported: %s", req.Method)
	}

	return resp
}

func main() {
	fmt.Println("=== Deconstructing 'gopls' Language Server Protocol (LSP) ===")

	// Simulating an editor requesting hover information on 'fmt.Println'
	hoverReq := LSPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "textDocument/hover",
	}
	hoverReq.Params.URI = "file:///project/main.go"
	hoverReq.Params.Position.Line = 10
	hoverReq.Params.Position.Character = 6

	resp := HandleLSPRequest(hoverReq)
	rawJSON, _ := json.MarshalIndent(resp, "", "  ")

	fmt.Println("1. Editor -> gopls: JSON-RPC request for 'textDocument/hover'")
	fmt.Println("2. gopls parses AST, type-checks symbols, and replies:")
	fmt.Println(string(rawJSON))

	fmt.Println("\nInsight: gopls is a standalone background daemon communicating via JSON-RPC.")
	fmt.Println("It decouples language intelligence from the editor, powering VS Code, Neovim, and Zed equally.")
}
