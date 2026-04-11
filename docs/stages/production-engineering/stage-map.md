# Production Engineering Stage Map

## Stage Goal

This stage teaches learners how to think about runtime visibility, shutdown safety, and deployment
workflow as part of engineering correctness.

## Public Stage Shape

| Surface | Role In Stage |
| --- | --- |
| [structured-logging](../../../14-application-architecture/structured-logging/) | live beta path for operational log shape and redaction |
| [graceful-shutdown](../../../14-application-architecture/graceful-shutdown/) | live beta path for drain order and lifecycle safety |
| [docker-and-deployment](../../../14-application-architecture/docker-and-deployment/) | reference surface for packaging and deployment workflow |

## Live Milestone Backbone

| ID | Surface | Why It Matters |
| --- | --- | --- |
| `SL.5` | PII redactor exercise | proves that log pipelines can enforce operational policy |
| `GS.3` | shutdown capstone | proves that service stop behavior can be coordinated deliberately |

## Recommended Order

1. `SL.1` through `SL.5`
2. `GS.1` through `GS.3`
3. `docker-and-deployment`

## Reference Surfaces That Still Matter

These are important production surfaces, but they are not the current public beta proof path:

- `docker-and-deployment`
