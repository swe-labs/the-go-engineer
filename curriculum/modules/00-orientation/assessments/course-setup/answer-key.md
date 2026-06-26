# Answer Key

## Part 1 — Environment verification

1. **Expected:** `go version go1.26.X <os>/<arch>`. Any Go 1.26+ is acceptable.
2. **Expected:** Two lines like `linux` / `amd64` or `darwin` / `arm64`. Any combination is valid.
3. **Expected:** VS Code with `gopls`, Go extension; or GoLand; or Neovim with `gopls` LSP. At least two extensions/plugins named.
4. **Expected:** "Yes, via `editor.formatOnSave`" or "`gofmt -w <file>` manually." Either is acceptable; the key is the learner knows the tool.

## Part 2 — Repository setup

5. **Expected:** `https://github.com/<username>/the-go-engineer` or similar.
6. **Expected:** `git clone https://github.com/<username>/the-go-engineer.git` (or SSH equivalent).
7. **Expected:** `go test ./...` from the repository root. Output should show `ok` lines and no `FAIL`.
8. **Expected:** Two remotes — `origin` pointing to the learner's fork and optionally `upstream` pointing to `swe-labs/the-go-engineer`.

## Part 3 — Course project

9. **Expected:** It uses `exec.LookPath` to check that Go and Git are installed and available in PATH. An editor check is left as a TODO.
10. **Expected:** `ok` with no failures. If tests fail, the environment is not correctly set up.
11. **Expected:** The solution checks five tools (Go, Git, VS Code, Vim, Nano), uses a loop instead of repeated calls, and prints a summary message. The starter checks only Go and Git with a TODO for the editor.
12. **Expected:** Two lines: `OK: Go is installed` and `OK: Git is installed`, preceded by a title and separator.

## Part 4 — Reflection

13. **Expected:** Any honest reflection. Common issues: GOPATH confusion, missing `gopls`, Windows path quoting. Resolution should be concrete.
14. **Expected:** Any honest answer. Common unclear areas: the relationship between v2 and v3 curriculum, elective modules, proof workflow.
