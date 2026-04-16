# Modules and IO Milestone Guidance

The `05 Packages & IO` stage has four proof surfaces.

These milestones show whether a learner can work across module, process, encoding, and filesystem
boundaries instead of staying inside toy in-memory examples.

## Milestone Backbone

| ID | Surface | What It Proves |
| --- | --- | --- |
| `MP.3` | versioning workshop | you understand module versioning, `/v2` import paths, and local `replace` reasoning |
| `CL.4` | file organizer | you can build a safe CLI workflow with arguments, flags, and clear output |
| `EN.6` | config parser | you can handle encoded configuration data through real JSON/config workflows |
| `FS.7` | log search tool | you can traverse files and directories without tangling search, reading, and matching logic |

## How To Use These Milestones

### Full Path

Complete all four in order.
They are the strongest proof that the stage really stuck.

### Bridge Path

If you move faster through the lessons, use these milestones as the non-skippable proof surfaces.
Do not claim the stage is done if you only skimmed the lessons and skipped the milestone work.

### Targeted Path

If you enter the stage late, choose the milestone that best matches your real boundary gap first,
then check backward honestly:

- weak on modules and dependency behavior: start with `MP.3`
- weak on CLI tooling: start with `CL.4`
- weak on encoding and config flows: start with `EN.6`
- weak on files and directory traversal: start with `FS.7`

## Honest Readiness Signals

You are likely ready to leave this stage when:

- you can complete the milestone without copying the solution line by line
- you can explain why the tool or workflow is shaped that way
- you can change an input, path, or output shape without the whole program collapsing
- you can describe which boundary the code is protecting and why

## If A Milestone Feels Too Hard

That is usually a routing signal, not a failure signal.

Go back one layer:

- `MP.3` too hard: revisit `go.mod`, dependency commands, and module versioning rules
- `CL.4` too hard: revisit args, flags, and small command dispatch
- `EN.6` too hard: revisit marshalling, unmarshalling, and config-loading patterns
- `FS.7` too hard: revisit file basics, directory traversal, and I/O composition
