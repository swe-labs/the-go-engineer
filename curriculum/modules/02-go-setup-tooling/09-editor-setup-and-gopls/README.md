# Editor Setup and gopls

## Learning objective

By the end of this lesson, you will have `gopls` installed, your editor configured for Go development, and you will use editor features like completion, hover documentation, and go-to-definition.

## Why this matters

A text editor without a language server is just a typewriter. `gopls` transforms your editor into a full-featured Go IDE with real-time feedback — errors underlined as you type, completions when you press `.`, documentation on hover, and instant navigation to definitions. This real-time feedback accelerates learning by showing you mistakes immediately.

## Mental model

`gopls` is the brain behind your Go editor. It reads your Go source files, builds an in-memory model of your codebase, and responds to editor requests (hover, go-to-definition, completion) via the Language Server Protocol (LSP). Your editor is the display; `gopls` is the intelligence.

## Core idea

`gopls` (Go Language Server) provides IDE features for any editor that supports the Language Server Protocol. It is developed by the Go team and is the official Go language server. Install it once, and every LSP-compatible editor gets the same high-quality Go experience.

## Under the hood

`gopls` implements the Language Server Protocol (LSP) — a JSON-RPC protocol for editor communication. It parses the module's Go files, resolves imports, builds a type-checked AST, and maintains a cache of symbols and their locations. When you hover over a symbol, `gopls` looks up its doc comment; when you go-to-definition, it returns the file, line, and column of the declaration. `gopls` also runs `gofmt` on save and provides diagnostics (errors and warnings) as you type.

## How Go uses it

The Go team develops `gopls` as the official language server. It powers all Go editor integrations. Key features:

- **Completion**: suggests symbols, methods, and keywords as you type
- **Hover documentation**: shows the doc comment for any symbol
- **Go-to-definition**: navigates to the declaration of any symbol
- **References**: finds all usages of a symbol
- **Rename**: renames a symbol across all files
- **Diagnostics**: shows compiler errors and vet warnings inline
- **Formatting**: runs `gofmt` on save

Installation: `go install golang.org/x/tools/gopls@latest`

## Go example

```go
package main

import "fmt"

func main() {
	fmt.Println("gopls is the Go language server.")
	fmt.Println("It powers your editor's Go features.")
}
```

## Step-by-step execution

1. Check if `gopls` is installed: `gopls version`
2. If not installed, run: `go install golang.org/x/tools/gopls@latest`
3. Install your editor's Go extension:
   - VS Code: install the Go extension from the marketplace
   - Vim/Neovim: install `vim-go` or configure `gopls` via `coc.nvim` or `vim-lsp`
   - Emacs: use `eglot` or `lsp-mode`
4. Open any Go file and verify features:
   - Hover over `fmt.Println` — you should see its documentation
   - Trigger completion by typing `fmt.` — you should see a list of functions
   - Go to definition on `Println` — it should navigate to the standard library source
5. Introduce a deliberate syntax error (remove a closing brace) and observe the red squiggly underline.

## Common mistakes

- **Installing the editor extension without installing `gopls`**: The extension is just a UI — it needs `gopls` as a backend. Install `gopls` first, then the extension.
- **Running an outdated `gopls`**: `gopls` is updated frequently. Run `go install golang.org/x/tools/gopls@latest` periodically to get the latest features and fixes.
- **Assuming the editor catches all errors**: `gopls` shows compiler errors, but it may not catch all runtime issues. Always run `go vet` and `go test` before committing.
- **Disabling `gopls` because it feels slow on first start**: The first load indexes the entire module and takes 2-10 seconds. Subsequent changes are near-instant.

## Debugging walkthrough

**Scenario**: The editor shows `gopls: no packages found` for every file.

1. Ensure the module has a `go.mod` file. Run `go mod init <module-name>` in the module root if it does not.
2. Ensure you opened the workspace root (the directory containing `go.mod`), not a subdirectory.
3. Reload the window after creating `go.mod`.

**Scenario**: `gopls` and `go build` show different errors.

1. Check the `gopls` version: `gopls version`. Update if outdated.
2. Check that your editor's Go version matches your terminal's Go version: `go version` in both.
3. Restart the language server: in VS Code, run `Go: Restart Language Server` from the command palette.

## Production notes

- Most teams pin the `gopls` version in their development tooling configuration.
- VS Code settings should include `"go.useLanguageServer": true` and `"gopls": { "analyses": { "unusedparams": true } }` for additional checks.
- `gopls` memory usage is typically 100-300 MB for a medium-sized module. This is normal.

## Performance implications

- `gopls` indexes the module on first open, taking 2-10 seconds for large modules.
- Incremental updates after the initial index are sub-second.
- `gopls` watches file changes and re-indexes only modified files.

## Practice task

1. Verify `gopls` is installed with `gopls version`.
2. Open the `main.go` file from this lesson in your editor.
3. Hover over `fmt.Println` and read the documentation.
4. Place your cursor on `Println` and use go-to-definition (Ctrl+click or F12).
5. Type `fmt.` and observe the completion list.
6. Add a syntax error and watch the diagnostic appear.

## Tests / verification

```bash
gopls version
```

Expected output: a version string like `golang.org/x/tools/gopls v0.16.0`

```bash
go run .
```

Expected output:
```
gopls is the Go language server.
It powers your editor's Go features.
```

Verify with your editor that hover, completion, and diagnostics all work on the `main.go` file.

## Review questions

1. What is `gopls` and what problem does it solve?
2. What protocol does `gopls` use to communicate with editors?
3. How do you install or update `gopls`?
4. Why might your editor show different diagnostics than `go build`?
5. What should you do if `gopls` shows "no packages found"?

## NEXT UP

[Reading compiler errors](../10-reading-compiler-errors/README.md) — understand and fix Go compiler error messages.
