# Opslane Threat Model

## Goal

Document the minimum threat model for Opslane so learners can connect implementation choices to concrete abuse cases.

## Assets

- tenant-scoped user identity and tokens
- order/payment workflow integrity
- database confidentiality and durability
- service availability under abusive traffic

## Trust Boundaries

- internet client to HTTP API boundary
- API to PostgreSQL boundary
- API to background worker boundary
- optional observability export boundary (OTLP endpoint)

## Main Threats

1. **Auth token abuse**
   - stolen tokens reused across tenants
   - expired/forged token acceptance
2. **Tenant isolation failure**
   - cross-tenant reads/writes via ID confusion
3. **Rate-limit bypass**
   - proxy header spoofing or shared-key evasion
4. **Workflow replay/idempotency abuse**
   - duplicate submissions creating inconsistent financial state
5. **Operational visibility gaps**
   - missing traces/metrics hide incident root cause

## Existing Controls

- signed auth tokens and identity middleware
- tenant-scoped repository queries
- distributed rate-limit support with trusted proxy parsing
- structured request logging and correlation IDs
- CI vulnerability scan (`govulncheck`) and race tests

## Gaps To Continue Hardening

- refresh token rotation and revocation policy
- centralized audit log policy for auth-sensitive actions
- stronger role/permission model
- full distributed trace propagation and exporter verification
- periodic dependency and static-analysis automation

## Abuse-Case Checklist For New Features

For each new endpoint or workflow:

1. What prevents cross-tenant access?
2. Is the operation idempotent for retries?
3. What limits abusive request volume?
4. What gets logged for incident reconstruction?
5. Which failure mode is fail-open vs fail-closed, and why?
