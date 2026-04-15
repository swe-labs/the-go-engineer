# GT.1 Installation Verification

## Mission

Confirm that Go is installed and that this machine can actually run a Go program.

This lesson is intentionally simple.
The point is not to impress the learner.
The point is to prove that the environment is real.

## Why This Lesson Exists Now

Beginners often lose confidence before programming even starts.
They are not stuck on logic yet.
They are stuck on setup.

So the first lesson should answer one clear question:

Can this computer run Go code successfully?

## Mental Model

Running a Go lesson means:

1. the `go` tool reads the source files
2. Go compiles them
3. the compiled program runs
4. the terminal shows the output

If this lesson runs, the whole learning loop becomes much more trustworthy.

## Visual Model

```text
you type:
go run ./01-foundations/01-getting-started/1-installation
```

```text
Go tool:
source code -> compile -> execute -> terminal output
```

```text
successful output means:
- Go is installed
- the repo path is correct
- the terminal can run the toolchain
```

## Machine View

This lesson calls values from Go's `runtime` package.
That package exposes information about the program that is currently running.

When the lesson prints:

- Go version
- operating system
- architecture
- CPU count

it is reading facts from the running binary and the current machine environment.

## Run Instructions

```bash
go run ./01-foundations/01-getting-started/1-installation
```

## Code Walkthrough

### `package main`

This file belongs to the special `main` package.
That tells Go this code should build into an executable program, not a reusable library.

### `import ("fmt" "runtime")`

The program needs two standard-library packages:

- `fmt` to print output
- `runtime` to inspect the running program and machine details

### `fmt.Println("Go installation looks healthy.")`

This prints the first human-facing success message.
It gives the learner a fast confidence signal before any detail lines appear.

### `runtime.Version()`

This returns the Go version that built the running program.
It answers: "Which Go toolchain is active right now?"

### `runtime.GOOS` and `runtime.GOARCH`

These identify the operating system and processor architecture.

Examples:

- `windows/amd64`
- `linux/amd64`
- `darwin/arm64`

### `runtime.NumCPU()`

This shows how many logical CPUs the program can see.
The beginner does not need concurrency yet.
They only need to see that a running program can inspect its environment.

### `fmt.Println("NEXT UP: GT.2 hello-world")`

The footer keeps the lesson path explicit.
A beginner should never have to guess where to go next.

## Try It

1. Run `go version` in the terminal before or after this lesson and compare it to the program output.
2. Change the first message text and rerun the lesson.
3. Add one more `fmt.Println(...)` line and confirm the program still runs.

## Common Questions

- Why does a setup lesson need code?
  Because the goal is not only to install Go. The goal is to prove the machine can run a real Go
  program from this repo.

- Why print system facts this early?
  Because it helps the learner connect "program output" with "real machine state."

## Next Step

Continue to `GT.2` hello world.
