# Section 14: Testing

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| Unit Testing | Beginner | **Critical** | `testing.T`, `go test` |
| Benchmarking | Advanced | High | `testing.B`, allocation profiling |
| Table-Driven | Intermediate| High | Data-isolated test matrices |

## Engineering Depth
Go builds testing directly into the language via the `testing` package rather than relying on external frameworks like Jest or JUnit. 
- **Zero-Allocation Tracking:** `go test -bench=. -benchmem` gives instant feedback on Heap allocations `allocs/op`. The lower the allocations, the more you avoid triggering the Go Garbage Collector (GC), leading to microsecond sub-latency.

## References
1. **[Go Docs]** [Package testing](https://pkg.go.dev/testing)
