# 06 Types and Interfaces

## Mission

This section teaches you to model data with structs, attach behavior with methods, and define behavior contracts with interfaces instead of inheritance.

By the end of this section, a learner should be able to:

- write structs that model real domain data
- attach methods with clear receiver intent
- define interfaces that describe behavior boundaries
- write small generic helpers built on interface-aware constraints
- write polymorphic code that works across multiple concrete types

## Why This Section Exists Now

The learner already knows:

- functions and multiple return values
- explicit error handling
- basic pointer usage
- slices and maps

That is enough to ask the next engineering questions:

- how do I model a whole thing with multiple related attributes?
- how do I add behavior to that data?
- how do I write code that works with multiple types that share behavior?

Those are types and interfaces questions.

## Zero-Magic Boundary

This section intentionally stays inside foundations-ready ideas:

- struct definition and field access
- method definition with value/pointer receivers  
- interface definition and implicit satisfaction
- basic generics with constraints

**Canonical Foundational Path** (`TI.1`-`TI.7`, `TI.9`, then `TI.8`):
- Core types and interfaces concepts every learner needs
- Ends with the payroll processor milestone

**Advanced/Optional Path** (TI.10 - TI.17):
- Advanced generics, type patterns, data structures
- Marked as stretch/optional - not required for foundations
- Can be explored after completing core path or in later stages

## Beta Stage Ownership

This section belongs to [04 Types & Design](../../docs/stages/04-types-design.md).

Inside the new repo architecture, it is the sixth foundations section:

1. `01-getting-started`
2. `02-language-basics`
3. `03-control-flow`
4. `04-data-structures`
5. `05-functions-and-errors`
6. `06-types-and-interfaces`

## Section Map

### Canonical Foundational Path (Core - Required)

| ID | Type | Surface | Core Job |
| --- | --- | --- | --- |
| `TI.1` | Lesson | [structs](./1-struct) | teach struct as grouped data |
| `TI.2` | Lesson | [methods](./2-methods) | teach method attachment and receiver choice |
| `TI.3` | Lesson | [interfaces](./3-interfaces) | teach implicit contracts and polymorphism |
| `TI.4` | Lesson | [interface embedding](./4-interface-embedding) | teach embedding interfaces in interfaces |
| `TI.5` | Lesson | [Stringer](./5-stringer) | teach fmt.Stringer implementation |
| `TI.6` | Lesson | [type switch](./6-type-switch) | teach handling multiple types from interface |
| `TI.7` | Lesson | [receiver sets](./7-receiver-sets) | teach value vs pointer receiver method sets |
| `TI.9` | Lesson | [generics](./9-generics) | teach type parameters and constraints |
| `TI.10` | Exercise | [payroll processor](./10-payroll-processor) | prove structs+methods+interfaces+generics together |

### Advanced/Optional Path (Stretch)

| ID | Type | Surface | Core Job |
| --- | --- | --- | --- |
| `TI.8` | Lesson | [custom errors](./8-custom-errors) | teach custom error type definition |
| `TI.11` | Lesson | [advanced generics](./10-advanced-generics) | push generics further |
| `TI.12` | Lesson | [empty interface](./11-empty-interface) | teach any/interface{} usage |
| `TI.13` | Lesson | [type assertions](./12-type-assertions) | teach extracting concrete types from interfaces |
| `TI.14` | Lesson | [nil interfaces](./13-nil-interfaces) | teach nil vs typed-nil interface behavior |
| `TI.15` | Lesson | [functional options](./14-functional-options) | teach configurable API pattern |
| `TI.16` | Lesson | [method values](./15-method-values) | teach methods as first-class values |
| `TI.17` | Lesson | [complex constraints](./16-complex-generics) | teach interface constraints, comparable |
| `TI.18` | Lesson | [generic structures](./17-generic-structures) | teach Stack, Queue, Set |

## Current Rebuild Goal

This section is being rebuilt so that:

- learner-facing explanation stays in lesson README.md files
- main.go stays runnable and clean
- the milestone proves earned foundations concepts only
- old Section 05 material can be kept as legacy reference without confusing the new primary path

## Suggested Learning Flow

1. Read each lesson README.md first.
2. Open main.go only after the lesson mission and machine view are clear.
3. Run the code and compare the output with the walkthrough.
4. Change one thing at a time using the Try It prompts.
5. Move to the milestone only after the early lessons stop feeling mechanical.

## Next Step

After completing the canonical path (`TI.1` through `TI.7`, then `TI.9`, then `TI.8`), the learner should be ready to move into the next section: **Composition** (`s06`).

The advanced lessons (TI.10-TI.17) are optional/stretch and can be explored later or in parallel with other sections.
