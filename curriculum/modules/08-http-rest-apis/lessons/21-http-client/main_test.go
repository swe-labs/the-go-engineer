package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientGetUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]User{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	users, err := client.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("GetUsers: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", users[0].Name)
	}
}

func TestClientCreateUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(User{ID: 3, Name: body["name"]})
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	user, err := client.CreateUser(context.Background(), "Charlie")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	if user.Name != "Charlie" {
		t.Errorf("expected Charlie, got %s", user.Name)
	}
	if user.ID != 3 {
		t.Errorf("expected ID 3, got %d", user.ID)
	}
}

func TestClientGetUsersError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	_, err := client.GetUsers(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention status 500, got: %v", err)
	}
}

func TestClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, 50*time.Millisecond)
	_, err := client.GetUsers(context.Background())
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestClientGetUsersRaw(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`raw response`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 5*time.Second)
	body, err := client.GetUsersRaw(context.Background())
	if err != nil {
		t.Fatalf("GetUsersRaw: %v", err)
	}
	if body != "raw response" {
		t.Errorf("body = %q, want %q", body, "raw response")
	}
}

func TestClientTransportConfig(t *testing.T) {
	client := NewClient("http://localhost:9999", 5*time.Second)

	transport, ok := client.httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatal("transport is not *http.Transport")
	}

	if transport.MaxIdleConns != 100 {
		t.Errorf("MaxIdleConns = %d, want 100", transport.MaxIdleConns)
	}
	if transport.MaxIdleConnsPerHost != 10 {
		t.Errorf("MaxIdleConnsPerHost = %d, want 10", transport.MaxIdleConnsPerHost)
	}
	if transport.IdleConnTimeout != 90*time.Second {
		t.Errorf("IdleConnTimeout = %v, want 90s", transport.IdleConnTimeout)
	}
	if transport.TLSHandshakeTimeout != 5*time.Second {
		t.Errorf("TLSHandshakeTimeout = %v, want 5s", transport.TLSHandshakeTimeout)
	}
}

func TestClientContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client := NewClient(server.URL, 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetUsers(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestLessonCompiles(t *testing.T) {
}
