# RBAC

## Learning objective

Design and implement role-based access control in Go using roles, permissions, and role hierarchies, and enforce authorization decisions in HTTP middleware.

## Why this matters

RBAC is the dominant authorization model for enterprise applications. Kubernetes, AWS IAM, GitHub, and PostgreSQL all use RBAC. It is the natural choice when access decisions depend on a user's function (admin, editor, viewer) rather than fine-grained resource attributes. In Go services, RBAC middleware provides a declarative way to protect routes: `r.Handle("/admin", rbac.RequireRole("admin"), handler)`. Engineers who understand RBAC design patterns build authorization systems that are auditable, maintainable, and extensible.

## Mental model

RBAC is a permission matrix with three dimensions: users, roles, and permissions. A user is assigned one or more roles. Each role grants a set of permissions. A permission is an action on a resource (e.g., "delete document"). The check is simple: does the user have a role that includes the required permission? If yes, allow. If no, deny. Roles can inherit from other roles (role hierarchy), where a senior role includes all permissions of a junior role.

## Core idea

RBAC has three core entities:

- **User**: a person or service account.
- **Role**: a named collection of permissions (e.g., "editor", "admin").
- **Permission**: a specific action on a resource (e.g., "document:delete").

Permissions are typically expressed as `action:resource` pairs or `resource.action` strings. This naming convention is arbitrary but must be consistent across the system.

Role hierarchy (role inheritance):

```
admin (inherits editor + viewer + own permissions)
  -> editor (inherits viewer + own permissions)
    -> viewer (documents:read)
```

When checking permission for an admin, the system checks the admin's permissions plus all inherited permissions from editor and viewer.

Default-deny is the fundamental principle: if no role explicitly grants a permission, access is denied.

## Under the hood

A minimal RBAC implementation in Go needs:

```go
type RBAC struct {
    roles       map[string]Role
    userRoles   map[string][]string // userID -> role names
}

type Role struct {
    Name        string
    Permissions []string
    Parents     []string // role hierarchy
}

func (r *RBAC) HasPermission(userID, requiredPerm string) bool {
    roleNames := r.userRoles[userID]
    for _, name := range roleNames {
        if r.roleHasPermission(name, requiredPerm) {
            return true
        }
    }
    return false
}

func (r *RBAC) roleHasPermission(roleName, perm string) bool {
    role, ok := r.roles[roleName]
    if !ok {
        return false
    }
    for _, p := range role.Permissions {
        if match(p, perm) {
            return true
        }
    }
    // Check parent roles (hierarchy).
    for _, parent := range role.Parents {
        if r.roleHasPermission(parent, perm) {
            return true
        }
    }
    return false
}
```

## How Go uses it

In production Go services, RBAC is implemented as middleware:

```go
func RequirePermission(perm string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := getUserFromContext(r.Context())
            if !rbac.HasPermission(user.ID, perm) {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

This middleware can be composed per route:

```go
r.Handle("/api/documents", authMiddleware(RequirePermission("documents:read")(listHandler)))
r.Handle("/api/documents", authMiddleware(RequirePermission("documents:create")(createHandler)))
```

Third-party packages like Casbin provide policy storage and enforcement engines, but a custom RBAC implementation in Go is simple enough for most applications.

## Go example

```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Permission represents a single action on a resource.
type Permission string

// RBACEngine manages roles and permissions.
type RBACEngine struct {
	mu        sync.RWMutex
	roles     map[string]*Role
	userRoles map[string][]string
}

// Role defines a named set of permissions with optional inheritance.
type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
	Inherits    []string     `json:"inherits,omitempty"`
}

// NewRBACEngine creates an empty RBAC engine.
func NewRBACEngine() *RBACEngine {
	return &RBACEngine{
		roles:     make(map[string]*Role),
		userRoles: make(map[string][]string),
	}
}

// AddRole registers a role with permissions and inheritance.
func (e *RBACEngine) AddRole(name string, permissions []Permission, inherits []string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.roles[name] = &Role{
		Name:        name,
		Permissions: permissions,
		Inherits:    inherits,
	}
}

// AssignRole assigns a role to a user.
func (e *RBACEngine) AssignRole(userID, roleName string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.roles[roleName]; !ok {
		return fmt.Errorf("role %s does not exist", roleName)
	}
	e.userRoles[userID] = append(e.userRoles[userID], roleName)
	return nil
}

// HasPermission checks if a user has a permission through any assigned role.
func (e *RBACEngine) HasPermission(userID string, perm Permission) bool {
	e.mu.RLock()
	roleNames := e.userRoles[userID]
	e.mu.RUnlock()

	for _, name := range roleNames {
		e.mu.RLock()
		role := e.roles[name]
		e.mu.RUnlock()
		if role == nil {
			continue
		}
		if e.roleHasPermission(role, perm) {
			return true
		}
	}
	return false
}

func (e *RBACEngine) roleHasPermission(role *Role, perm Permission) bool {
	for _, p := range role.Permissions {
		if permissionMatches(p, perm) {
			return true
		}
	}
	for _, parentName := range role.Inherits {
		e.mu.RLock()
		parent := e.roles[parentName]
		e.mu.RUnlock()
		if parent != nil && e.roleHasPermission(parent, perm) {
			return true
		}
	}
	return false
}

// permissionMatches supports wildcard matching: "documents:*" matches "documents:read".
func permissionMatches(allow, required Permission) bool {
	if allow == Permission("*") {
		return true
	}
	if strings.HasSuffix(string(allow), ":*") {
		prefix := strings.TrimSuffix(string(allow), ":*")
		return strings.HasPrefix(string(required), prefix+":")
	}
	return allow == required
}

// Middleware returns an HTTP middleware that enforces a permission.
func (e *RBACEngine) Middleware(perm Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract user ID from context (set by auth middleware).
			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				http.Error(w, `{"error":"unauthenticated"}`, http.StatusUnauthorized)
				return
			}
			if !e.HasPermission(userID, perm) {
				http.Error(w, `{"error":"insufficient permissions"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserPermissions returns all permissions for a user (for debugging/audit).
func (e *RBACEngine) GetUserPermissions(userID string) []Permission {
	var result []Permission
	e.mu.RLock()
	roleNames := e.userRoles[userID]
	e.mu.RUnlock()
	seen := make(map[Permission]bool)
	for _, name := range roleNames {
		e.mu.RLock()
		role := e.roles[name]
		e.mu.RUnlock()
		if role == nil {
			continue
		}
		e.collectPermissions(role, seen, make(map[string]bool))
	}
	for p := range seen {
		result = append(result, p)
	}
	return result
}

func (e *RBACEngine) collectPermissions(role *Role, seen map[Permission]bool, visited map[string]bool) {
	if visited[role.Name] {
		return
	}
	visited[role.Name] = true
	for _, p := range role.Permissions {
		seen[p] = true
	}
	for _, parentName := range role.Inherits {
		e.mu.RLock()
		parent := e.roles[parentName]
		e.mu.RUnlock()
		if parent != nil {
			e.collectPermissions(parent, seen, visited)
		}
	}
}

// Permissions as constants.
const (
	PermDocRead   Permission = "documents:read"
	PermDocWrite  Permission = "documents:write"
	PermDocDelete Permission = "documents:delete"
	PermUserAdmin Permission = "users:admin"
)

func main() {
	rbac := NewRBACEngine()

	// Define roles with hierarchy.
	rbac.AddRole("viewer", []Permission{PermDocRead}, nil)
	rbac.AddRole("editor", []Permission{PermDocWrite, PermDocDelete}, []string{"viewer"})
	rbac.AddRole("admin", []Permission{PermUserAdmin, "documents:*"}, []string{"editor"})

	// Assign roles to users.
	rbac.AssignRole("alice", "viewer")
	rbac.AssignRole("bob", "editor")
	rbac.AssignRole("carol", "admin")

	// Test permissions.
	users := []string{"alice", "bob", "carol"}
	for _, u := range users {
		fmt.Printf("User %s permissions: %v\n", u, rbac.GetUserPermissions(u))
	}

	// Permission checks.
	fmt.Println("\n=== Permission Checks ===")
	tests := []struct {
		user string
		perm Permission
	}{
		{"alice", PermDocRead},
		{"alice", PermDocWrite},
		{"bob", PermDocRead},
		{"bob", PermDocDelete},
		{"carol", PermDocDelete},
		{"carol", PermUserAdmin},
		{"unknown", PermDocRead},
	}
	for _, t := range tests {
		result := rbac.HasPermission(t.user, t.perm)
		fmt.Printf("%s can %s: %v\n", t.user, t.perm, result)
	}

	// HTTP server with RBAC middleware.
	mux := http.NewServeMux()
	mux.Handle("/api/docs", rbac.Middleware(PermDocRead)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"docs": "list"})
	})))
	mux.Handle("/api/admin", rbac.Middleware(PermUserAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"panel": "admin"})
	})))

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Ensure errors is used.
var _ = errors.New
```

## Step-by-step execution

For a request by `bob` (editor) to `GET /api/docs`:

1. Middleware reads `X-User-ID: bob` from headers.
2. `rbac.HasPermission("bob", "documents:read")` is called.
3. Looks up roles for bob: `["editor"]`.
4. Gets the `editor` role: permissions `["documents:write", "documents:delete"]`, inherits `["viewer"]`.
5. `roleHasPermission` checks editor's direct permissions: no match for `documents:read`.
6. Recursively checks parent roles: `viewer`. The `viewer` role has `["documents:read"]`. Match found.
7. Returns true. Request proceeds.

For a request by `alice` (viewer) to `DELETE /api/docs/123`:

1. `rbac.HasPermission("alice", "documents:delete")`.
2. Role `viewer` has `["documents:read"]` only. No inheritance.
3. Returns false. 403 Forbidden.

## Common mistakes

- Hardcoding role checks in business logic instead of using a policy layer — changes require code deploys.
- Using only roles without resource-level scoping — an admin role should not grant access to all resources across all tenants.
- Storing role assignments in application code rather than a database — roles should be dynamic.
- Implementing RBAC without default-deny — missing role checks implicitly allow access.
- Making role hierarchy too deep — more than 3 levels of inheritance makes permission audits difficult.

## Debugging walkthrough

Consider this RBAC check:

```go
func requireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := getUser(r)
            if user.Role != role {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

**Symptom**: Users with the "admin" role cannot access routes that require "editor". The role hierarchy is not implemented.

**Investigation**: The middleware checks exact role name (`user.Role != role`). There is no hierarchy or permission lookup. An admin is not an editor, so the check fails.

**Root cause**: The middleware checks role equality, not permission membership. It does not support role hierarchy.

**Fix**: Use a proper RBAC engine with hierarchy:

```go
func requirePermission(perm string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := getUser(r)
            if !rbac.HasPermission(user.ID, perm) {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## Production notes

- Store roles and assignments in a database (PostgreSQL, Redis). Cache frequently accessed roles in memory.
- Implement permission audit logging: log every authZ decision with user, resource, permission, and result.
- Use permission strings consistently: `action:resource` is the most common convention (e.g., `document:delete`, `user:create`).
- RBAC is best for coarse-grained access control. For fine-grained control (resource ownership, attribute conditions), use ABAC in combination with RBAC.
- Test permission changes with table-driven tests: define a matrix of roles x permissions x expected results.

## Performance implications

- RBAC permission check: O(r * p) where r = number of roles per user, p = permissions per role. Typically < 1 microsecond.
- Role hierarchy traversal: O(d) where d = depth of hierarchy. Depth is usually 2-3 levels.
- Database-backed role storage: cache the role-permission mapping in memory (invalidate on role change).
- Wildcard matching: `"documents:*"` matching `"documents:delete"` is O(n) on the number of permission segments. Negligible.

## Practice task

Write Go functions `BuildRBAC(roleDefs map[string][]string, hierarchy map[string][]string) *RBACEngine` and `CheckAccess(engine *RBACEngine, userID, resource, action string) bool` where:
- `roleDefs` maps role name to permission strings (e.g., `"editor": ["write", "delete"]`).
- `hierarchy` maps role name to parent role names (e.g., `"admin": ["editor"]`).
- `CheckAccess` checks if the user has a role with permission `resource:action`.
- Then write a `main()` that defines 3 roles (viewer, editor, admin), assigns them to 2 users, and tests 6 access scenarios.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/11-rbac
go test ./curriculum/modules/10-auth-security/lessons/11-rbac
```

The test file `main_test.go` contains table-driven tests that verify:
- `CheckAccess` returns true for a user with the matching role and permission.
- `CheckAccess` returns true for inherited permissions.
- `CheckAccess` returns false for a user without the required permission.
- `CheckAccess` returns false for an unknown user.
- Wildcard permission matches any action on the resource.

## Review questions

1. What are the three core entities in RBAC and how do they relate to each other?
2. How does role hierarchy differ from assigning multiple roles to a single user?
3. Why is default-deny important in RBAC? What happens if a new permission is added but no role grants it?
4. What is the difference between a permission and a role? Can a user have both?
5. How would you test that an RBAC system correctly denies a revoked user's previously granted permissions?

## NEXT UP

ABAC and policy checks — attribute-based access control for fine-grained, context-aware authorization decisions.
