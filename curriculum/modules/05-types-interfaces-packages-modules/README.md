# Types, Interfaces, Packages, and Modules

## Goal

Turn scripts into maintainable Go programs.

## Summary

Moves from scripts to maintainable Go systems by teaching types, methods, interfaces, packages, modules, and dependency boundaries.

## Prerequisites

- [Module 04: Functions, Errors, and Data Semantics](../04-functions-errors-data-semantics/README.md) — complete that module's completion checkpoint before starting this one.

## Learning Sessions & Lessons

### Session 1: Structs & Composition
- `core-05-01` Structs
- `core-05-02` Methods
- `core-05-06` Composition
- `core-05-07` Embedding

### Session 2: Method Receivers & Mutation Semantics (Consolidated Workshop)
*A unified three-part workshop on method receiver design and interface satisfaction:*
- `core-05-03` **Pointer receivers** *(Part 1: In-place mutation and struct pointer semantics)*
- `core-05-04` **Value receivers** *(Part 2: Immutability, value copying, and safety)*
- `core-05-05` **Receiver sets** *(Part 3: Method set rules and interface satisfaction)*

### Session 3: Interface Abstractions & Dynamic Type Inspection
- `core-05-08` Interfaces
- `core-05-09` Small interfaces
- `core-05-10` Interface embedding
- `core-05-11` Stringer
*Unified Dynamic Type Inspection Workshop:*
- `core-05-12` **Type assertions** *(Single dynamic type extraction with comma-ok idiom)*
- `core-05-13` **Type switches** *(Multi-type branch discrimination across interface values)*
- `core-05-14` Nil interfaces

### Session 4: Packages, Modules & Code Boundaries
- `core-05-15` Package names
- `core-05-16` Export rules
- `core-05-17` Project layout
- `core-05-18` Modules
- `core-05-19` go.mod
- `core-05-20` go.sum
- `core-05-21` Dependency management
- `core-05-22` Workspaces
- `core-05-23` Documentation comments

### Session 5: Modern Type Parameters
- `core-05-24` Generics fundamentals *(Type parameters, constraints, monomorphization, and interface tradeoffs)*

## Projects

- [`projects/bank-account/README.md`](./projects/bank-account/README.md) — **Bank Account**
- [`projects/payroll-processor/README.md`](./projects/payroll-processor/README.md) — **Payroll Processor**

## Assessments

- [`assessments/bank-account/README.md`](./assessments/bank-account/README.md) — **Bank Account**
- [`assessments/checkpoint/README.md`](./assessments/checkpoint/README.md) — **Module 05 Checkpoint — Types, Interfaces, Packages, and Modules**
- [`assessments/payroll-processor/README.md`](./assessments/payroll-processor/README.md) — **Payroll Processor**

## Completion

Complete every required lesson, lab, project, and assessment for this module.

## Next module

After passing this module's exit gate, continue to [Module 06: Testing, Debugging, and Refactoring](../06-testing-debugging-refactoring/README.md).
