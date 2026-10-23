# Opslane Milestone 14 Assessment — Deployment: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** What is the distinct operational purpose of a Kubernetes Liveness probe versus a Readiness probe?
- a) Liveness probes run on Linux; readiness probes run on Windows.
- b) A Liveness probe determines if the process is alive (restarting the container on failure); a Readiness probe determines if the container can accept incoming traffic (removing the pod from load balancer endpoints on failure without restarting it).
- c) Liveness probes check CPU usage; readiness probes check disk space.
- d) Both probes trigger immediate pod restarts on failure.

**2.** Why must `/health/live` NEVER query the database or external cache?
- a) Database queries in Go cannot return HTTP status codes.
- b) If the database experiences a brief 5-second connection saturation, all application liveness probes will fail simultaneously, causing Kubernetes to restart every pod across the cluster in a cascading outage.
- c) Kubernetes does not allow database drivers inside containers.
- d) Liveness probes only accept plaintext responses.

**3.** What HTTP status code does `/health/ready` return when one or more dependency checks fail?
- a) HTTP 200 OK.
- b) HTTP 503 Service Unavailable.
- c) HTTP 500 Internal Server Error.
- d) HTTP 404 Not Found.

**4.** Why does `handleReadiness` execute component health checks concurrently using goroutines?
- a) To ensure the total readiness evaluation duration is bounded by the slowest individual component check, rather than accumulating latency sequentially.
- b) Because Go forbids calling more than one function sequentially.
- c) To bypass CPU thread scheduling.
- d) To write to multiple database tables simultaneously.

**5.** In `handleReadiness`, how is context timeout enforced across component health checks?
- a) Via `time.Sleep`.
- b) By deriving `ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)` and passing `ctx` to all checkers.
- c) Through operating system cron jobs.
- d) Using HTTP 1.1 keep-alive timeouts.

**6.** In `handleReadiness`, how do concurrent goroutines store their results safely without race conditions?
- a) By appending to a shared slice with a global mutex lock.
- b) By preallocating the slice `results := make([]ComponentHealth, len(s.checkers))` and assigning each goroutine its own unique index `results[idx] = ...`.
- c) By writing results to disk.
- d) By sending results over an unbuffered channel.

---

## Short Answer (2 points each)

**7.** Explain what happens during a Kubernetes rolling update when readiness probes are configured properly versus when they are omitted.

**8.** What is a Cascading Failure in microservices? How does an improperly designed liveness probe cause a cascading cluster reboot?

**9.** Explain the significance of `initialDelaySeconds`, `periodSeconds`, `timeoutSeconds`, and `failureThreshold` in Kubernetes container probe definitions.

**10.** In a microservice deployment with dependencies on PostgreSQL, Redis, and a background Worker Pool, which of these dependencies should affect the Liveness probe, which should affect the Readiness probe, and why?
