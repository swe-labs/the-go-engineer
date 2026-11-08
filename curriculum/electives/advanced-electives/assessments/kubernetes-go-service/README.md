# Project Assessment — Kubernetes Go Service Capstone

## Purpose

Validate that the learner can deploy, operate, roll back, and explain a Go service on Kubernetes without relying on magic. The capstone is the proof that the Kubernetes stream's mechanisms landed, not just its vocabulary.

## Targets

- `project-kubernetes-go-service`

## Evidence

Submit the deployment manifests, the hardened image reference, the captured `kubectl rollout status` output, the `/version` response after a rollout and after a rollback, the probe results, the node-drain output, and the written runbook extension plus defense.

---

## 1. Visual Model: Stream E Verification Pipeline

```
+-----------------------------------------------------------------------------------------+
|                     KUBERNETES GO SERVICE CAPSTONE PIPELINE                            |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|   Deploy                Verify                          Recover                         |
|   +----------------+    +--------------------------+     +----------------------+        |
|   | kind/k3d       |    | kubectl rollout status   |     | kubectl rollout undo |        |
|   | image loaded   |--->| GET /version (revision)  |---> | previous revision    |        |
|   | objects applied|    | probes / resources /sec  |     | PDB holds minimum    |        |
|   +----------------+    +--------------------------+     +----------------------+        |
|              |                       |                          |                       |
|              v                       v                          v                       |
|   +----------------+    +--------------------------+     +----------------------+        |
|   | ConfigMap      |    | security context proof   |     | top 3 failure modes  |        |
|   | Secret (file)  |    | read-only, non-root      |     | symptom -> tool ->   |        |
|   | probes configured|  | resource limits set      |     | cause -> fix         |        |
|   +----------------+    +--------------------------+     +----------------------+        |
+-----------------------------------------------------------------------------------------+
```

Caption: the pass condition is not "the containers started." It is that a rollout was verified through the Service, a rollback restored the previous revision, and a node drain held the minimum availability.

---

## 2. Rubric

| Criterion | Weight | What earns full marks |
| :--- | :---: | :--- |
| **Correctness** | **35%** | The service runs in the cluster, traffic reaches it, and the rollout, rollback, and drain all complete as specified. |
| **Production readiness** | **25%** | Probes, resource requests and limits, the hardened security context, and the ConfigMap/Secret split are all present and correct. |
| **Proof** | **25%** | `rollout status`, `/version`, probe, and drain evidence are captured and interpreted, not just pasted. |
| **Explanation** | **15%** | A clear defense of each object, the tested failure modes, and the tradeoffs. |

---

## 3. Passing Bar

- **Passing score:** 80%.
- Each criterion is scored 0, 3, or 5 and weighted per the table. A zero on Correctness fails the assessment regardless of the total.
- The learner must demonstrate the rollback and the drain, not only describe them.

---

## 4. Defense Questions

1. Which revision is serving traffic right now, and which command proved it?
2. Why does `maxUnavailable: 0` require spare cluster capacity to produce zero downtime?
3. What does the PodDisruptionBudget protect against, and what does it not protect against?
4. You changed the schema and the code in the same release. Why is the rollback unsafe, and what should have happened instead?
5. Name the two metrics you would alert on for a Go service under a CPU limit, and why.

---

## 5. Retake Policy

The learner may retry after fixing the gaps identified by the rubric and writing a short correction note that explains what was wrong and what changed. The correction note is not scored but is required before a retake.

---

## NEXT UP

This assessment closes the Kubernetes stream. Return to the [Module 17 overview](../../README.md) to choose another stream or complete the elective track.
