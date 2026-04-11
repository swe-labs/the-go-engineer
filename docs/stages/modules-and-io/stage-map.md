# Modules and IO Stage Map

This stage teaches how Go programs cross package, tool, encoding, and filesystem boundaries.

It takes learners from module reasoning into command-line tools, encoded data flows, and real file
operations that interact with the outside world.

## Stage Flow

1. `modules-and-packages`
   - source: [08-modules-and-packages](../../../08-modules-and-packages/)
   - core job: teach module boundaries, dependency management, versioning, and build-surface reasoning
   - milestone: `MP.3` versioning workshop
2. `io-and-cli`
   - source surfaces:
     - [09-io-and-cli/cli-tools](../../../09-io-and-cli/cli-tools/)
     - [09-io-and-cli/encoding](../../../09-io-and-cli/encoding/)
     - [09-io-and-cli/filesystem](../../../09-io-and-cli/filesystem/)
   - core job: teach CLI flows, encoding workflows, and filesystem tooling patterns
   - milestone backbone: `CL.4`, `EN.6`, `FS.7`

## What Each Part Adds

### `modules-and-packages`

This is where the learner stops treating the Go toolchain as magic and starts understanding
package layout, dependency state, and version boundaries.

### `io-and-cli`

This is where the learner starts writing tools that interact with real inputs and outputs instead
of only in-memory examples.

## Recommended Full-Path Order

1. Finish the Section `08` milestone path first.
2. Complete the `cli-tools` track milestone.
3. Complete the `encoding` track milestone.
4. Complete the `filesystem` track milestone.

## Bridge-Path Reminder

If modules, packages, and standard-library tooling already feel familiar, you can move faster
through the early repetition.
What you should not skip is proof:

- `MP.3`
- `CL.4`
- `EN.6`
- `FS.7`

## Exit Condition

You are ready for `4 Backend Engineering` when you can finish the four stage milestones honestly
and explain how module boundaries, CLI interfaces, encoded data, and filesystem logic connect in a
real tool.
