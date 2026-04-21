# 05 Packages & IO

## Mission

This stage teaches how Go code crosses package, module, file, encoding, and command-line
boundaries without turning into navigation chaos.

By the end of this stage, a learner should be able to:

- explain what modules and packages are doing in a real repo
- manage dependencies and versioned imports intentionally
- build small CLI workflows that accept input and produce useful output
- read and write structured data through encoding surfaces
- work with files and directories without hiding I/O boundaries

## Stage Map

| Track | Surface | Core Job |
| --- | --- | --- |
| `PKG.1` | [modules and packages](./01-modules-and-packages) | teach module identity, dependency management, versioning, and build surfaces |
| `IO.1` | [I/O and CLI](./02-io-and-cli) | teach CLI input, encoding, filesystem work, and practical I/O tools |

## Why This Stage Exists Now

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

After this stage, move to [06 Backend & DB](../06-backend-db).
