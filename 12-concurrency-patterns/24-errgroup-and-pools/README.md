# Section 24: errgroup & sync.Pool

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
|-------|-------|------------|---------------------|
| `errgroup.Group` | Intermediate | **Critical** | Concurrent error collection, idiomatic WaitGroup replacement |
| `errgroup` + context | Advanced | **Critical** | Automatic cancellation on first error |
| `sync.Pool` | Advanced | High | Object reuse, GC pressure reduction |
| Bounded worker pool | Expert | High | Semaphore + errgroup, back-pressure |

## Engineering Depth

`errgroup` is the missing piece between WaitGroup and channels. WaitGroup tells you *when* goroutines finish but discards their errors. Channels collect errors but require manual synchronization. `errgroup.Group` does both: it waits for all goroutines and returns the first non-nil error, cancelling all remaining work via context.

`sync.Pool` is the most impactful memory optimization available in Go. A pool holds recycled objects. Instead of `make([]byte, 4096)` on every request (triggering GC), you `Get()` a pre-allocated buffer, use it, and `Put()` it back. The standard library uses pools everywhere: `fmt`, `encoding/json`, `net/http` all have internal pools.

**When to use sync.Pool:**
- Object allocation appears in pprof heap profile hot path
- Objects are temporary: used briefly, then discarded
- Object size is non-trivial: byte buffers, structs with multiple fields

**When NOT to use sync.Pool:**
- Objects hold state between requests (race condition)
- Allocation is not the bottleneck (profile first!)
- Objects outlive the goroutine that created them

## References

- [golang.org/x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup)
- [sync.Pool](https://pkg.go.dev/sync#Pool)
- [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)
