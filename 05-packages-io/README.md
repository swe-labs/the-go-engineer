# 05 Packages & IO

## Mission

This section teaches how Go code crosses package, module, file, encoding, and command-line
boundaries without turning into navigation chaos.

By the end of this section, a learner should be able to:

- explain what modules and packages are doing in a real repo
- manage dependencies and versioned imports intentionally
- build small CLI workflows that accept input and produce useful output
- read and write structured data through encoding surfaces
- work with files and directories without hiding I/O boundaries

## Section Map

| Track | Surface | Core Job |
| --- | --- | --- |
| `MP.1-MP.4` | [Modules and Packages](./01-modules-and-packages) | teach module identity, dependency management, versioning, and package boundaries |
| `CL.1-CL.4` | [CLI Workflows](./02-io-and-cli/cli-tools) | teach argument parsing, flags, standard streams, and command behavior |
| `EN.1-EN.6` | [Encoding](./02-io-and-cli/encoding) | teach JSON, CSV, and structured data movement |
| `FS.1-FS.8` | [Filesystem](./02-io-and-cli/filesystem) | teach files, directories, paths, and practical file tooling |

## Why This Section Exists Now

The learner already knows:

- values, control flow, and data structures
- function boundaries and explicit errors
- structs, interfaces, composition, and text-heavy modeling

That is enough to start asking engineering questions like:

- how does code stay organized across files and packages?
- how does a program accept real input from users or files?
- how does structured data move in and out of a process?

## Suggested Learning Flow

1. Start with [modules and packages](./01-modules-and-packages).
2. Move into [I/O and CLI](./02-io-and-cli) after module identity feels stable.
3. Complete the milestone surfaces in both tracks.

## Next Step

After this section, move to [06 Backend & DB](../06-backend-db).
