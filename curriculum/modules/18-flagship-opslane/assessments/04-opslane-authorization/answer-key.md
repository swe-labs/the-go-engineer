# Opslane Milestone 04 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) 401 means the caller has not supplied valid authentication credentials; 403 means the caller is authenticated, but lacks sufficient permissions to perform the requested operation.**  
*Explanation:* Under RFC 9110, 401 Unauthorized indicates unauthenticated requests (missing or invalid credentials), prompting authentication challenges. 403 Forbidden indicates the server understands who the caller is, but refuses authorization.

**2. b) To prevent context key collisions across different packages or third-party middleware that might also store values using string keys.**  
*Explanation:* In Go, interface equality for context keys checks both type and value. Defining an unexported package-level type (`type contextKey string`) ensures that even if another package stores a value under `"user_context"`, its type will differ, preventing accidental key overwrites or data corruption.

**3. b) To achieve loose coupling: route handlers require specific capabilities, allowing role mappings to change or custom roles to be created without modifying endpoint code.**  
*Explanation:* Binding endpoints to roles (e.g. `if role == "admin"`) tightly couples business routes to identity hierarchy. Gating by granular permission (`PermOrderCancel`) allows system administrators to adjust which roles have cancel permissions dynamically without altering code.

**4. b) HTTP 401 Unauthorized.**  
*Explanation:* When no token is present, the request is unauthenticated, warranting HTTP 401 with a descriptive JSON error.

**5. b) RoleOwner.**  
*Explanation:* In `RolePermissions`, only `RoleOwner` holds `PermTenantManage`. `RoleAdmin` can manage orders and view audits, but cannot modify tenant-level configuration.

**6. a) `order:read` only.**  
*Explanation:* `RoleViewer` is restricted strictly to read-only capabilities (`PermOrderRead`) to enforce read-only audit and monitoring views.

---

## Short Answer (2 points each)

**7. Principle of Least Privilege in Opslane:**
- **Principle:** Every module, user, and process must be granted only the minimum set of permissions necessary to perform its intended task, limiting the potential blast radius of compromised accounts or insider threats.
- **Opslane Implementation:** Opslane defines 4 distinct tiers: `viewer` has read-only access to orders; `member` can create orders but cannot cancel them; `admin` can cancel orders and view audit trails; only `owner` can manage tenant settings.
*Grading: 1 pt for definition of Least Privilege; 1 pt for demonstrating how the 4-tier matrix restricts capabilities.*

**8. Safe context retrieval in `UserFromContext`:**
In Go, extracting an interface value via `ctx.Value(key)` returns `any`. Direct type assertion `u := ctx.Value(key).(*UserContext)` will cause an immediate runtime panic if the key is missing or contains an unexpected type. Using the comma-ok idiom `u, ok := ctx.Value(key).(*UserContext)` safely verifies the type without panicking, and the function returns `nil, false` if the value is missing or nil.
*Grading: 1 pt for explaining panic prevention via comma-ok assertion; 1 pt for safe nil checking.*

**9. RBAC vs. Attribute-Based Access Control (ABAC):**
- **RBAC:** Grants access based solely on static role assignment (e.g. "Admins can cancel orders").
- **ABAC:** Grants access based on dynamic attributes of the user, resource, and environment (e.g. "Users can cancel orders only if `order.user_id == user.id` AND `order.total < $500` AND time is within business hours").
- **Opslane Scenario:** Opslane would require ABAC for self-service order cancellation (allowing regular members to cancel *their own* orders within 30 minutes of placement, while forbidding cancellation of other users' orders).
*Grading: 1 pt for distinguishing static roles from dynamic attributes; 1 pt for providing a valid business scenario.*

**10. Confused Deputy attacks in multi-tenancy:**
- **Attack Scenario:** A malicious user from Tenant A presents a valid authentication token for Tenant A, but crafts an API request specifying an entity ID belonging to Tenant B (e.g. `POST /api/v1/orders/cancel?order_id=tenant_b_order`). If the authorization middleware only verifies that the caller has `order:cancel` permission globally, the server acts as a "confused deputy," executing Tenant A's command against Tenant B's data.
- **Prevention:** Context-driven authorization must bind the caller's verified `TenantID` from their token claims, and repository queries must strictly filter by both `tenant_id` and entity ID (`WHERE tenant_id = ? AND id = ?`).
*Grading: 1 pt for explaining the cross-tenant confused deputy scenario; 1 pt for explaining mitigation via context-bound tenant scoping.*
