package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidStatus = errors.New("invalid status")
	ErrInvalidInput  = errors.New("invalid input")
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "TODO"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
	StatusArchived   TaskStatus = "ARCHIVED"
)

func IsValidStatus(s TaskStatus) bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone, StatusArchived:
		return true
	default:
		return false
	}
}

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskAuditLog struct {
	ID         string     `json:"id"`
	TaskID     string     `json:"task_id"`
	FromStatus TaskStatus `json:"from_status"`
	ToStatus   TaskStatus `json:"to_status"`
	ChangedAt  time.Time  `json:"changed_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateStatusRequest struct {
	Status TaskStatus `json:"status"`
}

type JSONResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
	Meta  any       `json:"meta,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func GenerateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Rebind translates standard '?' placeholders to PostgreSQL '$1, $2' if driver is postgres.
func Rebind(query string, driver string) string {
	if driver != "postgres" {
		return query
	}
	var sb strings.Builder
	paramIdx := 1
	for _, ch := range query {
		if ch == '?' {
			sb.WriteString(fmt.Sprintf("$%d", paramIdx))
			paramIdx++
		} else {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}

func ConfigurePool(db *sql.DB) {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)
}

func RunMigrations(ctx context.Context, db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id VARCHAR(64) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		status VARCHAR(32) NOT NULL DEFAULT 'TODO',
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS task_audit_logs (
		id VARCHAR(64) PRIMARY KEY,
		task_id VARCHAR(64) NOT NULL,
		from_status VARCHAR(32) NOT NULL,
		to_status VARCHAR(32) NOT NULL,
		changed_at TIMESTAMP NOT NULL
	);
	`
	_, err := db.ExecContext(ctx, schema)
	return err
}

type TaskRepository struct {
	db     *sql.DB
	driver string
}

func NewTaskRepository(db *sql.DB, driver string) *TaskRepository {
	return &TaskRepository{db: db, driver: driver}
}

func (r *TaskRepository) Create(ctx context.Context, title, description string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, fmt.Errorf("%w: title cannot be empty", ErrInvalidInput)
	}

	now := time.Now().UTC().Truncate(time.Second)
	task := Task{
		ID:          GenerateID(),
		Title:       title,
		Description: strings.TrimSpace(description),
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := Rebind(`INSERT INTO tasks (id, title, description, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`, r.driver)
	_, err := r.db.ExecContext(ctx, query, task.ID, task.Title, task.Description, string(task.Status), task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return Task{}, fmt.Errorf("failed to insert task: %w", err)
	}

	return task, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (Task, error) {
	query := Rebind(`SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?`, r.driver)
	row := r.db.QueryRowContext(ctx, query, id)

	var t Task
	var st string
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &st, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, fmt.Errorf("failed to query task: %w", err)
	}
	t.Status = TaskStatus(st)
	return t, nil
}

func (r *TaskRepository) List(ctx context.Context, status TaskStatus, limit, offset int) ([]Task, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var countQuery string
	var listQuery string
	var args []any

	if status != "" {
		countQuery = Rebind(`SELECT COUNT(*) FROM tasks WHERE status = ?`, r.driver)
		listQuery = Rebind(`SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, r.driver)
		args = append(args, string(status))
	} else {
		countQuery = `SELECT COUNT(*) FROM tasks`
		listQuery = Rebind(`SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY created_at DESC LIMIT ? OFFSET ?`, r.driver)
	}

	var total int
	var err error
	if len(args) > 0 {
		err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	} else {
		err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}

	queryArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, listQuery, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list query failed: %w", err)
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var t Task
		var st string
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &st, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}
		t.Status = TaskStatus(st)
		tasks = append(tasks, t)
	}
	return tasks, total, rows.Err()
}

// UpdateStatusTx transitions task status and writes an audit log within a single SQL transaction.
func (r *TaskRepository) UpdateStatusTx(ctx context.Context, id string, newStatus TaskStatus) (Task, error) {
	if !IsValidStatus(newStatus) {
		return Task{}, fmt.Errorf("%w: %q is not a valid status", ErrInvalidStatus, newStatus)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Task{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Query current status
	querySelect := Rebind(`SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?`, r.driver)
	row := tx.QueryRowContext(ctx, querySelect, id)

	var t Task
	var currentStatusStr string
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &currentStatusStr, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, fmt.Errorf("failed to select task in tx: %w", err)
	}
	fromStatus := TaskStatus(currentStatusStr)

	// Update task status
	now := time.Now().UTC().Truncate(time.Second)
	queryUpdate := Rebind(`UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?`, r.driver)
	if _, err := tx.ExecContext(ctx, queryUpdate, string(newStatus), now, id); err != nil {
		return Task{}, fmt.Errorf("failed to update status in tx: %w", err)
	}

	// Insert audit record
	auditID := GenerateID()
	queryAudit := Rebind(`INSERT INTO task_audit_logs (id, task_id, from_status, to_status, changed_at) VALUES (?, ?, ?, ?, ?)`, r.driver)
	if _, err := tx.ExecContext(ctx, queryAudit, auditID, id, string(fromStatus), string(newStatus), now); err != nil {
		return Task{}, fmt.Errorf("failed to insert audit log in tx: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return Task{}, fmt.Errorf("failed to commit tx: %w", err)
	}

	t.Status = newStatus
	t.UpdatedAt = now
	return t, nil
}

func (r *TaskRepository) GetAuditLogs(ctx context.Context, taskID string) ([]TaskAuditLog, error) {
	query := Rebind(`SELECT id, task_id, from_status, to_status, changed_at FROM task_audit_logs WHERE task_id = ? ORDER BY changed_at ASC`, r.driver)
	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	logs := make([]TaskAuditLog, 0)
	for rows.Next() {
		var l TaskAuditLog
		var from, to string
		if err := rows.Scan(&l.ID, &l.TaskID, &from, &to, &l.ChangedAt); err != nil {
			return nil, fmt.Errorf("audit scan failed: %w", err)
		}
		l.FromStatus = TaskStatus(from)
		l.ToStatus = TaskStatus(to)
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	query := Rebind(`DELETE FROM tasks WHERE id = ?`, r.driver)
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

type APIHandler struct {
	repo   *TaskRepository
	router *http.ServeMux
}

func NewAPIHandler(repo *TaskRepository) http.Handler {
	api := &APIHandler{
		repo:   repo,
		router: http.NewServeMux(),
	}

	api.router.HandleFunc("GET /healthz", api.handleHealthz)
	api.router.HandleFunc("POST /tasks", api.handleCreateTask)
	api.router.HandleFunc("GET /tasks", api.handleListTasks)
	api.router.HandleFunc("GET /tasks/{id}", api.handleGetTask)
	api.router.HandleFunc("PATCH /tasks/{id}/status", api.handleUpdateStatus)
	api.router.HandleFunc("GET /tasks/{id}/audit", api.handleGetAuditLogs)
	api.router.HandleFunc("DELETE /tasks/{id}", api.handleDeleteTask)

	return api.router
}

func (a *APIHandler) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if err := a.repo.db.PingContext(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, JSONResponse{
			Error: &APIError{Code: "DB_DOWN", Message: "Database is unreachable"},
		})
		return
	}
	writeJSON(w, http.StatusOK, JSONResponse{Data: map[string]string{"status": "ok", "db": "healthy"}})
}

func (a *APIHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			Error: &APIError{Code: "BAD_REQUEST", Message: "Malformed JSON payload"},
		})
		return
	}

	task, err := a.repo.Create(r.Context(), req.Title, req.Description)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeJSON(w, http.StatusUnprocessableEntity, JSONResponse{
				Error: &APIError{Code: "VALIDATION_FAILED", Message: err.Error()},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to create task"},
		})
		return
	}

	writeJSON(w, http.StatusCreated, JSONResponse{Data: task})
}

func (a *APIHandler) handleListTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	status := TaskStatus(query.Get("status"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	tasks, total, err := a.repo.List(r.Context(), status, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to list tasks"},
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		Data: tasks,
		Meta: map[string]int{"total": total, "limit": limit, "offset": offset},
	})
}

func (a *APIHandler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := a.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, JSONResponse{
				Error: &APIError{Code: "NOT_FOUND", Message: fmt.Sprintf("Task %s not found", id)},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to retrieve task"},
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{Data: task})
}

func (a *APIHandler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			Error: &APIError{Code: "BAD_REQUEST", Message: "Malformed JSON payload"},
		})
		return
	}

	task, err := a.repo.UpdateStatusTx(r.Context(), id, req.Status)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, JSONResponse{
				Error: &APIError{Code: "NOT_FOUND", Message: fmt.Sprintf("Task %s not found", id)},
			})
			return
		}
		if errors.Is(err, ErrInvalidStatus) {
			writeJSON(w, http.StatusUnprocessableEntity, JSONResponse{
				Error: &APIError{Code: "VALIDATION_FAILED", Message: err.Error()},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to update task status"},
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{Data: task})
}

func (a *APIHandler) handleGetAuditLogs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logs, err := a.repo.GetAuditLogs(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to retrieve audit logs"},
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{Data: logs})
}

func (a *APIHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := a.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, JSONResponse{
				Error: &APIError{Code: "NOT_FOUND", Message: fmt.Sprintf("Task %s not found", id)},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to delete task"},
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func main() {
	driver := "sqlite"
	dsn := "file::memory:?cache=shared"

	if envURL := os.Getenv("DATABASE_URL"); envURL != "" {
		driver = "postgres"
		dsn = envURL
	} else if envDriver := os.Getenv("DB_DRIVER"); envDriver == "postgres" {
		driver = "postgres"
		dsn = "postgres://postgres:postgrespassword@localhost:5432/taskdb?sslmode=disable"
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		slog.Error("failed to open database", "driver", driver, "error", err)
		os.Exit(1)
	}
	defer db.Close()

	ConfigurePool(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}

	if err := RunMigrations(ctx, db); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	repo := NewTaskRepository(db, driver)
	handler := NewAPIHandler(repo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("task api server running", "port", port, "driver", driver)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdown
	slog.Info("shutting down task api gracefully")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
