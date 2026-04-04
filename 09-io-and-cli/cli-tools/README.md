# Section 19: CLI Tools

## Learning Objectives

Every Go engineer needs to build command-line tools. Go excels at CLI development — the standard library provides `flag` for argument parsing, and the single-binary deployment model means your CLI runs anywhere without dependencies.

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
| ----- | ----- | ---------- | ------------------- |
| os.Args | Beginner | Medium | Raw argument access |
| flag package | Intermediate | High | Typed argument parsing |
| os.Exit codes | Intermediate | High | Unix process conventions |
| Subcommands | Advanced | High | Building complex CLIs |

## Contents

| Directory | Topic | Level |
| --------- | ----- | ----- |
| `1-args/` | Raw arguments with `os.Args`, environment variables | Beginner |
| `2-flags/` | Typed arguments with the `flag` package | Intermediate |
| `3-subcommands/` | Building multi-command CLIs (like `git`) | Advanced |
| `4-file-organizer/` | **Exercise:** CLI file organizer by extension | Intermediate |

## How to Run

```bash
go run ./19-cli-tools/1-args hello world
go run ./19-cli-tools/2-flags -name="Go Mastery" -count=3
go run ./19-cli-tools/3-subcommands greet -name="Gopher"
```

---

## 🏗 Exercise: File Organizer (`4-file-organizer`)

Build a CLI tool that organizes files in a directory by extension, with a `--dry-run` flag. Try it yourself first!

```bash
go run ./19-cli-tools/4-file-organizer/_starter --dir=./my-folder  # Try the exercise
go run ./19-cli-tools/4-file-organizer --dir=./my-folder           # See the solution
```

## References

- [Go Package: flag](https://pkg.go.dev/flag)
- [Go Package: os](https://pkg.go.dev/os)
