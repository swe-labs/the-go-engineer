package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type sessionKeyType string

const sessionKey sessionKeyType = "session"

type Session struct {
	ID        string
	UserID    string
	Role      string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{sessions: make(map[string]*Session)}
}

func (ss *SessionStore) Create(userID, role string, ttl time.Duration) (*Session, error) {
	id, err := generateSessionID()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	s := &Session{
		ID: id, UserID: userID, Role: role,
		CreatedAt: now, ExpiresAt: now.Add(ttl),
	}
	ss.mu.Lock()
	ss.sessions[id] = s
	ss.mu.Unlock()
	return s, nil
}

func (ss *SessionStore) Get(id string) (*Session, error) {
	ss.mu.RLock()
	s, ok := ss.sessions[id]
	ss.mu.RUnlock()
	if !ok {
		return nil, errors.New("session not found")
	}
	if time.Now().After(s.ExpiresAt) {
		ss.Delete(id)
		return nil, errors.New("session expired")
	}
	return s, nil
}

func (ss *SessionStore) Delete(id string) {
	ss.mu.Lock()
	delete(ss.sessions, id)
	ss.mu.Unlock()
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func sessionCookie(sessionID string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name: "session_id", Value: sessionID,
		HttpOnly: true, Secure: false,
		SameSite: http.SameSiteLaxMode, Path: "/",
		MaxAge: maxAge,
	}
}

var store = NewSessionStore()

func CreateSession(store *SessionStore, userID, role string, ttl time.Duration) (*Session, error) {
	return store.Create(userID, role, ttl)
}

func GetSession(store *SessionStore, sessionID string) (*Session, error) {
	return store.Get(sessionID)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if creds.Username != "alice" || creds.Password != "secret123" {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		s, err := store.Create(creds.Username, "user", 7*24*time.Hour)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, sessionCookie(s.ID, 86400*7))
		json.NewEncoder(w).Encode(map[string]string{"status": "logged_in"})
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			store.Delete(cookie.Value)
		}
		http.SetCookie(w, sessionCookie("", -1))
		json.NewEncoder(w).Encode(map[string]string{"status": "logged_out"})
	})
	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(sessionKey).(*Session)
		if !ok || session == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": session.UserID, "role": session.Role,
			"expires_at": session.ExpiresAt,
		})
	})

	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			session, err := store.Get(cookie.Value)
			if err == nil {
				ctx := context.WithValue(r.Context(), sessionKey, session)
				r = r.WithContext(ctx)
			}
		}
		mux.ServeHTTP(w, r)
	})

	fmt.Println("Session server on :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
