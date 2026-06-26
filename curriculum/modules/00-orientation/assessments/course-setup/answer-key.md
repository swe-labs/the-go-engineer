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

9. **Expected:** It checks that the Go toolchain is installed, `go version` works, the repository is correctly cloned, and `go build` succeeds.
10. **Expected:** `ok` with no failures. If tests fail, the environment is not correctly set up.
11. **Expected:** The solution uses a more robust checker, handles errors gracefully, and prints formatted output. The exact diff depends on the implementation.
12. **Expected:** A success message indicating the environment passed all checks.

## Part 4 — Reflection

13. **Expected:** Any honest reflection. Common issues: GOPATH confusion, missing `gopls`, Windows path quoting. Resolution should be concrete.
14. **Expected:** Any honest answer. Common unclear areas: the relationship between v2 and v3 curriculum, elective modules, proof workflow.
