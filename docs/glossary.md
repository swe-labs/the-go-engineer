# Glossary

- **Proof Surface**: A runnable or checkable artifact that verifies a claim (tests, validator output, CI checks, scripts).
- **Curriculum Registry**: The machine-readable course map in `curriculum.v2.json`.
- **Architecture Lock**: Contract that public root sections stay fixed unless maintainers explicitly approve architecture work.
- **Teaching Stub**: Intentional simplified implementation that demonstrates concepts but is not fully production-grade.
- **Production-Shaped**: Code that reflects real system boundaries and operational concerns, even when intentionally scoped.
- **Dirty Migration State**: A migration state requiring operator intervention before additional schema changes are applied.
- **Idempotency**: Safe retry behavior where repeated requests do not create duplicate side effects.
- **Outbox Pattern**: Reliable event publication strategy that stores events in the primary DB transaction and dispatches asynchronously.
