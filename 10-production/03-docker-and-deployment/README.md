# Track E: Docker and Deployment Reference

## Mission

This surface shows how Go applications move from local execution into containerized deployment workflows.

It remains a production-oriented reference track inside Stage 10, not the first milestone path for the section.

## Stage Ownership

This track belongs to [10 Production Operations](../README.md).

## Why This Surface Matters

Deployment changes the shape of a system:

- packaging decisions affect startup, security, and rebuild speed
- container structure changes how applications are configured and run
- rollout workflow exposes whether shutdown and observability decisions were actually sound

## Current Surface Shape

| Area | Focus |
| --- | --- |
| `1-docker-basics/` | basic Dockerfile structure |
| `2-multi-stage-builds/` | builder/runtime separation for Go binaries |
| `3-docker-compose/` | multi-service development setup |
| `4-cicd-pipelines/` | delivery workflow |
| `5-blue-green-and-rollback/` | rollout and recovery patterns |
| `6-dockerised-service/` | integrated deployment proof surface |

## Next Step

After you use this surface, return to the [10 Production Operations overview](../README.md) or continue into [11 Flagship](../../11-flagship).
