# Track A: CLI Tools

## Mission

This track teaches you how to build small command-line programs that accept input clearly, expose
safe options, and scale from one command to multiple subcommands.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `CL.1` | Lesson | [args](./1-args) | Introduces raw command-line arguments, environment variables, and exit codes. | entry |
| `CL.2` | Lesson | [flags](./2-flags) | Shows typed option parsing with the `flag` package. | `CL.1` |
| `CL.3` | Lesson | [subcommands](./3-subcommands) | Builds multi-command CLIs without third-party frameworks. | `CL.1`, `CL.2` |
| `CL.4` | Exercise | [file organizer](./4-file-organizer) | Combines flags, directory reading, and safe filesystem changes in one milestone. | `CL.1`, `CL.2`, `CL.3` |

## Suggested Order

1. Work through `CL.1` to `CL.3` in order.
2. Complete `CL.4` without copying the finished solution line by line.

## Track Milestone

`CL.4` is the current CLI track milestone.

If you can complete it and explain:

- why `flag.Parse()` gives you safer CLI input than manual `os.Args` slicing
- why `--dry-run` is a core safety feature for filesystem-changing tools
- why subcommands are a better scaling pattern than one giant `main()`

then the CLI part of Section 09 is doing its job.

## Next Step

After `CL.4`, continue to the [Encoding track](../encoding) or back to the
[Section 09 overview](../README.md).
