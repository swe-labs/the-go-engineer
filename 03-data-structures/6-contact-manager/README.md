# DS.6 Contact Directory

## Mission

Build a small in-memory contact directory that turns the Section 03 data-structure lessons into one
coherent runnable exercise.

This is the Section 03 milestone.
It is where slices, maps, and pointers come together without depending on later-section abstractions
like helper-function design or struct-heavy modeling.

## Prerequisites

Complete these first:

- `DS.1` arrays
- `DS.2` slices
- `DS.3` maps
- `DS.4` pointers
- `DS.5` advanced slicing

## What You Will Build

Implement a small contact directory that:

1. stores names, emails, and phone numbers in parallel slices
2. uses a map for efficient name lookup
3. updates stored data through a pointer to a slice element
4. prints a small demonstration flow in `main()`
5. stays inside the concepts already taught in Section 03

## Files

- [main.go](./main.go): complete solution with teaching comments
- [_starter/main.go](./_starter/main.go): starter file with TODOs

## Run Instructions

Run the completed solution:

```bash
go run ./03-data-structures/6-contact-manager
```

Run the starter scaffold:

```bash
go run ./03-data-structures/6-contact-manager/_starter
```

## Success Criteria

Your finished solution should:

- add and retrieve contact data safely
- show at least one update that persists correctly
- keep the flow simple enough that the data-structure choices are visible
- avoid hiding the exercise behind helper functions that belong more naturally to Section 04

## Common Failure Modes

- appending to a slice without reusing the updated slice value
- using a pointer before you are sure the lookup index exists
- letting one contact's slice positions drift out of sync with the others
- hiding the data-structure lesson under architecture that has not been taught yet

## Next Step

After this milestone, continue to Section 04: Functions and Errors.
