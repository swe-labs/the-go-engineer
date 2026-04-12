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
- `DS.5` slice sharing and capacity

## What You Will Build

Implement a small contact directory that:

1. stores names, emails, and phone numbers in parallel slices
2. uses a map for efficient name lookup
3. updates stored data through a pointer to a slice element
4. prints a small demonstration flow in `main()`
5. stays inside the concepts already taught in Section 03

## Visual Model

```text
index  name               email                phone
0      Alice Wonderland  alice@example.com    111-2222
1      Bob The Builder   bob@example.com      333-4444
2      Charlie Brown     charlie@example.com  555-6666
```

```text
indexByName

"Alice Wonderland" -> 0
"Bob The Builder"  -> 1
"Charlie Brown"    -> 2
```

## Why This Milestone Avoids Structs

Structs matter, but they belong to a later section.

This milestone is intentionally simpler than a real application because the goal here is to prove
that the learner understands slices, maps, pointers, shared indexing, and persistent updates before
adding new modeling tools.

## Files

- [main.go](./main.go): complete runnable solution
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

Run the automated verification surface:

```bash
go test ./03-data-structures/6-contact-manager
```

## Recommended Learning Flow

1. Read this README first.
2. Open [_starter/main.go](./_starter/main.go) and list the pieces you need.
3. Try to build the milestone yourself.
4. Run the starter and your solution often.
5. Compare your result with [main.go](./main.go) only after you have attempted the design.

## Solution Walkthrough

The solution stays inside Section 03 concepts on purpose.

### 1. Parallel slices set the storage model

The solution starts with:

- `names`
- `emails`
- `phones`

Each contact uses the same index in all three slices.
That is why the first design rule is "keep the indices aligned."

### 2. The map turns a name into an index

`indexByName` is a `map[string]int`.

It does not store the contact data itself.
It stores the slice position where that contact's data lives.

That lets the solution answer:

- "Where is Bob?"

before it answers:

- "What is Bob's phone number?"

### 3. Each append must keep all slices in sync

When the solution adds Alice, Bob, and Charlie, it appends to:

- `names`
- `emails`
- `phones`

Then it stores `len(names) - 1` in the map.

That line matters because the newly added contact always lands at the last valid index.

### 4. The duplicate check uses the map first

The duplicate guard asks:

```go
if _, exists := indexByName["Alice Wonderland"]; exists {
```

This uses the comma-ok pattern from the maps lesson.
It avoids creating a second contact entry with the same name.

### 5. Listing proves the shared index model

The `for i := 0; i < len(names); i++ { ... }` loop is deliberately simple.

It prints:

- `names[i]`
- `emails[i]`
- `phones[i]`

side by side so the learner can see that one contact is really "one index across several slices."

### 6. Pointer-based update is the real milestone proof

The important sequence is:

1. use the map to get Bob's index
2. take `&phones[bobIndex]`
3. write through that pointer
4. read the slice again to prove the update persisted

That is the exact Section 03 idea chain:

- maps find the right position
- slices hold the stored values
- pointers let the update stick

### 7. Missing lookup still uses comma-ok

The final `Zack` check reminds the learner that not every key exists.
The map lesson still matters here, even inside the milestone.

## Try It

1. Add one more contact and update the map with the new index.
2. Change the contact you update from Bob to Charlie.
3. Break the alignment on purpose by skipping one append, then explain why the output becomes wrong.
4. Change the duplicate check to another name and watch how the guard behaves.

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

## Verification Surface

Use these three proof surfaces together:

1. `go run ./03-data-structures/6-contact-manager`
2. `go run ./03-data-structures/6-contact-manager/_starter`
3. `go test ./03-data-structures/6-contact-manager`

The tests are not the lesson.
They are a confidence check that the visible milestone behavior still matches the README contract.

## Questions This Milestone Should Answer

- Why use slices and a map together instead of only one of them?
- Why is the map value an index instead of a phone number?
- Why does taking `&phones[index]` let the update persist?
- Why does this milestone stop before structs and helper functions?

## Next Step

After this milestone, continue to Section 04: Functions and Errors.
