# DS.5 Slice Sharing And Capacity

## Mission

Learn why sub-slices can still affect the original data and why `append` can reuse spare capacity in
ways that surprise beginners.

## Prerequisites

- `DS.2` slices
- `DS.4` pointers

## Mental Model

A sub-slice is usually another view over the same backing array.
That makes slicing cheap, but it also means two slices can still touch the same stored data.

## Visual Model

```text
original := [0 1 2 3 4 5]
shared   := original[1:4]

original: 0 1 2 3 4 5
shared:     1 2 3
```

```text
shared[0] = 100

original: 0 100 2 3 4 5
shared:     100 2 3
```

```text
growth := original[2:4]
append may still write into the original backing array
if spare capacity exists
```

## Run Instructions

```bash
go run ./03-data-structures/5-slices-2
```

## Code Walkthrough

### `original := []int{0, 1, 2, 3, 4, 5}`

This creates the base slice that the rest of the lesson will inspect.

### `shared := original[1:4]`

This creates a sub-slice view.
It does not copy values into a separate collection.

### `len(shared)` and `cap(shared)`

The lesson prints both measurements so the learner can see that a sub-slice may have:

- a smaller length
- but still a larger capacity than expected

That spare capacity is what makes later `append` behavior surprising.

### `shared[0] = 100`

This changes the first visible element of `shared`.
Because `shared` and `original` still share the same backing array, `original` changes too.

That is the first big warning in the lesson:

- a new slice view is not the same as a copied slice

### `growth := original[2:4]`

This creates another sub-slice to demonstrate append behavior.

### `growth = append(growth, 200)`

This is the second big warning.

If the append fits inside the available capacity, Go can reuse the backing array.
That means the append can overwrite part of the original data.

The next print proves that the original slice changed.

### `independent := make([]int, len(original[2:4]))`

This line starts the safe copy approach.
Instead of sharing, the code allocates a new slice with its own backing storage.

### `for i, value := range original[2:4]`

This loop copies values one by one into the new independent slice.

### `independent[0] = 500`

This final mutation proves the new slice no longer shares storage with the original.

## Try It

1. Change the sub-slice ranges and watch how `len` and `cap` change.
2. Append more than one value to `growth` and inspect whether the original still changes.
3. Change the manual copy loop to copy a different slice range and confirm the independent slice
   stays separate.

## Common Questions

- Why did `append` change the original slice?
  Because the sub-slice still had spare capacity in the same backing array.

- How do I avoid accidental sharing?
  Copy the values into a new slice before making independent changes.

## Production Relevance

This lesson prevents one of the most common slice bugs in Go: changing shared data accidentally
because two slices still point at the same backing array.

## Next Step

Continue to `DS.6` contact directory.
