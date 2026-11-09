package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (*sql.DB, *TaskRepository) {
	t.Helper()
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	ConfigurePool(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RunMigrations(ctx, db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	repo := NewTaskRepository(db, "sqlite")
	return db, repo
}

func TestPoolConfiguration(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	stats := db.Stats()
	if stats.MaxOpenConnections != 25 {
		t.Fatalf("expected MaxOpenConnections 25, got %d", stats.MaxOpenConnections)
	}
}

func TestTaskRepository_CreateAndGet(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()
	ctx := t.Context()

	task, err := repo.Create(ctx, "Build Database Layer", "Implement repositories with parameter binding")
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	if task.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if task.Status != StatusTodo {
		t.Fatalf("expected initial status TODO, got %s", task.Status)
	}

	fetched, err := repo.GetByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("failed to get task: %v", err)
	}
	if fetched.Title != task.Title {
		t.Fatalf("expected title %q, got %q", task.Title, fetched.Title)
	}
}

func TestSQLInjectionResistance(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()
	ctx := t.Context()

	// Malicious payload attempting SQL injection
	maliciousTitle := "Important Task'; DROP TABLE tasks; --"
	task, err := repo.Create(ctx, maliciousTitle, "SQL Injection attempt")
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	// Verify task title is saved literally
	fetched, err := repo.GetByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("failed to query task: %v", err)
	}
	if fetched.Title != maliciousTitle {
		t.Fatalf("expected title %q, got %q", maliciousTitle, fetched.Title)
	}

	// Verify table still exists and has tasks
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil || count == 0 {
		t.Fatalf("table tasks was compromised by injection payload: %v", err)
	}
}

func TestUpdateStatusTx_SuccessAndAudit(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()
	ctx := t.Context()

	task, err := repo.Create(ctx, "Refactor Queries", "Move to prepared statements")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Transition to IN_PROGRESS
	updated, err := repo.UpdateStatusTx(ctx, task.ID, StatusInProgress)
	if err != nil {
		t.Fatalf("update status failed: %v", err)
	}
	if updated.Status != StatusInProgress {
		t.Fatalf("expected status %s, got %s", StatusInProgress, updated.Status)
	}

	// Verify audit log has 1 record
	logs, err := repo.GetAuditLogs(ctx, task.ID)
	if err != nil {
		t.Fatalf("failed to get audit logs: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 audit log, got %d", len(logs))
	}
	if logs[0].FromStatus != StatusTodo || logs[0].ToStatus != StatusInProgress {
		t.Fatalf("unexpected audit transition: %s -> %s", logs[0].FromStatus, logs[0].ToStatus)
	}
}

func TestUpdateStatusTx_RollbackOnInvalidStatus(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()
	ctx := t.Context()

	task, _ := repo.Create(ctx, "Test Rollback", "Ensure atomicity")

	// Attempt transition with invalid status
	_, err := repo.UpdateStatusTx(ctx, task.ID, TaskStatus("INVALID_STATUS"))
	if err == nil {
		t.Fatal("expected error on invalid status, got nil")
	}

	// Verify task status remained unchanged
	fresh, _ := repo.GetByID(ctx, task.ID)
	if fresh.Status != StatusTodo {
		t.Fatalf("expected status to stay TODO, got %s", fresh.Status)
	}

	// Verify no audit log was inserted
	logs, _ := repo.GetAuditLogs(ctx, task.ID)
	if len(logs) != 0 {
		t.Fatalf("expected 0 audit logs after rollback, got %d", len(logs))
	}
}

func TestHTTPAPI_EndToEnd(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	handler := NewAPIHandler(repo)

	// 1. Healthz
	reqHealth := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recHealth := httptest.NewRecorder()
	handler.ServeHTTP(recHealth, reqHealth)
	if recHealth.Code != http.StatusOK {
		t.Fatalf("healthz failed: %d", recHealth.Code)
	}

	// 2. Create Task
	payload := CreateTaskRequest{
		Title:       "Set up Grafana dashboards",
		Description: "Include latency and error rate",
	}
	body, _ := json.Marshal(payload)
	reqCreate := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	recCreate := httptest.NewRecorder()
	handler.ServeHTTP(recCreate, reqCreate)

	if recCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", recCreate.Code, recCreate.Body.String())
	}

	var resp JSONResponse
	_ = json.Unmarshal(recCreate.Body.Bytes(), &resp)
	raw, _ := json.Marshal(resp.Data)
	var task Task
	_ = json.Unmarshal(raw, &task)

	// 3. Patch Status
	updatePayload := UpdateStatusRequest{Status: StatusDone}
	updateBody, _ := json.Marshal(updatePayload)
	reqPatch := httptest.NewRequest(http.MethodPatch, "/tasks/"+task.ID+"/status", bytes.NewReader(updateBody))
	recPatch := httptest.NewRecorder()
	handler.ServeHTTP(recPatch, reqPatch)

	if recPatch.Code != http.StatusOK {
		t.Fatalf("expected 200 on patch, got %d: %s", recPatch.Code, recPatch.Body.String())
	}

	// 4. Get Audit Logs
	reqAudit := httptest.NewRequest(http.MethodGet, "/tasks/"+task.ID+"/audit", nil)
	recAudit := httptest.NewRecorder()
	handler.ServeHTTP(recAudit, reqAudit)

	if recAudit.Code != http.StatusOK {
		t.Fatalf("expected 200 on audit, got %d", recAudit.Code)
	}

	// 5. Delete Task
	reqDelete := httptest.NewRequest(http.MethodDelete, "/tasks/"+task.ID, nil)
	recDelete := httptest.NewRecorder()
	handler.ServeHTTP(recDelete, reqDelete)

	if recDelete.Code != http.StatusNoContent {
		t.Fatalf("expected 204 on delete, got %d", recDelete.Code)
	}
}
