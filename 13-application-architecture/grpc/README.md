# Track D: gRPC Architecture Reference

## Mission

This surface shows how schema-first service boundaries look in Go when transport contracts matter
as much as handler code.

It is intentionally treated as an architecture reference surface in beta, not as the main proof
path for the stage.

## Beta Stage Ownership

This track belongs to [09 Architecture](../../docs/stages/09-architecture.md).

Within the beta public shell, it is a reference surface for service contracts, generated boundary
code, and RPC-style architecture after the package-design path is already clear.

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

## How To Use It In Beta

1. complete the live `package-design` path first
2. read this surface when you want stronger intuition for service boundaries
3. treat it as architectural reinforcement, not as a required milestone before continuing

## References

1. [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
2. [Protocol Buffers Language Guide](https://protobuf.dev/programming-guides/proto3/)
3. [google.golang.org/grpc](https://pkg.go.dev/google.golang.org/grpc)

## Next Step

After you use this surface, return to the
[Architecture stage](../../docs/stages/09-architecture.md)
or continue into
[10 Production](../../docs/stages/10-production.md)
if your current gap is operating systems safely.


