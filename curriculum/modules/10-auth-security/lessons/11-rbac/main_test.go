package main

import (
	"testing"
)

func TestBuildRBAC(t *testing.T) {
	roleDefs := map[string][]string{
		"viewer": {"documents:read"},
		"editor": {"documents:write", "documents:delete"},
		"admin":  {"users:admin"},
	}
	hierarchy := map[string][]string{
		"editor": {"viewer"},
		"admin":  {"editor"},
	}
	engine := BuildRBAC(roleDefs, hierarchy)
	if engine == nil {
		t.Fatal("expected non-nil engine")
	}
}

func TestCheckAccess_AdminHasAll(t *testing.T) {
	engine := BuildRBAC(
		map[string][]string{"viewer": {"documents:read"}, "admin": {"users:admin"}},
		map[string][]string{"admin": {"viewer"}},
	)
	engine.AssignRole("carol", "admin")
	if !CheckAccess(engine, "carol", "documents", "read") {
		t.Error("expected admin to have read access")
	}
}

func TestCheckAccess_ViewerReadOnly(t *testing.T) {
	engine := BuildRBAC(
		map[string][]string{"viewer": {"documents:read"}},
		nil,
	)
	engine.AssignRole("alice", "viewer")
	if !CheckAccess(engine, "alice", "documents", "read") {
		t.Error("expected viewer to have read access")
	}
	if CheckAccess(engine, "alice", "documents", "write") {
		t.Error("expected viewer to NOT have write access")
	}
}

func TestCheckAccess_UnknownUser(t *testing.T) {
	engine := BuildRBAC(
		map[string][]string{"viewer": {"documents:read"}},
		nil,
	)
	if CheckAccess(engine, "unknown", "documents", "read") {
		t.Error("expected unknown user to have NO access")
	}
}

func TestCheckAccess_InheritedHierarchy(t *testing.T) {
	engine := BuildRBAC(
		map[string][]string{"viewer": {"documents:read"}, "editor": {"documents:write"}},
		map[string][]string{"editor": {"viewer"}},
	)
	engine.AssignRole("bob", "editor")
	if !CheckAccess(engine, "bob", "documents", "read") {
		t.Error("expected editor to inherit read from viewer")
	}
	if !CheckAccess(engine, "bob", "documents", "write") {
		t.Error("expected editor to have write")
	}
}

func TestHasPermission_Wildcard(t *testing.T) {
	e := NewRBACEngine()
	e.AddRole("admin", []Permission{"documents:*"}, nil)
	e.AssignRole("carol", "admin")
	if !e.HasPermission("carol", "documents:delete") {
		t.Error("expected wildcard match")
	}
	if !e.HasPermission("carol", "documents:create") {
		t.Error("expected wildcard match for create")
	}
}
