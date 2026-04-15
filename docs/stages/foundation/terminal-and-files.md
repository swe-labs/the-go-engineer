# Terminal And Files Mental Model

If the terminal still feels mysterious, use this page before you assume you are "bad at coding."

## Three Things To Keep Straight

### 1. The terminal is a place where you give commands

You type a command like:

```bash
go run ./01-foundations/01-getting-started/2-hello-world
```

That tells the system to run the code inside that lesson folder.

### 2. Files hold code

A `.go` file is a source file.
It is not running by itself.
The Go tool reads it and turns it into execution.

### 3. Folders organize related code

In this repo, lessons live inside folders.
When you see a path like:

```text
01-foundations/01-getting-started/2-hello-world
```

that path is pointing you to one lesson surface.

## What `go run` Is Doing

At a beginner level, think of `go run` like this:

1. find the Go code at the path you gave it
2. build a runnable version temporarily
3. execute it
4. show you the result

You do not need the deep compiler details yet.
You only need the mental model that Go is turning files into a running program.

## What Usually Confuses Beginners

- being in the wrong folder
- copying the wrong path
- expecting a file to "run itself"
- thinking terminal errors mean permanent failure

## Beginner Rule

When something fails, ask:

1. Am I in the repo root?
2. Did I copy the path correctly?
3. Does `go version` still work?
4. Am I reading a normal tool error, or am I panicking too early?

Most early problems are simple path or environment problems, not deep programming failure.
