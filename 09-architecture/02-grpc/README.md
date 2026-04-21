# Track B: gRPC Architecture Reference

## Mission

This surface shows how schema-first service boundaries look in Go when transport contracts matter as much as handler code.

It remains a reference-heavy architecture track inside Stage 09, not the first proof path a learner should take.

## Stage Ownership

This track belongs to [09 Architecture & Security](../README.md).

## Why This Surface Matters

gRPC is useful here because it forces architecture decisions into explicit contracts:

- messages and services are schema-first
- client and server boundaries are typed
- transport seams become visible and reviewable
- code generation supports the contract instead of replacing design thinking

## Current Surface Shape

| Area | Focus |
| --- | --- |
| `proto/` | service and message definitions |
| `gen/` | generated Go code committed for runnable examples |
| `1-unary/` | basic request/response RPC |
| `2-streaming/` | streaming flows and interceptors |

## Next Step

After you use this surface, return to the [09 Architecture & Security overview](../README.md) or continue to [10 Production Operations](../../10-production).
