# Opslane deployment

## Mission

Understand and apply Opslane deployment in the context of professional Go software engineering.

## Prerequisites

- opslane-13

## Mental Model

The deployment pipeline is the application's delivery conveyor belt — it takes raw code from a repository and transforms it into a running, tested, observable production service through a series of automated stages.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

Under the hood, Docker multi-stage build copies only the compiled binary to the final image. Health checks are HTTP endpoints returning 200 OK — Kubernetes uses them via livenessProbe and readinessProbe. Rolling updates replace pods one at a time, waiting for readiness checks between each replacement.

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/14-opslane-deployment
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Building the binary manually on the production server instead of using CI/CD.
- Tagging Docker images with :latest instead of versioned tags.
- Deploying without health checks or rollback capability.

## In Production

Docker and GitHub Actions (or GitLab CI) are the industry standard for Go service deployment. Kubernetes orchestrates containerized Go services at scale — used by Google, Netflix, Shopify, and thousands of companies.

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-15`.
