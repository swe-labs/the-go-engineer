package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Permission string

type RBACEngine struct {
	mu        sync.RWMutex
	roles     map[string]*Role
	userRoles map[string][]string
}

type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
	Inherits    []string     `json:"inherits,omitempty"`
}

func NewRBACEngine() *RBACEngine {
	return &RBACEngine{
		roles:     make(map[string]*Role),
		userRoles: make(map[string][]string),
	}
}

func (e *RBACEngine) AddRole(name string, permissions []Permission, inherits []string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.roles[name] = &Role{Name: name, Permissions: permissions, Inherits: inherits}
}

func (e *RBACEngine) AssignRole(userID, roleName string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.roles[roleName]; !ok {
		return fmt.Errorf("role %s does not exist", roleName)
	}
	e.userRoles[userID] = append(e.userRoles[userID], roleName)
	return nil
}

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

func BuildRBAC(roleDefs map[string][]string, hierarchy map[string][]string) *RBACEngine {
	e := NewRBACEngine()
	for name, perms := range roleDefs {
		var ps []Permission
		for _, p := range perms {
			ps = append(ps, Permission(p))
		}
		e.AddRole(name, ps, hierarchy[name])
	}
	return e
}

func CheckAccess(engine *RBACEngine, userID, resource, action string) bool {
	return engine.HasPermission(userID, Permission(resource+":"+action))
}

const (
	PermDocRead   Permission = "documents:read"
	PermDocWrite  Permission = "documents:write"
	PermDocDelete Permission = "documents:delete"
	PermUserAdmin Permission = "users:admin"
)

func main() {
	rbac := NewRBACEngine()
	rbac.AddRole("viewer", []Permission{PermDocRead}, nil)
	rbac.AddRole("editor", []Permission{PermDocWrite, PermDocDelete}, []string{"viewer"})
	rbac.AddRole("admin", []Permission{PermUserAdmin, "documents:*"}, []string{"editor"})

	rbac.AssignRole("alice", "viewer")
	rbac.AssignRole("bob", "editor")
	rbac.AssignRole("carol", "admin")

	for _, u := range []string{"alice", "bob", "carol"} {
		fmt.Printf("User %s permissions: %v\n", u, rbac.GetUserPermissions(u))
	}

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

	// HTTP
	mux := http.NewServeMux()
	mux.Handle("/api/docs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, `{"error":"unauthenticated"}`, http.StatusUnauthorized)
			return
		}
		if !rbac.HasPermission(userID, PermDocRead) {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"docs": "list"})
	}))

	fmt.Println("\nRBAC server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
