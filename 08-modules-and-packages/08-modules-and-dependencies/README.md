# Section 8: Modules & Dependency Management

## Learning Objectives

By the end of this section, you will understand:
- How Go modules work and what `go.mod`/`go.sum` contain
- How to add, update, and remove dependencies
- Semantic versioning and the major version import rule
- The `replace` directive for local development
- Multi-module workspaces with `go.work`

## Contents

| Directory | Topic | Level |
|-----------|-------|-------|
| `1-module-basics/` | `go.mod` anatomy, import paths, essential commands | Beginner |
| `2-managing-deps/` | Adding/removing/inspecting dependencies | Intermediate |
| `3-versioning/` | Semantic versioning, `replace`, `exclude`, vendoring | Intermediate |

## How to Run

```bash
go run ./08-modules-and-dependencies/1-module-basics
go run ./08-modules-and-dependencies/2-managing-deps
go run ./08-modules-and-dependencies/3-versioning
```

## Exercises

1. Create a separate Go module in `/tmp/mathlib` with `Add()` and `Multiply()` functions
2. Use `replace` to import it locally into this project
3. Remove the `replace` directive and discuss what would happen if you tried to build

## References

- [Using Go Modules](https://go.dev/blog/using-go-modules)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Module version numbering](https://go.dev/doc/modules/version-numbers)


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| MP.1 | [module basics](./1-module-basics) | go.mod anatomy · module path · go.sum checksums | 🟢 entry |
| MP.2 | [managing deps](./2-managing-deps) | go get · go mod tidy · go list · go mod why | MP.1 |
| MP.3 | [versioning](./3-versioning) | semver · v2 import path rule · replace · exclude · vendoring | MP.1, MP.2 |
