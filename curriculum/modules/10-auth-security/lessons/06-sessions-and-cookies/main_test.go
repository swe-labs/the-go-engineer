package main

import (
	"testing"
	"time"
)

func TestCreateSession_Valid(t *testing.T) {
	s := NewSessionStore()
	session, err := CreateSession(s, "user1", "admin", 1*time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if session.ID == "" {
		t.Error("expected non-empty session ID")
	}
	if len(session.ID) != 64 {
		t.Errorf("expected 64-char hex ID, got %d chars", len(session.ID))
	}
	if session.UserID != "user1" {
		t.Errorf("expected user1, got %s", session.UserID)
	}
}

func TestGetSession_Valid(t *testing.T) {
	s := NewSessionStore()
	created, _ := CreateSession(s, "user1", "user", 1*time.Hour)
	fetched, err := GetSession(s, created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, fetched.ID)
	}
}

func TestGetSession_NotFound(t *testing.T) {
	s := NewSessionStore()
	_, err := GetSession(s, "nonexistent-id")
	if err == nil {
		t.Fatal("expected error for nonexistent session")
	}
}

func TestGetSession_Expired(t *testing.T) {
	store := NewSessionStore()
	session, _ := store.Create("user1", "user", -1*time.Hour) // expired
	_, err := GetSession(store, session.ID)
	if err == nil {
		t.Fatal("expected error for expired session")
	}
}

func TestGenerateSessionID(t *testing.T) {
	id1, _ := generateSessionID()
	id2, _ := generateSessionID()
	if id1 == id2 {
		t.Error("expected unique session IDs")
	}
	if len(id1) != 64 {
		t.Errorf("expected 64 chars, got %d", len(id1))
	}
}
