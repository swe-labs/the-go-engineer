# Track 2: Control Flow (CF)

## Mission

Learn how a Go program chooses what to do next and how it repeats work without duplicated code.

By the end of this track, a learner should be able to:
- Branch execution paths with `if`, `else if`, and `else`.
- Repeat work safely using the `for` keyword.
- Choose between multiple cases efficiently with `switch`.
- Control loop iteration with `break` and `continue`.
- Schedule cleanup work and resource management with `defer`.

## Zero-Magic Boundary

This track stays focused on **Execution Logic**.
It does **not** formally teach:
- Data structure internals (covered in Track 3).
- Helper function design (covered in Section 03).
- Multi-file package organization (covered in Section 05).

> [!NOTE]
> You may see simple slices (lists) used as loop targets. These are a preview of the [Data Structures](../04-data-structures/README.md) track.

## Track Map

| ID | Lesson | What It Unlocks |
| --- | --- | --- |
| `CF.1` | [If / Else](./01-if-else/) | Introduces branching and decision-making. |
| `CF.2` | [For Basics](./02-for-basics/) | Teaches Go's only loop keyword and repeated work. |
| `CF.3` | [Break / Continue](./03-break-continue/) | Early exit and selective skipping inside loops. |
| `CF.4` | [Switch](./04-switch/) | Readable multi-branch decision logic. |
| `CF.5` | [Defer Basics](./05-defer-basics/) | Cleanup scheduling and LIFO (Last-In, First-Out) order. |
| `CF.6` | [Defer Use Cases](./06-defer-use-cases/) | Real-world patterns for resource management. |
| `CF.7` | [Pricing Checkout](./07-pricing-checkout/) | **Milestone**: Logic-heavy exercise for state transitions. |

## Next Step

After completing this track, continue to [Track 3: Data Structures](../04-data-structures/README.md).
