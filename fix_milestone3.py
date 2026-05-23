"""Milestone 3: Add core-10-25 Refresh tokens, promote core-11-28 to full lesson."""
import json, copy

with open('curriculum/path.core.json', 'r', encoding='utf-8') as f:
    core = json.load(f)

items = core['items']
item_map = {it['id']: it for it in items}

# ════════════════════════════════════════════════════════════════
# 1. Add core-10-25: Refresh tokens
# ════════════════════════════════════════════════════════════════
print("=== Adding core-10-25 Refresh tokens ===")

# Check if it already exists
if 'core-10-25' in item_map:
    print("  core-10-25 already exists, skipping")
else:
    # Get last module-10 item for the chain fix
    m10_items = sorted(
        [it for it in items if it.get('module_id') == 'module-10'],
        key=lambda x: x.get('order', 0)
    )
    last_m10 = m10_items[-1] if m10_items else None
    print(f"  Last module-10 item: {last_m10['id']} ({last_m10['title']})")
    print(f"  Its next_item_ids currently: {last_m10.get('next_item_ids', [])}")

    # Create new item
    new_item = {
        "id": "core-10-25",
        "module_id": "module-10",
        "slug": "refresh-tokens",
        "title": "Refresh tokens",
        "type": "lesson",
        "subtype": "concept",
        "status": "planned",
        "difficulty": "intermediate",
        "phase": "backend",
        "order": 25,
        "estimated_minutes": 45,
        "learning_objective": "Implement a secure refresh token rotation flow in Go.",
        "required_prior_knowledge": ["Complete prerequisite lesson(s): core-10-24."],
        "prerequisites": ["core-10-24"],
        "next_item_ids": last_m10.get('next_item_ids', []),
        "zero_magic_status": "golden",
        "crossrefs": {
            "builds_on": [
                {
                    "target_id": "core-10-07",
                    "label": "Module 10.07 \u2014 JWT implementation",
                    "display_text": "Builds on: Module 10.07 \u2014 JWT implementation",
                    "reason": "Refresh tokens are built on top of JWT access token issuance.",
                    "required": True
                }
            ],
            "preview_only": [],
            "related": [],
            "reinforced_in": [
                {
                    "target_id": "project-secure-authenticated-api",
                    "label": "Project \u2014 Secure Authenticated API",
                    "display_text": "Reinforced in: Project \u2014 Secure Authenticated API",
                    "reason": "This project implements a production-grade auth system with refresh tokens.",
                    "required": True
                }
            ]
        },
        "proof": {
            "assessment_id": "assessment-module-10",
            "expected_artifact": "Runnable code demonstrating a refresh token rotation flow.",
            "mastery_checks": [
                "Can implement access token expiry with refresh token rotation.",
                "Can revoke refresh tokens and detect replay attacks.",
                "Can distinguish between access token and refresh token security requirements."
            ],
            "practice_task": "Implement a token service that issues access tokens (15min expiry) and refresh tokens (7d rotation). Include a revocation endpoint, test that old refresh tokens are invalidated on rotation, and verify concurrent rotation requests don't cause race conditions.",
            "project_id": "project-secure-authenticated-api",
            "rubric_ids": ["rubric-module-10"]
        },
        "content_contract": {
            "common_mistakes_required": True,
            "machine_view_required": True,
            "portfolio_artifact_required": True,
            "production_notes_required": True,
            "readme_required": True,
            "review_questions_required": True,
            "runnable_required": True,
            "tests_required": True,
            "visual_model_required": True
        },
        "verification": {
            "expected_output": "Token service with rotation, revocation, and race-condition tests.",
            "manual_steps": [
                "Read the README.",
                "Run the token service example.",
                "Complete the practice task.",
                "Answer the review questions."
            ],
            "mode": "mixed",
            "race_command": "go test -race ./...",
            "run_command": "go run ./10-security/25-refresh-tokens",
            "test_command": "go test ./10-security/25-refresh-tokens"
        },
        "zero_magic": {
            "beginner_mistakes": [
                "Storing refresh tokens in the same database as access tokens without isolation.",
                "Not rotating refresh tokens on each use, allowing unlimited reuse of stolen tokens.",
                "Using opaque bearer tokens without a revocation check, making logout impossible."
            ],
            "execution_timeline": [
                "Understand the access token expiry problem: short expiry forces re-login, long expiry is dangerous.",
                "Design the token rotation flow: issue access + refresh pair, rotate refresh on each use.",
                "Implement the token store with revocation support (database allowlist/denylist).",
                "Implement the refresh endpoint with rotation and race-condition protection.",
                "Add tests: concurrent rotation, token reuse detection, expiry enforcement."
            ],
            "failure_modes": [
                {
                    "scenario": "Two concurrent requests arrive with the same refresh token. Both succeed, producing two valid pairs. An attacker who stole the token can still access the system after the legitimate user's first refresh.",
                    "cause": "No concurrency protection on the rotation operation. The second request reads the old token before the first request's write completes.",
                    "fix": "Use a database transaction with SELECT FOR UPDATE or an optimistic concurrency check (token_version column incremented on each rotation)."
                },
                {
                    "scenario": "A revoked refresh token is still accepted because the revocation check reads from a stale cache.",
                    "cause": "The revocation store uses a cache with TTL longer than the acceptable revocation window.",
                    "fix": "Always check revocation from the primary database, or use a cache with sub-second invalidation for the revocation denylist."
                }
            ],
            "hidden_magic_checks": [
                "Learner must explain why access tokens should have shorter expiry than refresh tokens, and what happens when both have the same expiry.",
                "Learner must trace the token flow for a stolen-refresh-token scenario: attacker gets token, user rotates, attacker tries to reuse old token."
            ],
            "how_go_uses_it": "Go's standard library provides crypto/rand for secure token generation and database/sql for transactional token storage. The httptest package allows testing the full refresh flow without a real HTTP server. The sync.Mutex or database transactions provide the concurrency safety needed for rotation.",
            "mental_model": "Access tokens are hotel room keys that expire at checkout. Refresh tokens are the front desk: you show your ID to get a new room key. If someone steals your room key, it stops working at checkout. If someone steals your ID (refresh token), rotating the lock renders theirs useless.",
            "performance_implications": [
                "Token rotation adds latency to the refresh endpoint (DB write + read).",
                "Revocation denylist grows unbounded with frequent rotation; implement TTL-based cleanup.",
                "JWTs with short expiry (5-15 minutes) shift verification load to the client, reducing auth server load."
            ],
            "problem_solved": "JWTs cannot be revoked individually — once issued, they remain valid until expiry. Short access token lifetimes (5-15 min) require frequent re-authentication. Refresh tokens solve this by providing a long-lived, revocable credential that issues new short-lived access tokens on demand. The rotation pattern (issued on each refresh use) ensures that a stolen refresh token becomes useless as soon as the legitimate user performs their next refresh.",
            "proof_of_understanding": "The learner must produce a runnable token service that issues access + refresh token pairs, supports rotation with race-condition tests, provides a revocation endpoint, and rejects reused or expired tokens.",
            "real_world_usage": "Every production authentication system uses refresh tokens. GitHub issues access tokens (1hr expiry) and refresh tokens (6mo rotation). Auth0, Keycloak, and Firebase Auth all implement the OAuth2 refresh token grant type. Go microservices behind an API gateway use refresh tokens to maintain session state without forcing regular re-login.",
            "step_by_step_execution": [
                "Generate a cryptographically random refresh token (32 bytes from crypto/rand, hex-encoded).",
                "Store the token hash and expiry in the database alongside the user ID and a rotation counter.",
                "On refresh: look up the token hash, verify expiry and revocation status, and issue a new access + refresh pair.",
                "Increment the rotation counter and mark the old token as revoked in the same database transaction.",
                "On revocation: mark the token as revoked in the database. Check revocation status on every access token request.",
                "Add concurrent rotation protection: use SELECT ... FOR UPDATE or an optimistic version check."
            ],
            "under_the_hood": "The OAuth 2.0 refresh token grant type (RFC 6749 section 6) defines the protocol. The rotation pattern is specified in RFC 6819 (OAuth 2.0 Threat Model) and OAuth 2.0 Security BCP (draft). In Go, database/sql's Tx provides the atomicity needed for safe rotation. The crypto/rand package provides cryptographically secure token generation — math/rand must never be used for tokens.",
            "why_it_exists": "Without refresh tokens, applications face an impossible choice: short JWTs (frequent re-login) or long JWTs (irrevocable exposure). Refresh tokens with rotation provide revocability and short-lived access tokens without requiring the user to re-enter credentials every few minutes. This is not optional — it's the security baseline for any production authentication system."
        }
    }

    items.append(new_item)
    
    # Fix chain: core-10-24's next_item_ids currently points to core-11-01
    # Insert core-10-25 between them
    if last_m10:
        last_m10_next = last_m10.get('next_item_ids', [])
        if last_m10_next:
            new_item['next_item_ids'] = last_m10_next  # core-10-25 links to what 10-24 previously linked to
            last_m10['next_item_ids'] = ['core-10-25']  # 10-24 now links to 10-25
    
    # Update module terminal
    for m in core['modules']:
        if m['id'] == 'module-10':
            m['terminal_item_ids'] = ['core-10-25']
            print(f"  Updated module-10 terminal to core-10-25")
    
    print(f"  Created core-10-25 Refresh tokens")

# ════════════════════════════════════════════════════════════════
# 2. Promote core-11-28 from preview to full lesson
# ════════════════════════════════════════════════════════════════
print("\n=== Promoting core-11-28 ===")

c1128 = item_map.get('core-11-28')
if c1128:
    old_title = c1128['title']
    c1128['title'] = 'Idempotency keys'
    c1128['slug'] = 'idempotency-keys'
    c1128['status'] = 'planned'
    c1128['difficulty'] = 'intermediate'
    
    # Remove preview_only crossrefs (it's now a full lesson)
    crossrefs = c1128.get('crossrefs', {})
    if 'preview_only' in crossrefs:
        # Move preview_only entries to builds_on
        po_entries = crossrefs.get('preview_only', [])
        if po_entries:
            builds_on = crossrefs.get('builds_on', [])
            builds_on.extend(po_entries)
        crossrefs['preview_only'] = []
    
    # Fix learning_objective
    c1128['learning_objective'] = "Design and implement idempotency keys for HTTP APIs in Go."
    
    # Fix ZM mental_model and problem_solved if they're still template
    zm = c1128.get('zero_magic', {})
    if zm:
        mm = zm.get('mental_model', '')
        if 'preview' in mm.lower() or 'Think of' in mm:
            zm['mental_model'] = "An idempotency key is a unique identifier the client sends with a request. The server checks if it has already processed that key — if yes, it returns the stored result instead of executing the operation again. This turns 'at least once' delivery into 'exactly once' processing."
        ps = zm.get('problem_solved', '')
        if 'preview' in ps.lower():
            zm['problem_solved'] = "Network failures, retries, and duplicate delivery cause the same request to arrive multiple times. Without idempotency, duplicate payment charges, duplicate order creations, and duplicate email sends are guaranteed over any non-trivial network. Idempotency keys solve this by making repeated requests safe to process."
    
    print(f"  Promoted: '{old_title}' -> '{c1128['title']}'")

# Write
with open('curriculum/path.core.json', 'w', encoding='utf-8') as f:
    json.dump(core, f, indent=2, ensure_ascii=False)
print("\npath.core.json written")
