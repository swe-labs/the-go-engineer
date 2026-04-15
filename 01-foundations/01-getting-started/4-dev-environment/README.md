# GT.4 Development Environment

## Mission

Learn the basic command loop that makes day-to-day Go work predictable.

This lesson is the bridge between "I can run one program" and "I know the basic tools I will keep
using across the repo."

## Why This Lesson Exists Now

A beginner who can run one lesson still needs one more confidence layer:

- how to format code
- how to compile code
- how to run tests
- how to recognize whether editor support is installed

That is what this lesson establishes.

## Mental Model

The Go toolchain is a working loop, not a single command.

The beginner-safe version of that loop is:

1. write or change code
2. run `go fmt`
3. run `go build` or `go run`
4. run `go test` when tests exist

That loop repeats through the whole curriculum.

## Visual Model

```text
edit code
   |
   +--> go fmt
   +--> go build or go run
   +--> go test
```

```text
editor support:

gopls -> code intelligence
gofmt -> standard formatting
```

## Machine View

These commands do different jobs:

- `go fmt` rewrites source files into Go's standard style
- `go build` compiles packages to verify they are valid
- `go run` compiles and executes in one step
- `go test` builds and runs tests

This lesson also checks whether some tools can be found on the machine path.
That is why it uses `exec.LookPath(...)`.

## Run Instructions

```bash
go run ./01-foundations/01-getting-started/4-dev-environment
```

## Code Walkthrough

### `commands := []commandInfo{ ... }`

The lesson stores the important Go commands in a small slice of labeled records.
That keeps the output readable without repeating nearly identical `fmt.Printf(...)` blocks.

### `for _, command := range commands { ... }`

This loop prints the command list in a consistent format.
The learner does not need to master loops here.
They only need to recognize that repetition can print structured output cleanly.

### `tools := []toolInfo{ ... }`

This second list describes the editor-support tools the lesson wants to check.

### `exec.LookPath(tool.name)`

This asks the operating system whether a tool is available on the command path.

If the tool is found, the lesson prints where it lives.
If it is missing, the lesson prints a helpful note instead.

### `fmt.Println("NEXT UP: LB.1 variables")`

The footer connects this section to the next learning surface.
It tells the learner they are leaving setup and moving toward language fundamentals.

## Try It

1. Run `go fmt ./...` in the repo root after reading this lesson.
2. Run `go build ./...` in the repo root and see whether the command finishes quietly.
3. Temporarily add a fake tool name to the tool list and inspect the "not found" branch.

## Common Questions

- Why is `go fmt` such a big deal in Go?
  Because Go intentionally uses one standard formatting style so teams do not waste energy arguing
  about formatting rules.

- Why can `go build` finish with no output?
  In Go, a successful build often says nothing. Silence usually means success.

## Next Step

Continue to `LB.1` variables once `02-language-basics` is rebuilt, or move through the currently
available canonical foundations path from `03-control-flow` onward.
