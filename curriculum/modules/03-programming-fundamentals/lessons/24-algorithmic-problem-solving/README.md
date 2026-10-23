# Algorithmic Problem-Solving and Complexity

## Learning objective

Decompose programming challenges systematically into discrete algorithmic steps, trace state transformations using desk-checking tables, reason intuitively about time and space complexity ($O(1)$, $O(N)$, $O(\log N)$, $O(N^2)$), and apply idiomatic Go slice algorithms (order-preserving deduplication, two-pointer partitioning, binary search via `slices.BinarySearch`, and custom comparator sorting via `slices.SortFunc`).

## Why this matters

Writing software is more than memorizing syntax; it is the discipline of transforming input data into expected outputs predictably, correctly, and efficiently. Beginners often fall into accidental $O(N^2)$ bottlenecks (such as nested loops for deduplication or lookups) or struggle to debug edge cases because they execute code mentally without structured state tracing. Mastering structured decomposition and complexity reasoning early prevents architectural flaws and prepares you for real-world backend engineering and technical interviews.

## Mental model

An algorithm is an explicit recipe with three essential guarantees:
1. **Defined Boundaries:** Precise contracts on valid inputs and required outputs, including zero-value and boundary cases (empty slice, 1-element slice, duplicates, negative numbers).
2. **Loop Invariants:** Properties that remain true before and after each iteration of a loop.
3. **Resource Scaling (Big-O):** How CPU execution time and heap memory scale as the input size $N$ increases from 10 items to 1,000,000 items.

```
Input Collection (N items)
       │
       ▼
┌───────────────────────────────┐
│ Algorithmic Transformation    │  ◄── Invariants & Desk-checking
│ - Deduplicate (O(N))          │
│ - Partition in-place (O(N))   │
│ - Binary Search (O(log N))    │
└───────────────────────────────┘
       │
       ▼
Deterministic Output Collection
```

## Core idea

### 1. Systematic Problem Decomposition
Before writing a single line of Go:
1. **Clarify Inputs & Outputs:** What types are involved? Can the slice be empty? Does ordering matter?
2. **Identify Edge Cases:** $N=0$, $N=1$, identical elements, pre-sorted data.
3. **Desk-Checking (Trace Table):** Trace variables across 3–4 sample inputs manually on a table before typing code:
   | Iteration `i` | Variable `left` | Variable `right` | Slice State | Invariant Maintained? |
   | :---: | :---: | :---: | :--- | :---: |
   | 0 | 0 | 0 | `[3, 8, 2]` | Yes |
   | 1 | 0 | 1 | `[8, 3, 2]` | Yes (evens moved left) |

### 2. Practical Slice Algorithms
- **Order-Preserving Deduplication:** Use an auxiliary `map[K]struct{}` for $O(1)$ set membership checks while appending first-seen elements to an output slice. Runs in $O(N)$ time and $O(N)$ auxiliary space, compared to the naive $O(N^2)$ double loop.
- **Two-Pointer In-Place Partitioning:** Maintain a `left` pointer demarcating the boundary of processed target items, advancing a `right` pointer across the slice. Exchanges happen in-place without heap re-allocation ($O(N)$ time, $O(1)$ auxiliary space).
- **Binary Search (`slices.BinarySearch`):** On sorted data, halves the search space each step. Finds elements in $O(\log N)$ steps (e.g., at most 20 comparisons for $1,000,000$ elements).
- **Custom Multi-Criteria Sorting (`slices.SortFunc`):** Use standard library comparators (`cmp.Compare`) to cleanly sort compound structures by multiple fields (e.g., priority descending, timestamp ascending) without reflection or `sort.Interface` boilerplate.

## Under the hood

- **CPU Cache Locality:** Slices are backed by contiguous arrays. Iterating linearly over a slice leverages CPU hardware prefetchers and cache lines (64 bytes transferred per memory fetch). Algorithms that maintain contiguous slice layouts are vastly faster than pointer-chasing structures (like linked lists) in modern architectures.
- **Pattern-Defeating Quicksort (pdqsort):** Go's `slices.Sort` and `slices.SortFunc` use pdqsort under the hood. It runs in $O(N \log N)$ worst-case time, $O(N)$ best-case on already sorted or identical slices, and adapts dynamically to avoid worst-case quicksort quadratic degradation.
- **Binary Search Indexing:** `slices.BinarySearch` returns `(index, bool)`. If found, `bool` is `true`. If not found, `index` is the exact insertion point where the element can be inserted to keep the slice sorted.

## How Go uses it

- The Go standard library uses `slices.BinarySearch` and `slices.SortFunc` across compiler internals, HTTP routing tables, and encoding parsers.
- Database drivers and cache implementations use order-preserving deduplication when preparing batch query parameters and foreign key lists.

## Go example

```go
package main

import (
	"cmp"
	"fmt"
	"slices"
)

type Task struct {
	Name     string
	Priority int // 1 = Low, 3 = High
}

func main() {
	tasks := []Task{
		{Name: "Write docs", Priority: 1},
		{Name: "Fix critical outage", Priority: 3},
		{Name: "Review PR", Priority: 2},
		{Name: "Triage bug", Priority: 2},
	}

	// Sort by Priority descending, then Name ascending
	slices.SortFunc(tasks, func(a, b Task) int {
		if c := cmp.Compare(b.Priority, a.Priority); c != 0 {
			return c
		}
		return cmp.Compare(a.Name, b.Name)
	})

	for _, t := range tasks {
		fmt.Printf("Priority %d: %s\n", t.Priority, t.Name)
	}
}
```

## Failure modes & common mistakes

1. **Accidental $O(N^2)$ Lookup:** Checking slice containment with an inner loop inside another loop instead of using a lookup set or map. A list of 100,000 items takes fractions of a millisecond with a map, but over 30 seconds with nested slice searches.
2. **Binary Search on Unsorted Data:** Calling `slices.BinarySearch` on an unsorted slice yields bogus results without compile or runtime errors. Ensure data is sorted first.
3. **Index Mutation During Iteration:** Appending to or slicing the same slice while looping over it with `range` can cause index out of bounds or infinite loops.
4. **Mutating Slice Elements in Value Iteration:** In `for _, v := range slice`, `v` is a copy. Updating `v.Field` does not modify the underlying slice element.

## Practice task

1. Open `main.go` and inspect the implementations of `DeduplicatePreserveOrder`, `PartitionActive`, and `FindPriceTier`.
2. Add a new function `TopKTransactions(txs []Transaction, k int) []Transaction` that returns the top `k` highest-priority transactions without modifying the original slice.
3. Run the test suite:
   ```bash
   go test -v ./curriculum/modules/03-programming-fundamentals/lessons/24-algorithmic-problem-solving
   ```

## Review questions

1. What is the time and space complexity difference between deduplicating a slice using nested loops versus using a `map[string]struct{}`?
2. How many comparisons does `slices.BinarySearch` require in the worst case to locate an item in a sorted slice of $1,000,000$ elements?
3. What does `slices.BinarySearch` return if the target element is NOT present in the slice? How can this return value be used?
4. Why does `slices.SortFunc` outperform older reflection-based or interface-based sorting in Go?
