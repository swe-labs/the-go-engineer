# Section 3: Collections and Pointers

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Arrays | Beginner | Low | Contiguous fixed-size memory |
| Slices | Intermediate | **Critical** | Dynamic window over underlying arrays, capacity vs length |
| Maps | Intermediate | High | Hash table, $O(1)$ amortized lookups |
| Pointers | Advanced | High | Stack vs Heap allocation, escape analysis |

## Engineering Depth
In this section, you will learn the physical memory layout of Go collections.
Unlike other languages, Go gives you direct control over memory behavior. You will understand why `append` operations can be $O(N)$ when capacity is breached (causing Heap re-allocations), and how to pre-allocate slices to maintain $O(1)$ performance.

## References
1. **[Official Blog]** [Go Slices: usage and internals](https://go.dev/blog/slices-intro)
2. **[Book]** Effective Go - [Data](https://go.dev/doc/effective_go#data)
3. **[Official Blog]** [Go maps in action](https://go.dev/blog/maps)

---

## 🏗 Exercise: Contact Manager System (`6-contact-manager`)

This section's capstone forces you to manipulate slices and pointers safely.

### Step-by-Step Instructions & Hints
1. **Define the Data Model:** Create a `Contact` struct.
2. **Setup the Storage:** Create a slice `[]Contact` to act as your in-memory database.
   - *Hint:* Pre-allocate capacity using `make([]Contact, 0, 10)` to prevent early memory reallocation.
3. **Implement Add:** Write a function that appends a new contact.
   - *Hint:* Because `append` might allocate a new underlying array, return the mutated slice.
4. **Implement Delete (Advanced):** Remove an element from the middle of the slice.
   - *Hint:* Use `append(slice[:index], slice[index+1:]...)`. Understand that this shifts elements in memory and is an $O(N)$ operation.
5. **Implement Update:**
   - *Hint:* Pass the contact slice or a specific contact by **Pointer** (`*Contact`) so the modifications stick. Passing by value will only modify a copy!


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| DS.1 | [arrays](./1-array) | Fixed-size contiguous memory · value-type copy | 🟢 entry |
| DS.2 | [slices](./2-slices) | Slice header · make · append · capacity growth | DS.1 |
| DS.3 | [maps](./3-maps) | Hash table · O(1) lookup · comma-ok · delete | DS.2 |
| DS.4 | [pointers](./4-pointers) | &amp; · * · pass-by-value vs reference · nil | DS.1, DS.2 |
| DS.5 | [slices-2](./5-slices-2) | Sub-slicing · shared backing array · re-allocation trap | DS.2, DS.4 |
