# Section 09: CLI Tools

## Learning Objectives

Every Go engineer eventually builds command-line tools. Go is especially strong here because the
standard library already gives you argument parsing, environment access, filesystem support, and a
simple single-binary deployment model.

By the end of this section, you should understand:

- raw argument access with `os.Args`
- reading configuration from environment variables
- typed flag parsing with the `flag` package
- exit codes and process behavior
- how to structure subcommands for larger CLIs

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
| --- | --- | --- | --- |
| `os.Args` | Beginner | Medium | Raw argument access |
| `flag` package | Intermediate | High | Typed argument parsing |
| `os.Exit` codes | Intermediate | High | Unix process conventions |
| Subcommands | Advanced | High | Building multi-command CLIs |

## Contents

| Directory | Topic | Level |
| --- | --- | --- |
| `1-args/` | Raw arguments with `os.Args` and environment variables | Beginner |
| `2-flags/` | Typed arguments with the `flag` package | Intermediate |
| `3-subcommands/` | Building multi-command CLIs like `git` | Advanced |
| `4-file-organizer/` | Exercise: organize files by extension | Intermediate |

## How to Run

```bash
go run ./09-io-and-cli/cli-tools/1-args hello world
go run ./09-io-and-cli/cli-tools/2-flags -name="The Go Engineer" -count=3
go run ./09-io-and-cli/cli-tools/3-subcommands greet -name="Gopher"
```

## Exercise: File Organizer (`4-file-organizer`)

Build a CLI tool that organizes files in a directory by extension, with a `--dry-run` flag so you
can inspect the plan safely before moving anything.

```bash
go run ./09-io-and-cli/cli-tools/4-file-organizer/_starter --dir=./my-folder
go run ./09-io-and-cli/cli-tools/4-file-organizer --dir=./my-folder
```

## References

- [Package flag](https://pkg.go.dev/flag)
- [Package os](https://pkg.go.dev/os)

## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| CL.1 | [args](./1-args) | `os.Args` · `os.Getenv` · `os.Exit` codes · safe index access | entry |
| CL.2 | [flags](./2-flags) | `flag.String`/`Int`/`Bool` · `flag.Parse()` · pointer return · `-help` | CL.1 |
| CL.3 | [subcommands](./3-subcommands) | `flag.NewFlagSet` · `os.Args[2:]` slicing · switch routing | CL.1, CL.2 |
