# Section 02: Language Basics

## Mission

This section teaches the fundamental building blocks of Go logic: how to store state, control execution flow, and organize data in memory.

By the end of this section, a learner should be able to:
- Declare and update variables using appropriate types and zero values.
- Use constants and enums to model immutable and categorical data.
- Control program flow using branching (if/else, switch) and iteration (for).
- Master the mechanics of `defer` for resource cleanup.
- Organize data using arrays, slices, and maps.
- Understand how pointers enable memory efficiency and shared state.

## Section Map

### Track 1: Values and Types (LB)
| ID | Lesson | What It Unlocks |
| --- | --- | --- |
| `LB.1` | [Variables](./01-variables/) | Introduces typed values, zero values, and declaration styles. |
| `LB.2` | [Constants](./02-constants/) | Separates immutable values from ordinary variables. |
| `LB.3` | [Enums with iota](./03-enums/) | Shows how Go models enum-like values with typed constants. |
| `LB.4` | [Application Logger](./04-application-logger/) | **Milestone**: Combines LB track concepts into a practical logger. |

### Track 2: Control Flow (CF)
| ID | Lesson | What It Unlocks |
| --- | --- | --- |
| `CF.1` | [If / Else](./03-control-flow/01-if-else/) | Basic conditional branching and scope. |
| `CF.2` | [For Basics](./03-control-flow/02-for-basics/) | The single iteration keyword in Go. |
| `CF.3` | [Break / Continue](./03-control-flow/03-break-continue/) | Controlling loop execution. |
| `CF.4` | [Switch](./03-control-flow/04-switch/) | Clean multi-way branching without boilerplate. |
| `CF.5` | [Defer Basics](./03-control-flow/05-defer-basics/) | Mechanics and LIFO execution order. |
| `CF.6` | [Defer Use Cases](./03-control-flow/06-defer-use-cases/) | Real-world resource management patterns. |
| `CF.7` | [Pricing Checkout](./03-control-flow/07-pricing-checkout/) | **Milestone**: Logic-heavy exercise for state transitions. |

### Track 3: Data Structures (DS)
| ID | Lesson | What It Unlocks |
| --- | --- | --- |
| `DS.1` | [Arrays](./04-data-structures/01-array/) | Fixed-length contiguous memory. |
| `DS.2` | [Slices](./04-data-structures/02-slices/) | Dynamic, growable views into arrays (the most used structure). |
| `DS.3` | [Maps](./04-data-structures/03-maps/) | Key-value storage and fast lookups. |
| `DS.4` | [Pointers](./04-data-structures/04-pointers/) | Memory addresses, dereferencing, and shared state. |
| `DS.5` | [Slices in Depth](./04-data-structures/05-slices-2/) | Length, capacity, and the underlying backing array. |
| `DS.6` | [Contact Manager](./04-data-structures/06-contact-manager/) | **Milestone**: Integrated exercise using slices, maps, and pointers. |

## Zero-Magic Boundary

This section focuses on the **Foundations** level.
It does not try to teach:
- Functions as a design tool (covered in Section 03).
- Custom structs and interfaces (covered in Section 04).
- Concurrency (covered in Section 07).

## Next Step

After completing this section, continue to [Section 03: Functions & Errors](../03-functions-errors/README.md).
