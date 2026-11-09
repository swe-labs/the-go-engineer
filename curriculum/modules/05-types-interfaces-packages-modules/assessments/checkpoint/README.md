# Module 05 Checkpoint — Types, Interfaces, Packages, and Modules

## Purpose

This checkpoint verifies that you understand Go's type system, interface mechanics, package design, and module management before advancing to testing and debugging. You should be able to explain method sets, predict interface satisfaction, trace nil interface traps, and reason about package boundaries.

## Scope

This checkpoint covers:

- Struct definitions, field access, and struct literals
- Named types and type aliases
- Methods: value receivers vs pointer receivers
- Method sets and their effect on interface satisfaction
- Implicit interface satisfaction
- The empty interface (`any`) and type assertions
- Type switches
- Nil interface values vs interfaces holding typed nil pointers
- Struct embedding and method promotion
- Package declarations, imports, and visibility rules
- Go modules: `go.mod`, `go.sum`, and dependency management
- Generics: type parameters, constraints, and instantiation
- Composition over inheritance

## Time expectation

Plan for 45 to 60 minutes.

## Allowed resources

You may use:

- Module 05 lesson READMEs
- Your terminal and editor
- Your own notes

You may not copy from the answer key before attempting the questions.

## Evidence required

Submit:

1. Answers to `questions.md`
2. A code trace: given a struct with a pointer receiver method, predict whether it satisfies a given interface when passed as a value vs a pointer — then verify by running it
3. A short paragraph explaining why Go chose implicit interface satisfaction over explicit `implements` declarations

## Rubric

| Dimension | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------------|----------|----------------|
| **Method sets** | Correctly states that `*T` has methods of both `T` and `*T`, while `T` only has value receiver methods; predicts interface satisfaction | Knows receiver syntax but cannot predict satisfaction rules | Cannot explain method sets |
| **Interface mechanics** | Explains implicit satisfaction, empty interface, type assertions, and type switches accurately | Understands basic interface usage but uncertain about type assertions | Cannot define or use interfaces |
| **Nil interface traps** | Distinguishes `nil` interface from `(*T)(nil)` in an interface; predicts `== nil` check behavior | Knows nil interfaces exist but cannot predict behavior | Unaware of the nil interface trap |
| **Package and module design** | Explains visibility rules (exported vs unexported), `go.mod` purpose, and dependency management | Knows basic import syntax but uncertain about visibility | Cannot explain package boundaries |
| **Generics** | Explains type parameters, constraints, and when to use generics vs interfaces | Knows generics syntax but uncertain about constraints | Cannot write or read generic code |

**Passing threshold**: 3/5 average across all dimensions

## Retake policy

If you miss a category, return to the relevant lesson, revise your answer, and write one sentence explaining what changed in your understanding.

## NEXT UP

[Module 06 — Testing, Debugging, and Refactoring](../../../06-testing-debugging-refactoring/README.md)
