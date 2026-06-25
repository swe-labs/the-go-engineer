# Sessions and cookies

## Learning objective

Implement session-based authentication in Go using secure cookies (HttpOnly, Secure, SameSite), manage session state in a server-side store, and handle session expiry, rotation, and invalidation.

## Why this matters

Sessions are the original web authentication mechanism and remain the most practical choice for server-rendered applications and SPAs with a backend. Unlike JWTs, sessions support instant revocation: delete the session from the store and the user is logged out. Sessions also keep sensitive state on the server, avoiding the token size bloat and claim exposure of client-side tokens. Every major platform — GitHub, Stripe Dashboard, AWS Console — uses sessions for their web interfaces. Understanding session security (cookie flags, CSRF tokens, session fixation prevention) is essential for any Go web developer.

## Mental model

A session is a lockbox on the server. When a user logs in, the server creates a lockbox (session), puts their identity inside, and gives the user a key (session ID). The key is stored in a cookie. On every request, the user presents the key, and the server opens the lockbox to read the identity. If the user logs out, the server destroys the lockbox — the key becomes useless even if stolen. The lockbox never leaves the server, so the contents (user ID, roles, preferences) are safe even if the cookie is intercepted.

## Core idea

A session system has three components:

1. **Session store** — where session data lives on the server (memory, Redis, database).
2. **Session ID** — a cryptographically random token that maps to a session in the store.
3. **Cookie** — the mechanism for transporting the session ID between client and server.

Cookie security flags that protect the session ID:

| Flag | Effect | Why it matters |
|---|---|---|
| `HttpOnly` | JavaScript cannot read the cookie via `document.cookie` | Prevents XSS-based session theft |
| `Secure` | Cookie is only sent over HTTPS | Prevents MITM interception on HTTP |
| `SameSite=Lax` | Cookie is not sent on cross-site requests | Prevents CSRF attacks |
| `SameSite=Strict` | Cookie is not sent on any cross-site request | Stronger CSRF protection, may break OAuth flows |
| `Path=/` | Cookie is sent to all paths on the domain | Must scope to the minimum necessary path |
| `MaxAge` or `Expires` | Cookie lifetime | Browsers delete expired cookies |

## Under the hood

In Go, cookies are set via the `net/http` package:

```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_id",
    Value:    sessionID,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteLaxMode,
    Path:     "/",
    MaxAge:   86400 * 7, // 7 days
})
```

The session ID must be generated using a cryptographically secure random generator. Go's `crypto/rand` is the correct choice:

```go
func generateSessionID() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}
```

Session stores can be in-memory (`sync.Map`), Redis, or a database. In-memory stores are simple but lost on restart. Redis-backed sessions survive restarts and scale horizontally.

## How Go uses it

The standard library does not include a session manager, but several patterns are well-established:

- **Gorilla sessions**: `github.com/gorilla/sessions` provides cookie-backed and Redis-backed sessions.
- **SCS**: `github.com/alexedwards/scs` is a modern, lightweight session manager with Redis, PostgreSQL, and MySQL stores.
- **Custom**: Many production Go services implement their own session middleware using `net/http` and `crypto/rand`.

The typical middleware pattern:

```go
func sessionMiddleware(store SessionStore) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            cookie, err := r.Cookie("session_id")
            if err != nil {
                // No session — continue with empty context.
                next.ServeHTTP(w, r)
                return
            }
            session, err := store.Get(cookie.Value)
            if err != nil {
                http.SetCookie(w, clearSessionCookie())
                next.ServeHTTP(w, r)
                return
            }
            ctx := context.WithValue(r.Context(), sessionKey, session)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## Go example

```go
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

// Session holds data for an authenticated user.
type Session struct {
	ID        string
	UserID    string
	Role      string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SessionStore is a simple in-memory session store.
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
		ID:        id,
		UserID:    userID,
		Role:      role,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
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

// sessionCookie creates a secure session cookie.
func sessionCookie(sessionID string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   true, // set to false for local dev without TLS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   maxAge,
	}
}

var store = NewSessionStore()

// loginHandler authenticates credentials and creates a session.
func loginHandler(w http.ResponseWriter, r *http.Request) {
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
	// In production: verify against database with bcrypt.
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged_in"})
}

// logoutHandler destroys the session.
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		store.Delete(cookie.Value)
	}
	http.SetCookie(w, sessionCookie("", -1)) // Clear cookie
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged_out"})
}

// profileHandler shows the user's session data.
func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(sessionKey).(*Session)
	if !ok || session == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":    session.UserID,
		"role":       session.Role,
		"created_at": session.CreatedAt,
		"expires_at": session.ExpiresAt,
	})
}

// sessionMiddleware loads the session on every request.
func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// No cookie: continue without session.
			next.ServeHTTP(w, r)
			return
		}
		session, err := store.Get(cookie.Value)
		if err != nil {
			// Invalid or expired session: clear cookie.
			http.SetCookie(w, sessionCookie("", -1))
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), sessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/profile", profileHandler)

	var h http.Handler = mux
	h = sessionMiddleware(h)

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
```

## Step-by-step execution

For a login request `POST /login` with `{"username":"alice","password":"secret123"}`:

1. `loginHandler` decodes JSON. Validates credentials against database (simulated).
2. On success, `store.Create("alice", "user", 7*24*time.Hour)` generates a 64-character hex session ID.
3. Session is stored in `store.sessions` map.
4. Response includes `Set-Cookie: session_id=<id>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=604800`.
5. Browser stores the cookie and sends it on subsequent requests.

For `GET /profile` with the session cookie:
1. `sessionMiddleware` reads `Cookie: session_id=<id>`.
2. Calls `store.Get(<id>)`. Checks expiry: not expired.
3. Stores `*Session` in request context.
4. `profileHandler` reads session from context, returns user data.

For logout: `store.Delete(id)` removes the session. The cookie is cleared with `MaxAge: -1`.

## Common mistakes

- Storing sessions only in memory with no persistence — all sessions are lost on restart. Use Redis or a database.
- Using predictable session IDs — session tokens must be cryptographically random. `math/rand` is not acceptable.
- Not setting `HttpOnly` and `Secure` flags on session cookies — exposes the session ID to XSS and man-in-the-middle attacks.
- Not rotating session IDs after login — session fixation attack: an attacker sets a known session ID before the user logs in.
- Using `SameSite=None` without `Secure` — modern browsers reject SameSite=None without Secure.
- Not validating session expiry on every request — expired sessions should be deleted immediately.

## Debugging walkthrough

Consider this session code:

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    // ... authenticate user ...
    id := fmt.Sprintf("%d", time.Now().UnixNano())
    http.SetCookie(w, &http.Cookie{
        Name:  "session",
        Value: id,
    })
}
```

**Symptom**: Users are frequently logged in as other users. Session IDs are predictable.

**Investigation**: The session ID is the current Unix nanosecond timestamp — completely predictable. An attacker can enumerate session IDs. If a session is valid for 7 days, the attacker can try one ID per millisecond and hijack sessions.

**Root cause**: Using `time.Now().UnixNano()` instead of `crypto/rand`. The session ID has zero entropy — it is simply a timestamp.

**Fix**:

```go
func generateSessionID() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}
```

## Production notes

- Use Redis-backed sessions for horizontal scaling. All server instances read from the same Redis cluster.
- Set session TTL to the minimum practical value. 24 hours for most web apps, 15 minutes for banking applications.
- Implement session rotation: after login, generate a new session ID and delete the old one. This prevents session fixation.
- Monitor session store size. A Redis instance with millions of expired sessions wastes memory. Set `EXPIRE` on session keys.
- For APIs that cannot use cookies (mobile apps), accept the session ID via an `Authorization: Bearer <session_id>` header.

## Performance implications

- In-memory session lookup: O(1), sub-microsecond.
- Redis session lookup: ~1-5ms per request (network round trip). Use a local Redis or connection pooling.
- Database-backed sessions: ~5-50ms per request. Use an index on session_id. Consider a separate session table.
- Cookie parsing: negligible overhead. Go's `r.Cookie()` is O(n) on the number of cookies, but typically n < 10.
- Session ID generation via `crypto/rand`: ~1 microsecond. Negligible.

## Practice task

Write Go functions `CreateSession(store *SessionStore, userID, role string, ttl time.Duration) (*Session, error)` and `GetSession(store *SessionStore, sessionID string) (*Session, error)` that:
- `CreateSession` generates a cryptographically random 32-byte session ID (hex-encoded).
- `CreateSession` sets `CreatedAt` to now and `ExpiresAt` to now + ttl.
- `GetSession` returns an error if the session is not found or expired.
- `GetSession` automatically deletes expired sessions.
- Then write a `main()` that creates 3 sessions, retrieves 2, and attempts to retrieve an expired one.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/06-sessions-and-cookies
go test ./curriculum/modules/10-auth-security/lessons/06-sessions-and-cookies
```

The test file `main_test.go` contains table-driven tests that verify:
- `CreateSession` returns a session with a non-empty hex ID.
- `GetSession` returns the session for a valid ID.
- `GetSession` returns an error for a nonexistent ID.
- `GetSession` returns an error for an expired session.

## Review questions

1. What is the difference between a session and a session ID? Which one is stored on the client and which on the server?
2. What attack does the `HttpOnly` cookie flag prevent? What attack does `SameSite=Lax` prevent?
3. Why must session IDs be generated with `crypto/rand` rather than `math/rand`?
4. Describe a session fixation attack and how session rotation after login prevents it.
5. If your Go service runs on multiple instances behind a load balancer, what session store configuration is needed?

## NEXT UP

JWT implementation — self-contained tokens with signed claims for stateless authentication across services.
