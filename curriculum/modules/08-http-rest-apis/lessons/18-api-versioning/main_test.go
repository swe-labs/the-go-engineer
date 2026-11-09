package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersV1Handler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	usersV1Handler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var users []UserV1
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
	if users[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", users[0].Name)
	}
}

func TestUsersV2Handler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v2/users", nil)
	rec := httptest.NewRecorder()
	usersV2Handler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var users []UserV2
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
	if users[0].Email != "alice@example.com" {
		t.Errorf("expected alice@example.com, got %s", users[0].Email)
	}
}

func TestURLPathVersioning(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantLen int
	}{
		{"v1 returns UserV1", "/api/v1/users", 2},
		{"v2 returns UserV2", "/api/v2/users", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			usersV1Handler(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected 200, got %d", resp.StatusCode)
			}
		})
	}
}

func TestVersionMiddleware(t *testing.T) {
	handler := versionMiddleware("v2", usersV2Handler)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.Header.Get("X-API-Version") != "v2" {
		t.Errorf("expected X-API-Version: v2, got %s", resp.Header.Get("X-API-Version"))
	}
}

func TestHeaderVersioning(t *testing.T) {
	tests := []struct {
		name       string
		accept     string
		wantStatus int
	}{
		{
			name:       "v2 via accept header",
			accept:     "application/vnd.api.v2+json",
			wantStatus: http.StatusOK,
		},
		{
			name:       "v1 via accept header",
			accept:     "application/vnd.api.v1+json",
			wantStatus: http.StatusOK,
		},
		{
			name:       "default to v1",
			accept:     "application/json",
			wantStatus: http.StatusOK,
		},
	}

	vh := &VersionHandler{
		handlers: map[string]http.HandlerFunc{
			"v1": usersV1Handler,
			"v2": usersV2Handler,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
			req.Header.Set("Accept", tt.accept)
			rec := httptest.NewRecorder()
			vh.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/v1/users", "v1"},
		{"/v2/users", "v2"},
		{"/api/v1/users", "v1"}, // v1 is the default fallback
		{"/users", "v1"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := extractVersion(tt.path)
			if got != tt.want {
				t.Errorf("extractVersion(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestConvertV1ToV2(t *testing.T) {
	u := UserV1{ID: 1, Name: "Alice Smith"}
	v2 := convertV1ToV2(u)

	if v2.FirstName != "Alice" {
		t.Errorf("FirstName = %q, want %q", v2.FirstName, "Alice")
	}
	if v2.LastName != "Smith" {
		t.Errorf("LastName = %q, want %q", v2.LastName, "Smith")
	}
	if v2.Email != "" {
		t.Errorf("Email should be empty, got %q", v2.Email)
	}

	u2 := UserV1{ID: 2, Name: "Bob"}
	v2b := convertV1ToV2(u2)
	if v2b.LastName != "" {
		t.Errorf("LastName should be empty for single name, got %q", v2b.LastName)
	}
}

func TestLessonCompiles(t *testing.T) {
}
