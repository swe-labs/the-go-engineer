package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestTaskAPI (Function): runs the new test task api step and keeps its inputs, outputs, or errors visible.
func newTestTaskAPI() http.Handler {
	store := NewTaskStore()
	api := &TaskAPI{store: store}
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)
	return mux
}

func TestRESTAPIFlow(t *testing.T) {
	handler := newTestTaskAPI()

	listResp := serveRequest(handler, http.MethodGet, "/tasks", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listResp.Code, http.StatusOK)
	}

	var initial []Task
	if err := json.NewDecoder(listResp.Body).Decode(&initial); err != nil {
		t.Fatalf("decode initial list: %v", err)
	}
	if len(initial) != 0 {
		t.Fatalf("initial task count = %d, want 0", len(initial))
	}

	createResp := serveRequest(handler, http.MethodPost, "/tasks", `{"title":"Complete Go Section"}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", createResp.Code, http.StatusCreated)
	}

	var created Task
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created task: %v", err)
	}
	if created.ID != 1 || created.Title != "Complete Go Section" || created.Completed {
		t.Fatalf("created task = %+v, want id 1, title Complete Go Section, completed false", created)
	}

	getResp := serveRequest(handler, http.MethodGet, "/tasks/1", "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", getResp.Code, http.StatusOK)
	}

	var fetched Task
	if err := json.NewDecoder(getResp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode fetched task: %v", err)
	}
	if fetched.ID != created.ID || fetched.Title != created.Title {
		t.Fatalf("fetched task = %+v, want %+v", fetched, created)
	}

	deleteResp := serveRequest(handler, http.MethodDelete, "/tasks/1", "")
	if deleteResp.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleteResp.Code, http.StatusNoContent)
	}

	missingResp := serveRequest(handler, http.MethodGet, "/tasks/1", "")
	if missingResp.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d, want %d", missingResp.Code, http.StatusNotFound)
	}
}

func TestRESTAPIValidation(t *testing.T) {
	handler := newTestTaskAPI()

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{name: "invalid json", method: http.MethodPost, path: "/tasks", body: "{", wantStatus: http.StatusBadRequest},
		{name: "missing title", method: http.MethodPost, path: "/tasks", body: `{}`, wantStatus: http.StatusUnprocessableEntity},
		{name: "invalid get id", method: http.MethodGet, path: "/tasks/not-a-number", wantStatus: http.StatusBadRequest},
		{name: "invalid delete id", method: http.MethodDelete, path: "/tasks/not-a-number", wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := serveRequest(handler, tt.method, tt.path, tt.body)
			if resp.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", resp.Code, tt.wantStatus, resp.Body.String())
			}
		})
	}
}

// serveRequest (Function): runs the serve request step and keeps its inputs, outputs, or errors visible.
func serveRequest(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	return resp
}
