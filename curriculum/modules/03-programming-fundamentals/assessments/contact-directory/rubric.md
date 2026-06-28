# Rubric

Grade using the assessment criteria in metadata.

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------------|----------|----------------|
| **Correctness (35%)** | All CRUD operations work correctly; tests pass for normal cases, duplicates, not-found, empty directory, and search | Most operations work; some edge cases missed | Operations incomplete or incorrect |
| **Go quality (25%)** | Idiomatic Go with proper map/slice usage, clean error handling, readable sorting and search logic | Mostly readable with minor issues | Hard to read or unidiomatic |
| **Proof (25%)** | Tests cover add, get, update, delete, list, search; empty and not-found edge cases; case-insensitive search | Tests cover normal operations only | No meaningful tests |
| **Explanation (15%)** | Clearly explains map key choice, sorting approach, whitespace handling, search trade-offs, and add-vs-update semantics | Basic explanation with gaps | Unable to explain design decisions |

**Passing score:** 80%

**Evidence required:**
- Working Contact Directory implementation
- Test output showing all tests pass
- Written answers to review questions
