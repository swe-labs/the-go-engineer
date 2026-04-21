# Stage 05: I/O and CLI

## Mission

This track teaches how Go handles command-line input, structured text encoding, and filesystem work without hiding the underlying I/O model.

By the end of this track, you should be comfortable reading and writing:

- small CLIs that use arguments, flags, and subcommands
- JSON encoding and decoding flows for files, APIs, and config
- filesystem tools that traverse directories, stream content, and keep I/O testable
- runnable utilities that interact with the outside world without turning into fragile scripts

## Stage Ownership

This track belongs to [05 Packages & IO](../README.md).

## Track Map

| Track | Entry | Milestone | Focus |
| --- | --- | --- | --- |
| CLI Tools | [CL.1 args](./cli-tools) | `CL.4` | arguments, flags, subcommands, and a safe file-organizer CLI |
| Encoding | [EN.1 marshalling](./encoding) | `EN.6` | JSON encoding/decoding, streaming, base64, and config loading |
| Filesystem | [FS.1 files](./filesystem) | `FS.7` | file I/O, traversal, embed, I/O composition, and log search |

## Suggested Order

1. Complete the CLI track if you want the fastest path to a runnable tool.
2. Complete the Encoding track if you want file and API payload confidence.
3. Complete the Filesystem track if you want stronger operational tooling patterns.
4. Use `FS.8` as an optional stretch lesson after the filesystem milestone.

## Track Milestones

This track has three promoted milestone surfaces:

- `CL.4` file organizer
- `EN.6` config parser
- `FS.7` log search tool

`FS.8` remains a stretch lesson that extends the filesystem track with a testing seam based on `fs.FS`.

## Next Step

After you finish the track or milestone you care about here, continue to [06 Backend, APIs & Databases](../../06-backend-db).
