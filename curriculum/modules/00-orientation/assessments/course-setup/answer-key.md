# Answer Key

## Part 1 — Environment verification

1. **Expected:** `go version go1.26.X <os>/<arch>`. Any Go 1.26+ is acceptable.
2. **Expected:** Two lines like `linux` / `amd64` or `darwin` / `arm64`. Any combination is valid.
3. **Expected:** VS Code with `gopls`, Go extension; or GoLand; or Neovim with `gopls` LSP. At least two extensions/plugins named.
4. **Expected:** "Yes, via `editor.formatOnSave`" or "`gofmt -w <file>` manually." Either is acceptable; the key is the learner knows the tool.

## Part 2 — Repository setup

5. **Expected:** `https://github.com/<username>/the-go-engineer` or similar.
6. **Expected:** `git clone https://github.com/<username>/the-go-engineer.git` (or SSH equivalent).
7. **Expected:** `git remote -v` from the repository root. Output should show `origin` pointing to the learner's fork.
8. **Expected:** Two remotes — `origin` pointing to the learner's fork and optionally `upstream` pointing to `swe-labs/the-go-engineer`.

## Part 3 — Course project

9. **Expected:** Go, Git, and a text editor (VS Code, Vim, or Nano).
10. **Expected:** A `setup.txt` file containing the version output of `go version`, `git version`, and the editor version command (e.g., `code --version`).
11. **Expected:** `go version` — output should show `go version go1.26.X <os>/<arch>`.
12. **Expected:** `git version` — output should show `git version 2.4X.X`.

## Part 4 — Reflection

13. **Expected:** Any honest reflection. Common issues: GOPATH confusion, missing `gopls`, Windows path quoting. Resolution should be concrete.
14. **Expected:** Any honest answer. Common unclear areas: the relationship between v2 and v3 curriculum, elective modules, proof workflow.
