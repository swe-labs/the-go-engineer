# CLI, Files, JSON, Config, and Runtime Coordination

## Goal

Build useful local programs and understand process lifecycle before servers.

## Summary

Teaches practical local tools using CLI input, file I/O, JSON, configuration, signals, and runtime coordination before network services.

## Learning Sessions & Lessons

### Session 1: Command-Line Interfaces & Unix Streams
- `core-07-01` CLI args
- `core-07-02` Flags
- `core-07-03` Subcommands
- `core-07-04` Standard input and standard output
- `core-07-05` io.Reader
- `core-07-06` io.Writer

### Session 2: Filesystem Operations & Abstractions
- `core-07-07` Files
- `core-07-08` Paths
- `core-07-09` Directories
- `core-07-10` Temp files
- `core-07-11` fs.FS

### Session 3: JSON Data Serialization (Consolidated Workshop)
*Workshop A: In-Memory JSON Values (Working with byte slices and structs)*
- `core-07-12` **JSON marshal** *(Serializing Go structs and values to JSON bytes)*
- `core-07-13` **JSON unmarshal** *(Parsing JSON bytes into concrete Go structures)*
*Workshop B: Stream-Based JSON Processing (Working with `io.Reader` and `io.Writer` streams)*
- `core-07-14` **JSON encoder** *(Streaming objects directly to files, pipes, and network sockets)*
- `core-07-15` **JSON decoder** *(Streaming large, multi-token JSON payloads without memory bloat)*

### Session 4: Application Configuration & Templates
- `core-07-16` Config files
- `core-07-17` Environment variables for configuration
- `core-07-18` Config validation
- `core-07-19` Text templates
- `core-07-20` Base64 as encoding, not encryption
- `core-07-21` CLI testability

## Completion

Complete every required lesson, lab, project, and assessment for this module.
