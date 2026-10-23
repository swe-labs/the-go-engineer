# Module 07 Checkpoint — CLI, Files, JSON, and Config

## Purpose

This checkpoint verifies that you understand Go's I/O ecosystem — CLI argument parsing, file operations, JSON serialization, and configuration patterns — before advancing to HTTP and web APIs. You should be able to build complete command-line tools that read, process, and write data safely.

## Scope

This checkpoint covers:

- CLI flag parsing with the `flag` package
- `os.Args` for positional arguments
- File reading with `os.ReadFile` and `os.Open`
- File writing with `os.WriteFile` and `os.Create`
- Directory operations with `os.MkdirAll` and `os.Remove`
- `filepath.WalkDir` for recursive directory traversal
- JSON marshaling and unmarshaling with `encoding/json`
- Struct tags for JSON field mapping
- Streaming I/O with `bufio.Scanner` and `bufio.Writer`
- `io.Reader` and `io.Writer` interfaces
- Environment variables with `os.LookupEnv`
- Configuration precedence hierarchies
- Error handling for filesystem operations

## Time expectation

Plan for 45 to 60 minutes.

## Allowed resources

You may use:

- Module 07 lesson READMEs
- Your terminal and editor
- Your own notes

You may not copy from the answer key before attempting the questions.

## Evidence required

Submit:

1. Answers to `questions.md`
2. A small CLI program that reads a JSON file, extracts a specific field, and prints it — demonstrating `os.ReadFile`, `json.Unmarshal`, and error handling
3. A brief paragraph explaining when to use `os.ReadFile` vs `bufio.Scanner` and the memory implications of each

## Rubric

| Dimension | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------------|----------|----------------|
| **File I/O** | Correctly uses `os.ReadFile` for small files and `bufio.Scanner` for streaming; handles errors and file closing | Uses file I/O but ignores errors or forgets to close files | Cannot read a file in Go |
| **JSON handling** | Marshals and unmarshals with struct tags; handles missing fields, unknown fields, and malformed JSON | Basic JSON works but struct tags not used or error cases ignored | Cannot parse JSON into a struct |
| **CLI design** | Uses `flag` package for typed arguments with defaults and help text | Parses args manually with `os.Args` without validation | No argument handling |
| **Streaming and memory** | Explains O(1) vs O(N) memory; chooses `bufio.Scanner` for large files | Knows both approaches exist but cannot explain when to use each | Unaware of streaming I/O |
| **Configuration patterns** | Explains precedence (defaults < file < env); understands why env vars are preferred in production | Knows configs come from files but not env override pattern | No awareness of configuration hierarchy |

**Passing threshold**: 3/5 average across all dimensions

## Retake policy

If you miss a category, return to the relevant lesson, revise your answer, and write one sentence explaining what changed in your understanding.

## NEXT UP

[Module 08 — HTTP and REST APIs](../../../08-http-rest-apis/README.md)
