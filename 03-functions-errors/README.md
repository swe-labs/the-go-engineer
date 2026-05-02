# Section 03: Functions & Errors

## Mission

Learn how to model behavior boundaries and failure states by treating both functions and errors as first-class, manageable values in Go.

By the end of this section, a learner should be able to:
- Define function signatures with multiple returns and named return values.
- Handle errors explicitly as values instead of using exceptions.
- Use `defer` for deterministic resource cleanup.
- Implement first-class functions and state-capturing closures.
- Understand the mechanics of panic and recover boundaries.

## Section Map

### Track 1: Function Basics (FE)
| ID | Lesson | What It Unlocks |
| --- | --- | --- |
| `FE.1` | [Functions Basics](./1-functions-basics/) | Encapsulating logic into reusable boundaries. |
| `FE.2` | [Parameters and Returns](./2-parameters-and-returns/) | Passing data in and getting results out. |
| `FE.3` | [Multiple Return Values](./3-multiple-return-values/) | Go's idiomatic way to return data and error simultaneously. |
| `FE.4` | [Errors as Values](./4-errors-as-values/) | Why Go doesn't use try/catch and how to handle errors. |
| `FE.5` | [Validation](./5-validation/) | Guard clauses and defensive programming. |
| `FE.6` | [Orchestration](./6-orchestration/) | Coordinating multiple functions to solve a workflow. |

### Track 2: Advanced Functions (FE)
| ID | Lesson | What It Unlocks |
| --- | --- | --- |
| `FE.8` | [First-Class Functions](./7-first-class-functions/) | Treating functions as data and callbacks. |
| `FE.9` | [Closures Mechanics](./8-closures-mechanics/) | Functions that "carry" their environment with them. |
| `FE.7` | [Order Summary](./9-order-summary/) | **Milestone**: Building a multi-step pricing and tax engine. |
| `FE.10` | [Panic and Recover](./10-panic-and-recover/) | Handling unrecoverable failures and stopping crashes. |

## Zero-Magic Boundary

This section focuses on **Logical Boundaries**.
It does **not** formally teach:
- Struct methods or Receivers (covered in Section 04).
- Interface-based polymorphism (covered in Section 04).
- Package-level visibility and internal/external boundaries (covered in Section 05).

> [!NOTE]
> Functions are the fundamental unit of work in Go. Mastering them is a prerequisite for [Section 04: Type Design](../04-types-design/README.md), where you will learn to attach this behavior to custom data types.

## Next Step

After completing this section, continue to [Section 04: Type Design](../04-types-design/README.md).
