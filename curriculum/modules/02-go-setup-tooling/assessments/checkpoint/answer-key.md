# Module 02 Checkpoint Answer Key

Use this after attempting the checkpoint.

## Part 1

1. It proves the shell can find a Go toolchain and shows the installed version.
2. It proves which `go` executable the shell will run.
3. `GOROOT` is the Go installation root.
4. `GOPATH` is a workspace/cache-era path still used for defaults such as module cache placement.
5. `GOMODCACHE` stores downloaded modules; `GOCACHE` stores build/test cache artifacts.

## Part 2

6. Package declaration, import, `main` function, and output call.
7. It resolves packages and compiles a temporary executable.
8. Syntax/type/import/module errors can stop compilation first.
9. It produces an executable artifact.
10. To choose a predictable output path/name.

## Part 3

11. It compiles package code plus `_test.go` files into a temporary test binary.
12. It states actual vs expected behavior with enough context to debug.
13. It is the Go standard format and removes style debates.
14. It lists files whose formatting differs from gofmt output.

## Part 4

15. Suspicious patterns such as malformed format strings or likely incorrect constructs.
16. Vet is static analysis; it does not prove all behavior.
17. Documentation comments and exported declarations from packages.
18. Examples: `go doc fmt.Println`, `go doc testing.T`, `go doc errors`.

## Part 5

19. `gopls` is the Go language server that powers editor diagnostics, completion, navigation, rename, and quick fixes.
20. Terminal commands are the same evidence CI and release validation use.

## Part 6

21. Read file, line, column, message, then inspect the first reported error.
22. A compiler error happens before a program can run; a runtime error happens while it is executing.
23. Test name, expected value, actual value, assertion message, and likely source of disagreement.
24. Later errors may be consequences of the first one.

## Part 7

25. `go.mod` declares module path, Go version, and dependencies.
26. It shows the path to the active module file, or an empty/dev-null-like value depending on mode.
27. The Go tool cannot determine imports, module path, or dependency context the same way.
28. Errors about packages not found, wrong import paths, or `go.mod file not found`.

## Part 8

29. Strong checklists include `pwd`, `go version`, `go env GOMOD`, `gofmt -l`, `go test ./...`, `go vet ./...`, and project validators.
30. Good answers name one weak command and a practice plan.
