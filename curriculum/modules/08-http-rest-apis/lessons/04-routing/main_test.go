package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListTasks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	listTasks(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var body map[string]Task
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(body) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(body))
	}
}

func TestGetTask_Found(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()
	getTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var task Task
	json.NewDecoder(rec.Body).Decode(&task)
	if task.ID != "1" {
		t.Errorf("expected task ID 1, got %s", task.ID)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	req.SetPathValue("id", "999")
	rec := httptest.NewRecorder()
	getTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestCreateTask(t *testing.T) {
	body := `{"id":"3","text":"New task","done":false}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	createTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	var task Task
	json.NewDecoder(rec.Body).Decode(&task)
	if task.ID != "3" {
		t.Errorf("expected task ID 3, got %s", task.ID)
	}
	if _, ok := tasks["3"]; !ok {
		t.Error("expected task 3 to be stored")
	}
	delete(tasks, "3")
}

func TestCreateTask_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	createTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
