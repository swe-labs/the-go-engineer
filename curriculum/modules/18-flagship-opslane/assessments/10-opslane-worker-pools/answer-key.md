# Opslane Milestone 10 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) Unbounded goroutine creation under high traffic exhausts system memory, increases garbage collector scan times, and degrades CPU cache locality.**  
*Explanation:* While goroutines are lightweight (~2–4KB per stack), 100,000 active goroutines consume hundreds of megabytes of RAM. As goroutines churn through heap objects, GC cycles trigger latency spikes. Bounded worker pools constrain concurrency to an optimal number matching hardware capacity.

**2. b) When reading from `p.jobs` returns `job, ok := <-p.jobs` where `ok == false`, indicating that `close(p.jobs)` was called and the buffer has drained.**  
*Explanation:* In Go, reading from a closed channel does not block; it yields the remaining buffered values followed by the zero value and `ok == false`. This is the canonical idiom for worker termination.

**3. b) An atomic Compare-And-Swap (CAS) loop with `p.maxConcurrencySeen.CompareAndSwap(currentMax, active)`.**  
*Explanation:* In a CAS loop, a thread reads the current value, computes the desired update, and attempts an atomic swap. If another thread updated the value in the interim, the swap fails and the loop re-evaluates. This avoids mutex locking overhead.

**4. b) A deferred `recover()` block logs the panic, decrements active workers, increments completed jobs, and allows the worker goroutine to process subsequent jobs safely.**  
*Explanation:* Because panics are caught within `defer` inside `executeJob`, the worker goroutine does not terminate; it logs the error, updates telemetry counters, and loops to pick up the next available job.

**5. b) `ErrPoolClosed`.**  
*Explanation:* Once `Shutdown` is called, `closed.Swap(true)` causes `Submit` to immediately return `ErrPoolClosed` without enqueuing.

**6. b) Coordinates shutdown: each worker calls `defer p.wg.Done()`, and `Shutdown` waits on `p.wg.Wait()` until all workers have completely drained and terminated.**  
*Explanation:* `sync.WaitGroup` tracks the active lifecycle of the $N$ worker goroutines, ensuring `Shutdown()` only returns after all workers have finished processing their final jobs.

---

## Short Answer (2 points each)

**7. Bounded vs. Unbounded channels:**
- **Bounded Channels:** Created with an explicit capacity (e.g. `make(chan Job, 50)`). Holds up to 50 items in an internal ring buffer. When the buffer is full, sending on the channel blocks the sender until a worker reads an item.
- **Unbounded Behavior:** Go does not have native unbounded channels. Attempting to buffer indefinitely in memory without bounds eventually consumes all available RAM and triggers an OS OOM kill.
*Grading: 1 pt for explaining buffered capacity and sender blocking; 1 pt for explaining memory exhaustion risks.*

**8. Backpressure in worker pools:**
- **Definition:** A mechanism that slows down the ingestion rate of producers when consumers are saturated, preventing system overload.
- **Worker Pool Application:** When all worker goroutines are busy and the `jobs` channel buffer is full, subsequent calls to `Submit(j)` block. If an upstream HTTP handler calls `Submit()`, the HTTP handler's response is delayed, naturally throttling external clients or returning HTTP 503/429 rather than taking down the server.
*Grading: 1 pt for defining backpressure; 1 pt for explaining how channel blocking throttles upstream callers.*

**9. Mechanics of Compare-And-Swap (CAS):**
A CAS operation is a hardware-supported atomic CPU instruction (`CMPXCHG` on x86). It compares the contents of a memory location with a given value and, only if they are equal, modifies the contents to a new value. If another core modified the variable in the interim, the operation fails and the loop retries. It is significantly faster than mutex locks because it requires zero OS thread context switches and zero kernel sleep transitions.
*Grading: 1 pt for explaining compare-and-swap CPU instruction; 1 pt for highlighting performance advantages over mutexes (no context switching).*

**10. Channel closing rules and race prevention in `Shutdown`:**
- **Go Rule:** Sending to a closed channel causes an immediate, unrecoverable runtime panic (`panic: send on closed channel`).
- **Race Prevention:** Setting `p.closed.Swap(true)` first ensures that any subsequent or concurrent calls to `Submit()` will observe `closed == true` and abort immediately with `ErrPoolClosed`, rather than attempting to send into `p.jobs` right as `close(p.jobs)` is executed.
*Grading: 1 pt for identifying the `send on closed channel` panic; 1 pt for explaining how `closed.Swap(true)` gates sends before closing.*
