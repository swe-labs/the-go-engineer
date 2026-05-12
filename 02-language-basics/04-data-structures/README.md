# Track 3: Data Structures (DS)

## Mission

Learn how to organize, store, and manipulate collections of data using Go's fundamental data structures.

By the end of this track, a learner should be able to:
- Explain the relationship between fixed-size Arrays and dynamic Slices.
- Perform efficient key-based lookups using Maps.
- Manage memory addresses and shared state using Pointers.
- Reason about value copying vs. pointer mutation at the machine level.

## Zero-Magic Boundary

This track focuses on **Data Organization**.
It does **not** formally teach:
- Advanced struct design (covered in Section 04).
- Concurrency safety for maps (covered in Section 07).
- Complex algorithm complexity analysis (O-notation).

> [!NOTE]
> Pointers are introduced here specifically to explain how to mutate slices and maps efficiently, providing a bridge to [Section 04: Type Design](../../04-types-design/README.md).

## Track Map

| ID | Lesson | What It Unlocks |
| --- | --- | --- |
| `DS.1` | [Arrays](./01-array/) | Foundation for contiguous memory storage. |
| `DS.2` | [Slices](./02-slices/) | Go's most powerful tool for dynamic sequences. |
| `DS.3` | [Maps](./03-maps/) | Fast lookup and associative data storage. |
| `DS.4` | [Pointers](./04-pointers/) | Direct memory access and shared state control. |
| `DS.5` | [Slices in Depth](./05-slices-2/) | Length, capacity, and the underlying backing array. |
| `DS.6` | [Contact Manager](./06-contact-manager/) | **Milestone**: Building a searchable data registry. |

## Next Step

After completing this track, you have finished Section 02. Continue to [Section 03: Functions & Errors](../../03-functions-errors/README.md).
