# Opslane Milestone 14 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) A Liveness probe determines if the process is alive (restarting the container on failure); a Readiness probe determines if the container can accept incoming traffic (removing the pod from load balancer endpoints on failure without restarting it).**  
*Explanation:* Under Kubernetes pod lifecycle design, liveness failures trigger container restarts (`SIGKILL` followed by fresh boot). Readiness failures do NOT restart the container; they remove the pod from Service endpoints so it stops receiving traffic until healthy.

**2. b) If the database experiences a brief 5-second connection saturation, all application liveness probes will fail simultaneously, causing Kubernetes to restart every pod across the cluster in a cascading outage.**  
*Explanation:* Liveness probes must only verify that the Go process is alive and un-deadlocked. If liveness depends on an external database and the database stutters, every pod in the deployment will fail liveness and restart simultaneously, destroying all active connections and exacerbating the database overload.

**3. b) HTTP 503 Service Unavailable.**  
*Explanation:* HTTP 503 is the canonical HTTP status code indicating that the server is temporarily unable to handle requests due to overloaded or degraded downstream dependencies.

**4. a) To ensure the total readiness evaluation duration is bounded by the slowest individual component check, rather than accumulating latency sequentially.**  
*Explanation:* If checking the DB takes 30ms, Redis takes 20ms, and the worker pool takes 15ms: sequential execution takes $30+20+15 = 65\text{ms}$ (exceeding a 50ms deadline). Concurrent execution completes in $\max(30, 20, 15) = 30\text{ms}$.

**5. b) By deriving `ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)` and passing `ctx` to all checkers.**  
*Explanation:* Deriving a deadline context bounded to 50ms ensures that if any single dependency checker hangs, its context is cancelled and the overall readiness check completes within the timeout.

**6. b) By preallocating the slice `results := make([]ComponentHealth, len(s.checkers))` and assigning each goroutine its own unique index `results[idx] = ...`.**  
*Explanation:* Slices in Go are safe for concurrent writes if each goroutine writes to a distinct, preallocated index (`results[0]`, `results[1]`). Appending (`results = append(...)`) would cause data races because slice headers and lengths are mutated concurrently.

---

## Short Answer (2 points each)

**7. Rolling updates with and without readiness probes:**
- **Without Readiness Probes:** Kubernetes assumes a newly started pod is ready as soon as the container process starts. It immediately terminates an old pod and routes live customer traffic to the new pod before the new pod has established database pools or warmed up caches. Users experience dropped requests, connection refused errors, and cold-start latency spikes.
- **With Readiness Probes:** Kubernetes waits until the new pod's readiness probe returns HTTP 200 OK. Traffic is only routed to the new pod once it has verified all its dependencies, and old pods are only terminated once the new pod is confirmed healthy, guaranteeing true zero-downtime deployments.
*Grading: 1 pt for explaining premature traffic routing without readiness; 1 pt for explaining smooth zero-downtime transition with readiness.*

**8. Cascading failure and liveness probes:**
- **Cascading Failure:** When a minor failure in one component triggers failures across interconnected components, escalating into a complete cluster outage.
- **Liveness Trap:** If `/health/live` queries the database, a momentary database slowdown causes all app pods to fail liveness simultaneously. Kubernetes restarts all 20 pods at once. As the 20 pods restart, they all execute startup connection handshakes, run migrations, and flood the struggling database with hundreds of new connection attempts simultaneously. The database crashes completely, keeping the pods in an infinite crash-loop (CrashLoopBackOff).
*Grading: 1 pt for defining cascading failure; 1 pt for explaining the infinite restart / database connection flood storm.*

**9. Kubernetes probe tuning parameters:**
- `initialDelaySeconds`: How long the kubelet waits after container creation before polling the first probe (allows slow app startup / migrations to finish).
- `periodSeconds`: The frequency (interval) between subsequent probe polls (e.g. poll every 10s).
- `timeoutSeconds`: How long the kubelet waits for a response before treating the probe call as a failure (e.g. timeout after 2s).
- `failureThreshold`: How many consecutive probe failures must occur before Kubernetes takes action (restarting for liveness, removing from endpoints for readiness).
*Grading: 1 pt for defining initial delay and period; 1 pt for defining timeout and failure threshold.*

**10. Dependency assignment to Liveness vs Readiness:**
- **Liveness Probe:** None of the external dependencies! Liveness should only monitor internal process vitality (e.g. checking that the internal HTTP event loop is un-deadlocked).
- **Readiness Probe:** PostgreSQL, Redis, and Worker Pool! If PostgreSQL or Redis is down, or if the worker pool queue is completely jammed, the application cannot successfully serve customer requests and should report `503 Service Unavailable` so the load balancer routes requests to other healthy replicas.
*Grading: 1 pt for correctly isolating liveness from external dependencies; 1 pt for correctly assigning DB/Cache/Workers to readiness with operational justification.*
