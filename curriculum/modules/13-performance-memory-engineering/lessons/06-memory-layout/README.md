# Memory layout

## Learning objective

Explain how Go aligns struct fields in memory, use `unsafe.Sizeof` and `unsafe.Offsetof` to measure struct sizes and field offsets, reorder struct fields to minimize padding waste, and understand how struct layout affects CPU cache utilization.

## Why this matters

Every byte of memory in a hot struct is multiplied by the number of instances. A struct that is 24 bytes instead of 16 bytes wastes 50% space on every allocation. In a cache with 10 million entries, that is 80 MB of wasted RAM. Beyond memory usage, cache line utilization directly impacts CPU performance. A poorly laid-out struct that spans two cache lines causes two memory fetches instead of one. In high-throughput Go services, struct layout optimization yields measurable latency improvements.

## Mental model

Memory layout is like Tetris: each field is a block of a certain width (its size), and the compiler places them in a row (the struct). If a block cannot start at the current position because it would overlap a boundary (alignment), the compiler adds padding (empty space) before it. The row's total width is padded at the end to be a multiple of the widest block. Poor ordering wastes space (padding); optimal ordering fits all blocks with minimal gaps.

## Core idea

Every Go type has a natural alignment equal to its size (for primitive types). On 64-bit systems:

| Type | Size | Alignment |
|---|---|---|
| `bool`, `int8`, `uint8`, `byte` | 1 byte | 1 |
| `int16`, `uint16` | 2 bytes | 2 |
| `int32`, `float32` | 4 bytes | 4 |
| `int64`, `float64`, `uint64` | 8 bytes | 8 |
| `string` | 16 bytes (ptr + len) | 8 |
| `slice` | 24 bytes (ptr + len + cap) | 8 |
| `interface{}` | 16 bytes (type + data) | 8 |

The compiler ensures every field's offset is a multiple of its alignment by adding padding bytes before it. After all fields, the compiler adds trailing padding to make the total struct size a multiple of the maximum alignment among all fields.

## Under the hood

Go follows the System V AMD64 ABI alignment rules on most platforms. When compiling a struct, the compiler processes fields in declaration order, maintaining a current offset:

1. For each field, compute `offset % field.alignment`. If non-zero, add `field.alignment - (offset % field.alignment)` bytes of padding.
2. Place the field at the new offset.
3. After all fields, compute `finalSize % maxAlignment`. If non-zero, add `maxAlignment - (finalSize % maxAlignment)` bytes of trailing padding.

The resulting layout is visible via `unsafe.Offsetof(field)`. The compiler may reorder fields for optimization (Go 1.18+ attempts field reordering by size), but this does not change observable behavior — field reordering is only applied when it does not affect memory safety. To guarantee a specific layout, use `encoding/binary` or manual serialization.

## How Go uses it

- The standard library's `sync.Mutex` fits in 8 bytes (a single uint32 for non-sema implementations).
- `net/http.Server` field order is carefully arranged to minimize size and avoid false sharing.
- `runtime.g` (goroutine struct) is hundreds of bytes with fields ordered by access frequency and alignment.
- `encoding/binary` reads struct fields by computing offsets manually, relying on the layout guarantee within a single struct.
- `reflect` uses the same offset metadata to implement `reflect.Value.Field()`.

## Go example

```go
package main

import (
	"fmt"
	"unsafe"
)

type BadOrdered struct {
	Flag    bool    // 1 byte + 7 padding
	Amount  float64 // 8 bytes
	Counter int32   // 4 bytes + 4 padding
} // total: 24 bytes

type GoodOrdered struct {
	Amount  float64 // 8 bytes
	Counter int32   // 4 bytes
	Flag    bool    // 1 byte + 3 padding
} // total: 16 bytes

type User struct {
	ID    int64   // 8 bytes
	Name  string  // 16 bytes (pointer + len on 64-bit)
	Score float64 // 8 bytes
} // total: 32 bytes

func main() {
	fmt.Println("BadOrdered size:", unsafe.Sizeof(BadOrdered{}))
	fmt.Println("GoodOrdered size:", unsafe.Sizeof(GoodOrdered{}))

	var bad BadOrdered
	var good GoodOrdered

	fmt.Printf("BadOrdered offsets: Flag=%d, Amount=%d, Counter=%d\n",
		unsafe.Offsetof(bad.Flag), unsafe.Offsetof(bad.Amount), unsafe.Offsetof(bad.Counter))

	fmt.Printf("GoodOrdered offsets: Amount=%d, Counter=%d, Flag=%d\n",
		unsafe.Offsetof(good.Amount), unsafe.Offsetof(good.Counter), unsafe.Offsetof(good.Flag))

	fmt.Println("User size:", unsafe.Sizeof(User{}))
}
```

Output:

```
BadOrdered size: 24
GoodOrdered size: 16
BadOrdered offsets: Flag=0, Amount=8, Counter=16
GoodOrdered offsets: Amount=0, Counter=8, Flag=12
User size: 32
```

`GoodOrdered` saves 8 bytes (33%) by sorting fields by size descending: largest first.

## Step-by-step execution

For `BadOrdered` layout (fields: bool, float64, int32):

1. Start offset = 0. `bool` (size 1, align 1). Offset 0 % 1 = 0. Place `bool` at offset 0. Current offset = 1.
2. `float64` (size 8, align 8). Offset 1 % 8 = 1. Need 7 bytes padding. Place `float64` at offset 8. Current offset = 16.
3. `int32` (size 4, align 4). Offset 16 % 4 = 0. Place `int32` at offset 16. Current offset = 20.
4. Trailing padding: max alignment = 8. 20 % 8 = 4. Need 4 bytes padding. Total = 24.

For `GoodOrdered` (fields: float64, int32, bool):

1. `float64` at offset 0. Current = 8.
2. `int32` at offset 8 (8 % 4 = 0). Current = 12.
3. `bool` at offset 12 (12 % 1 = 0). Current = 13.
4. Trailing padding: max alignment = 8. 13 % 8 = 5. Need 3 bytes padding. Total = 16.

## Common mistakes

- **Declaring fields in logical order instead of size order** — Field order determines struct size. Sort fields by size descending (largest first) to minimize padding.
- **Assuming sizeof(struct) = sum(sizeof(fields))** — Alignment padding adds invisible bytes. Always verify with `unsafe.Sizeof`.
- **Using many bool fields** — Each `bool` is 1 byte but can cause 3-7 bytes of padding between it and a larger field. Use a bitmask (`uint8`) for multiple flags.
- **False sharing in concurrent structs** — If two goroutines access different fields of the same struct and the struct is smaller than a cache line (64 bytes), they contend on the cache line. Add padding to separate hot fields onto different cache lines.

## Debugging walkthrough

A caching layer uses a struct to store entries. Profiling shows high memory usage:

```go
type CacheEntry struct {
	Key       string    // 16 bytes
	ExpiresAt time.Time // 24 bytes (struct{sec,nsec,loc})
	Value     []byte    // 24 bytes (ptr + len + cap)
	HitCount  int32     // 4 bytes
	Flags     bool      // 1 byte
	Active    bool      // 1 byte
}
```

Check size:

```go
fmt.Println(unsafe.Sizeof(CacheEntry{})) // probably 88 bytes with padding
```

Reordered:

```go
type CacheEntry struct {
	ExpiresAt time.Time // 24 bytes — largest field first
	Value     []byte    // 24 bytes
	Key       string    // 16 bytes
	HitCount  int32     // 4 bytes
	Flags     uint8     // 1 byte (merge two bools into bitmask)
	_         [3]byte   // explicit trailing padding
}
```

After reordering, the size drops from 88 to 72 bytes, saving 18% per entry. For 1 million entries, that is 16 MB saved.

## Production notes

- Always verify struct sizes with `unsafe.Sizeof`. Compiler versions and architectures may differ.
- Field ordering is a documentation concern: sort by descending size, but keep related fields together within the same size tier.
- Cache line alignment matters for concurrent structs accessed by different goroutines. Add `_ [cacheLinePad]byte` to separate hot fields.
- Go 1.18+ performs automatic field reordering in some cases, but do not rely on it. Always order explicitly.
- For wire formats (network, disk), do not depend on Go's in-memory layout. Use `encoding/gob`, `encoding/binary`, or protobuf.

## Performance implications

Smaller structs improve performance in three ways: (1) less memory allocation per instance, reducing GC pressure; (2) better cache utilization — more entries fit in each 64-byte cache line; (3) fewer cache line misses when iterating slices of structs. A 24-byte struct vs a 16-byte struct means 33% fewer entries per cache line, which translates to approximately 33% more cache misses. In hot loops, this can be a 20-40% performance difference.

## Practice task

Create a struct `Event` with fields: `Timestamp int64`, `Name string`, `Priority byte`, `RetryCount int32`, `Scheduled bool`. Measure its size with `unsafe.Sizeof`. Reorder the fields for minimum size. What is the original size vs the optimized size? Verify with `unsafe.Offsetof` that each field is at the expected offset.

## Tests / verification

```bash
go run ./curriculum/modules/13-performance-memory-engineering/lessons/06-memory-layout
go test ./curriculum/modules/13-performance-memory-engineering/lessons/06-memory-layout
```

## Review questions

1. Why does `unsafe.Sizeof(BadOrdered{})` return 24 when the sum of field sizes is 13?
2. What rule determines the offset of a struct field in memory?
3. How would you reorder fields `bool`, `string`, `int64`, `int32` to minimize struct size?
4. What is false sharing, and how does struct padding help prevent it?
5. Can you rely on the Go compiler to always reorder fields optimally? Why or why not?

## NEXT UP

Caching basics — building in-memory caches in Go, using `sync.Map` for concurrent access, implementing LRU eviction with TTL, and comparing read-through vs write-through cache strategies.
