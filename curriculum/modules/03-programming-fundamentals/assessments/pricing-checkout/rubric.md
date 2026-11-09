# Rubric

Grade using the assessment criteria in metadata.

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------------|----------|----------------|
| **Correctness (35%)** | All pricing functions work correctly; tests cover normal and edge cases; output matches spec | Most functions work; some edge cases missed | Functions incomplete or incorrect |
| **Go quality (25%)** | Idiomatic Go with proper error handling, clean switch/loop usage, readable formatting | Mostly readable with minor issues | Hard to read or unidiomatic |
| **Proof (25%)** | Test file covers normal, empty, negative, bulk discount, and unknown country cases | Tests cover normal cases only | No meaningful tests |
| **Explanation (15%)** | Clearly explains discount-vs-tax ordering, cents-vs-float choice, country code handling, and edge cases | Basic explanation with gaps | Unable to explain design decisions |

**Passing score:** 80%

**Evidence required:**
- Working Pricing Checkout implementation
- Test output showing all tests pass
- Written answers to review questions
