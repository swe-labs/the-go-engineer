# Opslane Milestone 10 Assessment — Worker Pools: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** Why does `BoundedWorkerPool` constrain concurrency to a fixed number of workers rather than spawning a new goroutine for every job?
- a) Go compilers cannot execute more than 16 goroutines simultaneously.
- b) Unbounded goroutine creation under high traffic exhausts system memory, increases garbage collector scan times, and degrades CPU cache locality.
- c) Channels cannot be read by more than one goroutine.
- d) To disable operating system thread scheduling.

**2.** How does a worker goroutine know when all jobs have been submitted and it should exit?
- a) When the operating system sends `SIGTERM`.
- b) When reading from `p.jobs` returns `job, ok := <-p.jobs` where `ok == false`, indicating that `close(p.jobs)` was called and the buffer has drained.
- c) When the memory usage reaches 100%.
- d) Workers never terminate.

**3.** What technique does `executeJob` use to update `maxConcurrencySeen` without a mutex lock?
- a) Channel multiplexing.
- b) An atomic Compare-And-Swap (CAS) loop with `p.maxConcurrencySeen.CompareAndSwap(currentMax, active)`.
- c) AES-256 encryption.
- d) A global spinlock.

**4.** What happens if a job panics during `j.Execute(p.ctx)`?
- a) The entire Go program crashes immediately.
- b) A deferred `recover()` block logs the panic, decrements active workers, increments completed jobs, and allows the worker goroutine to process subsequent jobs safely.
- c) The channel is automatically closed.
- d) The database restarts.

**5.** If a client calls `Submit(j)` after `Shutdown` has been initiated, what does `Submit` return?
- a) `nil`.
- b) `ErrPoolClosed`.
- c) An HTTP 500 error.
- d) It blocks forever.

**6.** What does `sync.WaitGroup` do in `BoundedWorkerPool`?
- a) Limits the number of items in the channel.
- b) Coordinates shutdown: each worker calls `defer p.wg.Done()`, and `Shutdown` waits on `p.wg.Wait()` until all workers have completely drained and terminated.
- c) Encrypts data in transit.
- d) Calculates job completion duration.

---

## Short Answer (2 points each)

**7.** Explain the difference between bounded channels (`make(chan Job, cap)`) and unbounded channels. What happens when a bounded channel's buffer becomes completely full?

**8.** What is Backpressure in software engineering? How does a bounded worker pool apply backpressure to upstream HTTP callers during peak traffic?

**9.** Explain the mechanics of a Compare-And-Swap (CAS) loop. Why is CAS faster than mutex locking for tracking peak numeric telemetry?

**10.** Explain the graceful shutdown sequence of `BoundedWorkerPool`: why is `p.closed.Swap(true)` called before `close(p.jobs)`? What race condition would occur if `p.jobs` was closed while producers were still executing `Submit`?
