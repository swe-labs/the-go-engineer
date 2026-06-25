# Race detector

## Learning objective

Use `go build -race` and `go test -race` to detect data races in Go programs, interpret the report, and integrate race detection into CI.

## Why this matters

Data races are the most insidious bugs in concurrent Go programs. They can go undetected for months, producing wrong results only under specific load patterns or on particular hardware. The Go race detector is a compile-time instrumentation tool that catches races at runtime with near-zero false positives. Every professional Go team runs the race detector as part of CI. Without it, shipping concurrent Go code is reckless.

## Mental model

The race detector records every memory access (read or write) of every goroutine along with a happens-before timestamp. When it observes two unsynchronized accesses to the same memory location at overlapping times, it prints a report with the stack traces of both goroutines and the exact line of the access. It is a dynamic analysis tool: it only finds races that are triggered during the execution. Running tests without the race detector is like flying an airplane without a preflight checklist.

## Core idea

The race detector is not a separate binary. When you pass `-race` to `go build`, `go run`, or `go test`, the compiler instruments every memory access with a call to the runtime's race detector library (a C++ ThreadSanitizer). The instrumented binary runs 5-10x slower and uses 2-5x more memory. When a race is detected, the runtime prints a report to stderr and the program continues (or you can make it fail by using `go test -race` which treats races as test failures).

Components of a race report:

- `WARNING: DATA RACE` -- header.
- `Read at 0x... by goroutine N` -- the stack trace of the read.
- `Previous write at 0x... by goroutine M` -- the stack trace of the conflicting write.
- `Goroutine N (running/created by)` -- the creation stack trace of the first goroutine.
- `Goroutine M (running/created by)` -- the creation stack trace of the second goroutine.

## Under the hood

The race detector is based on ThreadSanitizer (TSan), the same technology used in Clang and GCC. TSan instruments each memory access with a shadow word (two 32-bit cells per 8 bytes of application memory). Each shadow cell stores a thread ID, a clock vector timestamp, and an access type. When a memory access occurs, TSan checks its shadow for conflicting accesses with no happens-before relation. The happens-before graph is maintained via vector clocks that advance on lock acquire/release, channel send/receive, and `sync/atomic` operations.

False positives are theoretically possible with custom assembly or `unsafe.Pointer` games, but they are extremely rare in practice. False negatives are common: if your test does not exercise the race window, the detector will not find it.

## How Go uses it

The Go project itself runs `go test -race` on every package in the standard library as part of its build dashboard. Many open-source Go projects run `-race` tests nightly or on every PR. The `go` command also runs a limited race detector internally in the runtime's test suite.

## Go example

```go
package main

import (
	"fmt"
	"sync"
)

type SharedCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SharedCounter) Add(n int) {
	c.mu.Lock()
	c.value += n
	c.mu.Unlock()
}

func (c *SharedCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func racyAdd(c *SharedCounter, n int) {
	c.value += n // deliberate race for demonstration
}

func main() {
	var wg sync.WaitGroup
	safe := &SharedCounter{}
	n := 1000

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			safe.Add(1)
		}()
	}
	wg.Wait()

	fmt.Printf("Safe counter: %d (expected %d)\n", safe.Value(), n)

	fmt.Println("\nRun with -race to verify no races in the safe version.")
	fmt.Println("Try adding racyAdd calls to trigger a report.")
}
```

Run with:

```bash
go run -race .
```

If you add a call to `racyAdd`, the output includes:

```
WARNING: DATA RACE
Write at 0x... by goroutine 7:
  main.racyAdd(...)
  main.main.func1(...)
Previous write at 0x... by goroutine 5:
  main.(*SharedCounter).Add(...)
  main.main.func1(...)
```

## Step-by-step execution

When `go test -race` runs:

1. The compiler wraps every 1-, 2-, 4-, and 8-byte memory access with a TSan callback.
2. The test starts and goroutines begin executing.
3. When goroutine A reads `c.value` (no lock held), TSan records the read with a timestamp.
4. When goroutine B writes `c.value` (no lock held), TSan checks A's shadow entry, sees the previous read with no happens-before, and fires the report.
5. TSan prints the report to stderr and continues. After all tests finish, `go test` exits with code 1 because a race was detected.

## Common mistakes

- Mistake: Running `-race` on a single test case and declaring the code race-free.
  - Why it happens: The race detector only detects races that occur during the run. A test with 100 goroutines might not trigger the race, but 1000 goroutines in production will.
  - Fix: Run `-race` on every test and use the `-count=1` flag to disable test caching. Consider a stress flag (`-race -count=5`).

- Mistake: Ignoring the race report because "the program works fine in production".
  - Why it happens: The race might only manifest under specific scheduler interleavings that are rare on the current hardware.
  - Fix: A data race is always undefined behavior. Fix it regardless of observed symptoms.

- Mistake: Shipping a `-race` binary to production to monitor for races.
  - Why it happens: The developer wants continuous race detection.
  - Fix: The `-race` binary is too slow (5-10x) and memory-heavy (2-5x) for production. Use staging or canary deployments with `-race`, or run integration tests with `-race` in CI.

- Mistake: Using `sync/atomic` for one access and a mutex for another on the same variable.
  - Why it happens: The developer thinks partial synchronization is sufficient.
  - Fix: All accesses to a shared variable must use the same synchronization mechanism.

## Debugging walkthrough

Buggy program:

```go
package main

var counter int

func main() {
	go func() {
		for i := 0; i < 1000; i++ {
			counter++
		}
	}()
	for i := 0; i < 1000; i++ {
		counter++
	}
}
```

Run with `go run -race .`:

```
WARNING: DATA RACE
Write at 0x... by goroutine 5:
  main.main.func1()
Previous write at 0x... by main goroutine:
  main.main()
```

Fix: both goroutines must use a mutex or atomic:

```go
var mu sync.Mutex

func add() {
	mu.Lock()
	counter++
	mu.Unlock()
}
```

## Production notes

Integrate the race detector in CI by running all unit and integration tests with `-race`. Use a separate CI step that runs `go test -race ./...` with a timeout multiplier (since -race is slower). For services written in Go, run canary instances with `-race` in staging to catch races that only appear under realistic traffic patterns. The Go team recommends running the detector during local development, in pre-submit CI, and in post-submit CI with a stress factor.

## Performance implications

The race detector adds 5-10x CPU overhead and 2-5x memory overhead. Garbage collection latency increases because TSan shadow memory must be scanned. Never ship a `-race` binary to production. The overhead is acceptable for testing: a test suite that takes 30 seconds normally may take 3-4 minutes with `-race`. Use `-count=1` to disable test caching when running with `-race`.

## Practice task

Write a test file for `SharedCounter` that intentionally introduces a race by calling `racyAdd` from one goroutine and `Add` from another. Run `go test -race .` and confirm the race report. Then fix the test by removing the racy call and verify the race detector passes. Then add a table-driven test that runs the safe version with increasing goroutine counts and verify all pass with `-race`.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/22-race-detector
go test -race ./curriculum/modules/11-lifecycle-context-concurrency/lessons/22-race-detector
```

## Review questions

1. What are the two main performance costs of enabling `-race`?
2. Can the race detector produce false positives? If so, when?
3. Why does `-race` not guarantee that all races are found?
4. What happens when `go test -race` detects a race? How does it affect the exit code?
5. Name three CI strategies for integrating the race detector.

## NEXT UP

Goroutine leaks
