# DS.2 Slices

## Mission

Learn how Go represents dynamic collections through slices, and why `len`, `cap`, `make`, and `append` are all part of one connected idea.

## Why This Lesson Exists Now

Arrays taught you about fixed-size collections and value copying. But real programs need to handle data that grows and shrinks.

That is what slices do. They are Go's primary way to work with dynamic collections.

This lesson builds on DS.1 by showing what happens when you need flexibility.

## Prerequisites

- `DS.1` arrays

## Mental Model

A slice is a small view over an underlying array.
It tracks:

- which array data it points to
- how many elements are currently in the slice
- how much capacity remains before growth needs a new backing array

## Visual Model

```mermaid
graph TD
    A["many values"] --> B["one collection type"]
    B --> C["read or update by position or key"]
```
```text
items := make([]int, 0, 3)

slice header:
- pointer -> backing array
- len     -> 0
- cap     -> 3
```

```text
after append 10, 20, 30

backing array: [10 20 30]
len = 3
cap = 3
```

```text
after append 40

Go may allocate a larger backing array
so the slice can keep growing
```

## Machine View

A slice is not an array. A slice is a data structure that points to an array.

When you create a slice, Go creates a small "header" with three fields:
- **pointer**: the address of the first element in the backing array
- **len**: how many elements are currently accessible
- **cap**: how many elements can be added before the backing array must be reallocated

When you pass a slice to a function, you are passing this small header, not copying the whole backing array.

When `append` needs more capacity, it allocates a new, larger backing array and copies the old elements over.

## Run Instructions

```bash
go run ./02-language-basics/04-data-structures/2-slices
```

## Code Walkthrough

### `names := []string{"Alice", "John", "Mark"}`

This line creates a slice literal.
Unlike the array lesson, there is no fixed size written in the type.

That tells the learner:

- this is a slice, not an array
- slices are the normal collection tool for variable-length lists

### `fmt.Printf("len=%d cap=%d\n", len(names), cap(names))`

This prints the two slice measurements:

- `len`: how many elements are currently visible
- `cap`: how much space is available in the current backing array view

### `items := make([]int, 0, 3)`

This is one of the most important slice lines in the lesson.

It means:

- make a slice of `int`
- start with length `0`
- reserve capacity `3`

The slice starts empty, but it already has room to grow.

### `items = append(items, 10)` and the next two append lines

Each `append` returns the updated slice, so the result must be assigned back into `items`.

That is why the code does not write only `append(items, 10)`.
If the learner forgets the reassignment, they lose the updated slice value.

### `items = append(items, 40)`

This append matters because it grows beyond the original capacity of `3`.

That is the first hint that:

- slices can grow
- growth may require a different backing array

### `firstTwo := items[:2]`

This creates a smaller view over the same data.
It does not make a copy.

### `lastTwo := items[2:]`

This creates another view, starting from index `2` to the end.

These two lines teach the learner that slicing syntax makes views, not independent new collections
by default.

## Try It

1. Change the `make` call to `make([]int, 1, 3)` and watch how the starting length changes.
2. Comment out one `items =` reassignment and see why `append` must be captured.
3. Change `firstTwo := items[:2]` to `firstThree := items[:3]` and inspect the output again.

## Common Questions

- Why does `append` return a slice?
  Because growth may change the slice header or even the backing array.

- Does `items[:2]` copy the first two values?
  No. It usually creates another view over the same underlying data.

## ⚠️ In Production
Slices are everywhere in Go.
Understanding `len`, `cap`, `make`, and `append` prevents a huge amount of confusion later in file
processing, HTTP work, concurrency, and general application code.

## 🤔 Thinking Questions

1. What problem is this lesson trying to solve?
2. What would change if you removed this idea from the program?
3. Where do you expect to see this pattern again in real Go code?
## Next Step

Continue to `DS.3` maps.
