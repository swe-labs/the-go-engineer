# 09 Architecture & Security

## Mission

Learn the engineering of large-scale Go systems. This section turns correct code into maintainable architecture by applying naming discipline, visibility rules, package boundaries, architecture patterns, and security controls.

## Section Map

| Track | Surface | Mission |
| --- | --- | --- |
| `PD.1-PD.3` | [Package Design](./01-package-design) | teach naming, visibility, and layout rules |
| `ARCH.1-ARCH.9` | [Architecture Patterns](./03-architecture-patterns) | teach service layering, DDD, CQRS, and decoupling |
| `SEC.1-SEC.11` | [Security](./04-security) | teach input validation, cryptographic safety, identity, and trust boundaries |

[gRPC Reference](./02-grpc) is supporting reference material for contract-first service boundaries. It is not a canonical public track in Section 09.

## Why This Section Matters

You now know how to write Go code that is fast and correct. In a team environment, the shape of the codebase matters as much as the logic inside each function.

1. **Maintainability**: Package design makes ownership, dependencies, and import boundaries visible.
2. **Evolution**: Architecture patterns keep business rules separate from transport, persistence, and infrastructure details.
3. **Trust**: Security controls protect user data and make failure modes explicit at system boundaries.

## Suggested Learning Flow

1. **Start with Track PD**: Learn how Go packages work before you design larger systems.
2. **Continue to Track ARCH**: Apply those package rules to service boundaries and application architecture.
3. **Finish with Track SEC**: Integrate security into validation, authentication, authorization, transport, and secret-handling boundaries.

## Next Step

After learning architecture and security, continue to [Section 10 Production Operations](../10-production) or return to the [Curriculum Overview](../README.md).
