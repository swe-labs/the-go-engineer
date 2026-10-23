# Build tags

## When you need this

Go code that depends on platform-specific features (system calls, file paths, network APIs) compiles only on specific operating systems or architectures. Without build tags, platform-specific files must use filename suffixes (_linux.go, _windows.go) or be excluded from cross-platform builds entirely. Build tags solve this by allowing conditional compilation based on OS, architecture, custom build constraints, or Go version.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-05

## Mental model

A build tag is a boolean expression evaluated at compile time. Files whose build tag expression evaluates to false are excluded from the compilation. Common constraints include: GOOS (linux, darwin, windows), GOARCH (amd64, arm64), Go version (go1.22), custom tags set via -tags flag, and cgo. The expression syntax supports &&, ||, and ! operators for complex constraints.

## How it works

Go's go/build package evaluates build constraints during package scanning. The //go:build comment (Go 1.17+) replaces the older // +build syntax. The expression is: //go:build linux && amd64 means the file only compiles on linux/amd64. Multiple //go:build lines are ORed. Filename-based constraints use the pattern *_GOOS.go, *_GOARCH.go, or *_GOOS_GOARCH.go — these are checked before build tag comments.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Build tags")
}
```

## Tradeoffs

- Build tags have zero runtime cost — excluded files do not exist in the compiled binary
- Build tag evaluation happens at compile time, not runtime — there is no conditional branching in the generated code
- Complex build tag expressions (&&, ||, !) are resolved entirely during compilation and add no runtime overhead
- Using too many platform-specific build tags increases build matrix complexity but does not affect runtime performance

## Common mistakes

- Using the old // +build syntax without the new //go:build syntax — Go 1.17+ requires both for migration; Go 1.21+ ignores // +build alone.
- Placing the build tag comment after the package declaration — the build tag must be at the top of the file, before the package statement.
- Forgetting that filenames with _linux.go are automatically constrained to linux — adding a //go:build windows to a _linux.go file creates an impossible constraint (file never compiles).
- Using custom build tags inconsistently across files — defining a feature in some files but not others causes partial compilation.
- Not testing all build tag combinations — a build tag combination with a missing file causes a compilation error that is only caught in CI.

## Practice task

The learner must write a platform-specific implementation using build tags (e.g., a signal handler that uses syscall.SIGUSR1 on Linux and os.Interrupt on Windows), demonstrate that the correct file compiles on each platform, and use a custom build tag to conditionally include debug logging.

## Production notes

Build tags are used in every Go project targeting multiple platforms: system-level libraries (database drivers with CGo), platform-specific implementations (file paths, process signals), testing utilities (integration tests with //go:build integration), and debug/logging features (//go:build debug). Production Go code uses build tags for OS-specific system monitoring, Windows service management, and Linux cgroup integration.

## Related core lessons

- Module 16.05 — Method values
- Module 17.04 — Functional options
