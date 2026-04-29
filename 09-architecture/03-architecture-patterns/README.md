# Track ARCH: Architecture Patterns

## Mission

Master the major service-shape and layering decisions that teams face as systems grow. This track covers high-level concepts like Domain-Driven Design (DDD), Hexagonal Architecture, and Event-Driven systems, providing a toolkit for building scalable and maintainable Go applications.

## Stage Ownership

This track belongs to [09 Architecture & Security](../README.md).

## Track Map

| ID | Type | Surface | Mission | Requires |
| --- | --- | --- | --- | --- |
| `ARCH.1` | Lesson | [Architecture Trade-offs](./1-architecture-trade-offs) | Choose the right system shape (Monolith vs. Micro). | entry |
| `ARCH.2` | Lesson | [DDD Basics](./2-ddd-basics) | Manage business complexity with Bounded Contexts. | `ARCH.1` |
| `ARCH.3` | Lesson | [Hexagonal Architecture](./3-hexagonal-architecture-in-go) | Isolate domain logic from technical details. | `ARCH.2` |
| `ARCH.4` | Lesson | [Repository Pattern](./4-repository-pattern-deep-dive) | Master data persistence boundaries. | `ARCH.3` |
| `ARCH.5` | Lesson | [Service Layer Pattern](./5-service-layer-pattern) | Coordinate complex business use cases. | `ARCH.4` |
| `ARCH.6` | Lesson | [Event-Driven Architecture](./6-event-driven-architecture) | Decouple systems using facts (Events). | `ARCH.5` |
| `ARCH.7` | Lesson | [CQRS Basics](./7-cqrs-basics) | Scale reads and writes independently. | `ARCH.6` |
| `ARCH.8` | Lesson | [When to Split Services](./8-when-to-split-services) | Identify the tipping point for Microservices. | `ARCH.7` |
| `ARCH.9` | Exercise | [Modular Refactor](./9-modular-refactor-exercise) | Practice untangling a monolith. | `ARCH.1-8` |

## Next Step

After completing this track, continue to [Track SEC: Security](../04-security) or return to the [09 Architecture & Security overview](../README.md).
