# Maps

## Learning objective

Create, read, update, and delete entries in Go maps (`map[K]V`), distinguish nil maps from empty maps, predict iteration order behavior, and explain why maps are reference types.

## Why this matters

Maps are the most commonly used associative data structure in Go. They appear in virtually every production Go service — counting requests, grouping records, caching lookups, storing configuration, and implementing sets. Unlike slices, maps provide O(1) amortized lookups by key rather than O(n) scanning. Understanding map internals (hash collisions, resizing, nil behavior) prevents panics, data races, and silent data loss in production systems.

## Mental model

Think of a map as a **hash-powered lookup table**. You give it a key — any comparable type — and it computes a hash code, then uses that hash to store or find the value in an internal bucket array. The map does not maintain insertion order; it scatters entries across buckets for fast access. A nil map is like a locked filing cabinet: you cannot write to it, but you can safely ask whether a key exists (it won't be there).

## Core idea

A map in Go has the type `map[K]V` where `K` is the key type (must be comparable: booleans, numbers, strings, pointers, channels, structs with comparable fields, and arrays) and `V` is the value type (any type, including another map or slice).

Three ways to create a map:

```go
var m map[string]int         // nil map, no storage allocated
m2 := make(map[string]int)   // empty map, ready to use
m3 := map[string]int{        // map literal with initial entries
    "alice": 30,
    "bob":   25,
}
```

Basic operations:

| Operation | Syntax | Nil map behavior |
|---|---|---|
| Insert/update | `m["key"] = val` | **Panic** |
| Lookup | `v := m["key"]` | Returns zero value (no panic) |
| Two-value lookup | `v, ok := m["key"]` | `v` = zero value, `ok` = false |
| Delete | `delete(m, "key")` | No-op (no panic) |
| Length | `len(m)` | Returns 0 |
| Range | `for k, v := range m` | No iterations |

## Under the hood

A Go map is implemented as `hmap` in the runtime (`runtime/map.go`). The structure contains:

- A pointer to a **bucket array** (8 key-value pairs per bucket).
- A **hash seed** (random per-map to prevent HashDoS attacks).
- A **count** of entries.
- **B** (log₂ of the bucket count) for incremental growth.

When you insert, Go hashes the key, masks it to the bucket count, and stores the key-value pair in the appropriate bucket (or an overflow bucket if the primary bucket is full). When the **load factor** exceeds ~6.5 entries per bucket, Go **grows** the bucket array by 2× and rehashes all entries — but it does this **incrementally**, a few buckets at a time, to avoid latency spikes.

Maps are **reference types**: a map variable holds a pointer to the underlying `hmap` structure. Copying the variable or passing it to a function shares the same underlying data.

## How Go uses it

Maps are idiomatic for:

- **Counting**: `counts[item]++` tallies occurrences.
- **Grouping**: `groups[key] = append(groups[key], val)` collects items by category.
- **Memoization / caching**: cache results of expensive computations.
- **Set simulation**: `set := map[string]bool{"a": true}`.
- **Lookup tables**: replacing long if-else chains or switch statements.
- **Struct field indexes**: map from a field value to the struct pointer for fast access.

The `sync.Map` type provides a concurrent-safe map for high-contention scenarios, but for most cases a regular `map` protected by a `sync.RWMutex` is preferred.

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	ages := map[string]int{
		"alice": 30,
		"bob":   25,
	}
	ages["carol"] = 35

	fmt.Println("Initial map:", ages)

	ages["bob"] = 26
	fmt.Println("Bob's updated age:", ages["bob"])

	delete(ages, "alice")
	fmt.Println("After delete:", ages)

	_, exists := ages["alice"]
	fmt.Println("Alice exists?", exists)

	v, ok := ages["carol"]
	fmt.Printf("Carol lookup: value=%d, ok=%v\n", v, ok)

	word := "mississippi"
	freq := make(map[rune]int)
	for _, ch := range word {
		freq[ch]++
	}
	fmt.Println("Character frequencies:", freq)

	var nilMap map[string]int
	fmt.Println("nil map len:", len(nilMap))
	fmt.Println("nil map lookup:", nilMap["x"])
	// nilMap["x"] = 1 // would panic
}
```

## Step-by-step execution

Consider `ages["bob"] = 26` on a map with bucket count 8:

1. Go computes Bob's hash from the key `"bob"` using the per-map seed.
2. The hash is masked: `hash & (8-1)` selects bucket index (e.g., bucket 3).
3. Go walks bucket 3's chain (the primary bucket and any overflow buckets), comparing keys.
4. It finds an existing entry with key `"bob"` at position 2 in the bucket.
5. Go overwrites the value at that position with `26`.
6. If the key had not existed, Go would find an empty slot and insert a new entry, incrementing the count.

For `delete(ages, "alice")`:

1. Hash `"alice"`, mask to bucket index.
2. Walk the bucket chain looking for a matching key.
3. Found: Go zeroes the key and value memory and sets a "tombstone" bit so the bucket is not considered full.
4. If not found: `delete` is a no-op.

## Common mistakes

- **Assigning to a nil map panics**: `var m map[string]int; m["x"] = 1` → panic. Always `make` the map first or use a literal.
- **Assuming iteration order is deterministic**: Map iteration order is intentionally randomized. Code that depends on order will break across runs, machines, or Go versions.
- **Reading from a nil map doesn't panic**: `var m map[string]int; _ = m["x"]` returns 0. This can mask bugs where you expect entries to exist.
- **Map is not safe for concurrent use**: Concurrent read+write or concurrent writes panic with `"concurrent map writes"`. Use `sync.RWMutex` or `sync.Map`.
- **Using non-comparable types as keys**: Slices, maps, and functions are not comparable and cannot be used as map keys.
- **Forgetting that map is a reference type**: Passing a map to a function and modifying it modifies the caller's map — no pointer needed.

## Debugging walkthrough

Consider this code that crashes:

```go
package main

import "fmt"

func main() {
	var users map[string]int
	users["alice"] = 30 // panic: assignment to entry in nil map
	fmt.Println(users)
}
```

**Symptom**: Panic at runtime.

**Investigation**: The declaration `var users map[string]int` creates a nil map. The map variable exists but points to no underlying `hmap` structure.

**Root cause**: Attempting to write to a nil map.

**Fix**: Initialize before use:

```go
users := make(map[string]int)
// or
users := map[string]int{}
```

**Another example**: A caching function that silently returns zero values for missing keys:

```go
func getUserScore(cache map[string]int, name string) int {
	return cache[name]
}
```

If `cache` is nil or missing entries, this returns 0 — indistinguishable from a user who actually scored 0. Use the comma-ok idiom to distinguish:

```go
func getUserScore(cache map[string]int, name string) (int, bool) {
	score, ok := cache[name]
	return score, ok
}
```

## Production notes

- Maps are **not safe for concurrent access**. Even concurrent reads are safe, but a concurrent write alongside any other access (read or write) causes a runtime panic. Use `sync.RWMutex`:
  ```go
  type SafeCache struct {
      mu   sync.RWMutex
      data map[string]Result
  }
  func (c *SafeCache) Get(k string) (Result, bool) {
      c.mu.RLock()
      defer c.mu.RUnlock()
      v, ok := c.data[k]
      return v, ok
  }
  ```
- Map creation with a size hint avoids resizing: `make(map[string]int, 1000)` pre-allocates buckets for 1000 entries.
- Avoid storing large values in maps — store pointers to structs instead.
- `delete` does not shrink map memory. If you need to free memory, recreate the map.
- For sets, `map[K]struct{}` is more memory-efficient than `map[K]bool` because `struct{}{}` uses zero bytes of storage.

## Performance implications

| Operation | Amortized cost | Notes |
|---|---|---|
| Insert | O(1) | O(n) during growth; incremental growth spreads this out |
| Lookup | O(1) | Degrades to O(n) with many hash collisions |
| Delete | O(1) | Tombstone marking |
| Range | O(n) | Order random; n = number of entries |
| len | O(1) | Stored as a field |

- Map lookup is roughly 2-5x slower than array/slice indexing.
- Each map allocation, growth, and rehash incurs GC scanning cost for keys and values.
- If performance-critical, consider whether a sorted slice + binary search (for small, stable datasets) or a custom hash table (for extreme requirements) is better.
- Integer keys are faster than string keys (less hashing overhead).

## Practice task

Write a word counter that:

1. Reads a multi-line string (hardcoded) containing at least 50 words with punctuation.
2. Strips common punctuation (`. , ! ? ; :`).
3. Converts all words to lowercase.
4. Counts how many times each word appears.
5. Prints the top 5 most frequent words and their counts.
6. Handles the case where the input string is empty (returns an empty result gracefully).

Write a function `WordFrequency(text string) map[string]int` and a `TopWords(freq map[string]int, n int) []struct{Word string; Count int}` that returns the top n entries.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/17-maps
go test ./curriculum/modules/03-programming-fundamentals/lessons/17-maps
```

## Review questions

1. What happens when you assign a value to a key in a nil map? What about reading from one?
2. Why is the iteration order of a map non-deterministic? How does Go enforce this?
3. What is the two-value form of map access and why should you use it?
4. How does `delete` behave when the key does not exist in the map?
5. How would you implement a set of integers in Go? Why would you choose `struct{}` over `bool` as the value type?

## NEXT UP

Comma-ok idiom (core-03-18): Learn the two-value form pattern used not just for maps but for type assertions and channel receives.
