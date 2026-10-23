# Questions — Module 16 Checkpoint: Portfolio & Interview Readiness

Answer each question with concrete reasoning and evidence.

---

### Question 1 (Multiple Choice)
When curating a senior software engineering portfolio on GitHub, what is the single most critical factor that distinguishes high-impact projects from generic tutorial clones?
- A) The repository having thousands of stars from automated bots.
- B) An extensive commit history spanning over ten years.
- C) Verifiable proof of production engineering: detailed architectural contracts, comprehensive automated test suites including race detection, reproducible benchmarks, and operational artifacts (Docker, CI/CD).
- D) Using as many external third-party libraries and frameworks as possible.

---

### Question 2 (Multiple Choice)
During a live technical coding interview in Go, an interviewer asks you to write a concurrent map processing function. Before writing any code, what should you do first?
- A) Immediately begin typing `func ProcessMap(m map[string]int)`.
- B) Ask clarifying questions to establish requirements: concurrency bounds, read/write ratios, memory constraints, error handling expectations, and define the function signature and expected edge cases.
- C) Complain that Go maps are not concurrently safe and insist on switching languages.
- D) Memorize a leetcode solution silently for 15 minutes without speaking.

---

### Question 3 (Multiple Choice)
When evaluating a take-home coding assignment for a senior backend position, what is typically the primary reason hiring teams fail a candidate whose solution produces correct output?
- A) The code was written in Go 1.22 instead of Go 1.21.
- B) Absence of automated tests, unstructured error handling (`panic` or ignored errors), lack of context cancellation, data races under `go test -race`, or lack of architectural documentation explaining trade-offs.
- C) The candidate provided too much documentation in the README.
- D) The candidate used table-driven tests instead of assertions.

---

### Question 4 (Multiple Choice)
In a behavioral interview using the STAR method, what is the primary objective of the "Result" (R) section?
- A) To describe how tired you were after the incident was resolved.
- B) To quantify the concrete impact of your actions using measurable metrics (e.g., "reduced P99 latency by 45%, eliminated deadlocks, zero customer data lost") and reflect on lessons learned.
- C) To assign blame to junior engineers or third-party vendors for the initial outage.
- D) To explain the syntax of the programming language used.

---

### Question 5 (Multiple Choice)
When writing a technical blog post or system design case study for public engineering visibility, why is it valuable to include a "What Went Wrong / Trade-offs" section?
- A) It shows potential employers that your code is unreliable.
- B) It demonstrates intellectual honesty, systems maturity, and deep architectural understanding of operational failure modes, which hiring managers value far more than simplistic "happy path" narratives.
- C) It is required by open-source licensing agreements.
- D) It fills word count requirements for SEO indexing.

---

### Question 6 (Multiple Choice)
When receiving an initial job offer, why should you avoid immediately accepting the first verbal compensation figure presented by the recruiter?
- A) Recruiters intentionally revoke offers if candidates accept immediately.
- B) Initial offers typically represent the midpoint or lower end of the approved salary band; expressing enthusiasm while professionally asking for time to review the complete package allows for benchmarking and value-based counter-offers.
- C) Labor laws require at least three negotiation rounds before an offer can be signed.
- D) It is considered disrespectful to accept without asking for a company car.

---

### Question 7 (System Design Interview Walkthrough)
You are in a 45-minute system design interview for a Senior Distributed Systems Engineer role. The prompt is:
*"Design a globally distributed rate limiting and API quota service that handles 500,000 requests per second across 4 geographic regions."*
Provide your complete structured interview response:
1. Requirements gathering (Functional and Non-Functional).
2. High-level architecture diagram and component breakdown.
3. Data storage and synchronization strategy (Redis clusters, local memory token buckets with async batching, CAP trade-offs).
4. Failure modes and blast radius containment (what happens if the central Redis cluster becomes unreachable?).

---

### Question 8 (Live Coding Interview Simulation)
Write a complete, thread-safe, generic in-memory Cache in Go that satisfies senior interview expectations:
1. Generic key and value types (`[K comparable, V any]`).
2. Set with TTL and Get methods.
3. Background expiration loop using `time.Ticker` that cleanly shuts down via `context.Context`.
4. RWMutex synchronization to optimize for read concurrency.
5. Provide a table-driven test and benchmark verifying thread safety.

---

### Question 9 (Incident Response & Postmortem Defense)
You are asked in an interview:
*"Describe an incident where a memory leak or goroutine leak in your Go service brought down production. How did you identify it, mitigate it, and prevent it from recurring?"*
Provide a detailed response following the STAR framework:
1. **Situation**: Service architecture, traffic profile, and symptoms (OOM kills, alerts).
2. **Task**: Your specific responsibility during the active incident.
3. **Action**: The investigative steps taken (heap profile analysis via `pprof`, goroutine dump inspection, identifying the unbuffered channel / missing context cancellation).
4. **Result**: Recovery MTTR, memory footprint reduction, and permanent guardrails instituted (CI race/leak tests, alerting thresholds).

---

### Question 10 (Senior Career & Offer Negotiation Strategy)
Explain the mechanics and strategy of negotiating a Senior/Staff Engineer compensation package:
1. Explain the difference between Restricted Stock Units (RSUs) at a publicly traded company vs Incentive Stock Options (ISOs) at an early-stage startup.
2. What questions must an engineer ask about a startup's equity (preferred share liquidation preferences, strike price, current 409A valuation, total outstanding diluted shares)?
3. Draft a script for countering an offer professionally using market benchmarking data and technical capability alignment.
