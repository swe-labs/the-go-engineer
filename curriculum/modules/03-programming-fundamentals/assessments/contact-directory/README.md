# Project Assessment — Contact Directory

## Purpose

Evaluate your ability to implement CRUD operations using Go maps, slices, sorting, and error handling. This assessment verifies that you can build a practical data management tool and prove its correctness through automated tests and edge-case reasoning.

## Targets

- `project-contact-directory`

## Scope

This assessment covers:

- Map declaration, initialization, and key-based lookups
- The comma-ok idiom for safe map access
- Slice operations including sorting with `sort.Slice`
- String manipulation with `strings.TrimSpace` and `strings.Contains`
- Case-insensitive comparison with `strings.ToLower`
- Custom struct types
- Error creation and return patterns
- Edge cases: duplicate keys, missing entries, empty inputs

## Evidence required

Submit:

1. Your completed `_starter/main.go` with all 7 function bodies implemented
2. Output of `go test -v .` from the `_starter/` directory showing all tests passing
3. Output of `go run .` showing a complete directory workflow (add, list, search, update, delete)
4. Output of `gofmt -l .` (must be empty) and `go vet .` (must be clean)

## Time expectation

Plan for 60 to 90 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Correctness** | 30% | All 7 CRUD operations work correctly; all tests pass including edge cases (not found, duplicate phone, empty fields) | 5–6 operations work; most tests pass | Fewer than 4 operations work |
| **Data handling** | 25% | Proper error handling for missing, duplicate, and empty inputs; comma-ok idiom used for map lookups | Basic error handling with gaps in edge cases | Missing or incorrect error handling; panics on missing keys |
| **Proof** | 25% | Tests cover normal cases and edge cases (not found, duplicate, empty, case-insensitive search, sorted output verification) | Tests cover normal CRUD operations only | No meaningful test coverage |
| **Code quality** | 20% | Idiomatic Go: clean map/slice usage, proper error returns, `sort.Slice` for ordering, consistent naming | Mostly readable with minor issues | Hard to read, unidiomatic, or uses global state |

**Passing score**: 70%

## Review questions

Be prepared to answer these questions during a review session:

1. What happens if you access a map key that does not exist without using the comma-ok idiom? Why is this dangerous for a contact directory?
2. Your `ListContacts` function returns contacts sorted by name. Explain how `sort.Slice` works — what does the `less` function receive and return?
3. Walk through your `DeleteContact` function: what Go built-in do you use to remove a map entry? What happens if you delete a key that does not exist?
4. Your `SearchContacts` uses case-insensitive matching. Why do you convert both the query and the contact fields to lowercase before comparing, instead of using `==`?
5. What is the zero value of a `map[string]Contact`? What happens if you try to write to it before calling `make()`?

## Passing standard

You pass when all tests pass, your directory workflow output demonstrates all CRUD operations, and you can explain map and slice behavior for edge cases without guessing.

## Retake policy

If you score below 70%, review the relevant Module 03 lessons (maps, slices, comma-ok idiom, sorting), fix the issues, and resubmit with a short correction note.

## NEXT UP

[Module 03 Checkpoint](../checkpoint/README.md)
