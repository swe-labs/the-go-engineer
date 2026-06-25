package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUsersHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()

	usersHandler(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var got []User
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("json.Decode: %v", err)
	}
	resp.Body.Close()

	if len(got) != 2 {
		t.Fatalf("expected 2 users, got %d", len(got))
	}
	if got[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", got[0].Name)
	}
}

func TestUsersHandlerTableDriven(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "GET users returns 200",
			method:     http.MethodGet,
			path:       "/api/users",
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":1,"name":"Alice","role":"admin"},{"id":2,"name":"Bob","role":"user"}]`,
		},
		{
			name:       "POST returns 405 Method Not Allowed",
			method:     http.MethodPost,
			path:       "/api/users",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			handler := http.HandlerFunc(usersHandler)
			handler.ServeHTTP(rec, req)

			resp := rec.Result()
			resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
			if tt.wantBody != "" {
				got := strings.TrimSpace(rec.Body.String())
				if got != tt.wantBody {
					t.Errorf("body = %s, want %s", got, tt.wantBody)
				}
			}
		})
	}
}

func TestGreetingHandler(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{"default greeting", "", "Hello, world!"},
		{"named greeting", "name=Go", "Hello, Go!"},
		{"empty name param", "name=", "Hello, world!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/greet?"+tt.query, nil)
			rec := httptest.NewRecorder()

			greetingHandler(rec, req)

			if got := strings.TrimSpace(rec.Body.String()); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		wantStatus int
		wantBody   string
	}{
		{"admin role allowed", "admin", http.StatusOK, "admin area\n"},
		{"user role forbidden", "user", http.StatusForbidden, "forbidden\n"},
		{"no role forbidden", "", http.StatusForbidden, "forbidden\n"},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "admin area")
	})
	handler := adminMiddleware(next)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
			req.Header.Set("X-Role", tt.role)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			resp := rec.Result()
			resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
			if rec.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestWithServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", usersHandler)
	mux.HandleFunc("/greet", greetingHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/users")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var got []User
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("json.Decode: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 users, got %d", len(got))
	}
}

func TestGoldenFile(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	usersHandler(rec, req)

	goldenPath := filepath.Join("testdata", "users.golden")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}

	got := strings.TrimSpace(rec.Body.String())
	wantStr := strings.TrimSpace(string(want))
	if got != wantStr {
		t.Errorf("response does not match golden file\n got:  %s\n want: %s", got, wantStr)
	}
}

func TestLessonCompiles(t *testing.T) {
}
