# Stage 05: I/O and CLI

## Mission

This section teaches you how Go handles command-line input, structured text encoding, and
filesystem work without hiding the underlying I/O model.

By the end of Stage 05, you should be comfortable reading and writing:

- small CLIs that use arguments, flags, and subcommands
- JSON encoding and decoding flows for files, APIs, and config
- filesystem tools that traverse directories, stream content, and keep I/O testable
- runnable utilities that interact with the outside world without turning into fragile scripts

## Beta Stage Ownership

This section belongs to [05 Packages & IO](../docs/stages/05-packages-io.md).

Within the beta public shell, it is the second and final part of that stage:

1. Stage 05 `modules-and-packages`
2. Stage 05 `io-and-cli`

## Who Should Start Here

### Full Path

Start here after completing Stage 05 in order.

### Bridge Path

You can move faster if you already understand:

- packages and import paths
- structs and methods
- slices, maps, and error handling

Even on the bridge path, do not skip the first lesson in any track.
Those entry points establish the standard-library vocabulary the rest of the section assumes.

### Targeted Path

This section is the first multi-track pilot in v2.
You can choose the track that matches your immediate goal:

- CLI track for command-line tooling
- Encoding track for JSON and config work
- Filesystem track for file and directory operations

## Track Map

| Track | Entry | Milestone | Focus |
| --- | --- | --- | --- |
| CLI Tools | [CL.1 args](./cli-tools) | `CL.4` | arguments, flags, subcommands, and a safe file-organizer CLI |
| Encoding | [EN.1 marshalling](./encoding) | `EN.6` | JSON encoding/decoding, streaming, base64, and config loading |
| Filesystem | [FS.1 files](./filesystem) | `FS.7` | file I/O, traversal, embed, I/O composition, and log search |

## Suggested Order

1. Complete the CLI track if you want the fastest path to a runnable tool.
2. Complete the Encoding track if you want file/API payload confidence.
3. Complete the Filesystem track if you want stronger operational tooling patterns.
4. Use `FS.8` as an optional stretch lesson after the filesystem milestone.

## Section Milestones

This pilot section has three live milestone surfaces:

- `CL.4` file organizer
- `EN.6` config parser
- `FS.7` log search tool

`FS.8` remains a stretch lesson that extends the filesystem track with a testing seam based on
`fs.FS`.

If you can complete the three milestone exercises and explain:

- why CLI tools should expose safe defaults like `--dry-run`
- why streaming JSON APIs and config files matters more than memorizing `Marshal` vs `Unmarshal`
- why filesystem tools should separate traversal, reading, and matching logic cleanly

then you are ready to move into web and database work in Stage 06.

## Pilot Role In V2

This live v2 slice keeps the current `05-packages-io/02-io-and-cli` filesystem layout and `CL.*`, `EN.*`, and
`FS.*` ids stable while upgrading the learner-facing structure:

- the section now has one top-level guide
- each track has a clearer mission and milestone
- the milestone exercises have explicit README contracts

That keeps the section useful for current learners while the broader v2 migration continues.

## References

1. [Package os](https://pkg.go.dev/os)
2. [Package io](https://pkg.go.dev/io)
3. [Package encoding/json](https://pkg.go.dev/encoding/json)
4. [Package flag](https://pkg.go.dev/flag)

## Next Step

After you finish the track or milestone you care about here, you have completed the core milestone
path for [05 Packages & IO](../docs/stages/05-packages-io.md).

From there, move to [06 Backend & DB](../docs/stages/06-backend-db.md).
The first source section in that next stage is
[Stage 06: Web and Database](../06-backend-db/01-web-and-database).



