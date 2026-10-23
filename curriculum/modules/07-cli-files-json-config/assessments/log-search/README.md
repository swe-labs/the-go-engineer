# Project Assessment — Log Search

## Purpose

Evaluate your ability to build a memory-efficient streaming search tool using `bufio.Scanner`, regular expressions, log-level hierarchy filtering, and contextual output buffering. This assessment verifies that you can process arbitrarily large files in constant memory — a critical production skill.

## Targets

- `project-log-search`

## Scope

This assessment covers:

- Streaming I/O with `bufio.Scanner` for O(1) memory processing
- Regular expression compilation and matching with `regexp`
- Case-sensitive and case-insensitive string matching
- Log-level severity hierarchies (DEBUG < INFO < WARN < ERROR)
- Ring buffer implementation for context-before lines
- After-context countdown tracking
- CLI flag parsing with the `flag` package
- Error handling for invalid regex patterns

## Evidence required

Submit:

1. Your completed `_starter/main.go` with `SearchStream`, `extractLogLevel`, and context buffering implemented
2. Output of `go test -v .` showing all tests passing
3. A brief explanation of why you use `bufio.Scanner` instead of `os.ReadFile` — what happens to memory when processing a 10GB log file with each approach?

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Streaming and memory** | 30% | Line-by-line processing via `bufio.Scanner`; memory remains constant O(1) regardless of input size; scanner buffer expanded for long lines | Scanner used with reasonable memory footprint | Entire file loaded into memory with `os.ReadFile` |
| **Pattern and regex** | 25% | Substring and regex searches handled correctly; case sensitivity toggle works; `ErrInvalidRegex` returned for bad patterns | Searches work but case sensitivity buggy or regex errors panic | Regex compilation errors cause panics |
| **Context window** | 25% | Before-context ring buffer correctly sized and flushed; after-context countdown accurate; match markers (`>` vs ` `) clean | Context mostly works with off-by-one line counts | Context lines missed, duplicated, or not implemented |
| **Severity hierarchy** | 20% | Log levels correctly ranked; `WARN` filter shows WARN + ERROR but not INFO + DEBUG; level extraction handles multiple formats | Levels checked but hierarchy ordering is incorrect | Misclassifies log severity or no filtering |

**Passing score**: 75%

## Review questions

1. What is the memory complexity of `bufio.Scanner` vs `os.ReadFile` for a file of size N? Why does this matter for production log analysis?
2. Your ring buffer stores context-before lines. What happens when a new line is added and the buffer is full? Why is this better than an unbounded slice?
3. Walk through your log-level hierarchy: if a user specifies `-level WARN`, which levels are included and which are excluded? How do you compare levels?
4. What happens when a user passes an invalid regex like `[unclosed`? Show the error path from `regexp.Compile` to the user.
5. How does your scanner handle lines longer than the default 64KB buffer? What does `scanner.Buffer()` do?

## Passing standard

You pass when all tests pass, your tool processes large inputs in constant memory, and you can explain streaming I/O and ring buffer mechanics.

## Retake policy

If you score below 75%, review Module 07 lessons on streaming I/O, `bufio.Scanner`, and ring buffers, then resubmit.

## NEXT UP

[Project Assessment: File Organizer](../file-organizer/README.md)
