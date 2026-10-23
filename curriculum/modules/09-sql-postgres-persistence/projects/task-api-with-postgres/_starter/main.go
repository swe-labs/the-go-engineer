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
	// TODO: Configure connection pool limits (MaxOpenConns, MaxIdleConns, ConnMaxLifetime)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Minute)
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
	// TODO: Validate non-empty title, insert task using parameterized query with Rebind
	return Task{}, errors.New("not implemented")
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (Task, error) {
	// TODO: Query task by id using parameterized query, return ErrNotFound if sql.ErrNoRows
	return Task{}, errors.New("not implemented")
}

func (r *TaskRepository) List(ctx context.Context, status TaskStatus, limit, offset int) ([]Task, int, error) {
	// TODO: Count total rows and query paginated tasks ordered by created_at DESC
	return nil, 0, errors.New("not implemented")
}

// UpdateStatusTx transitions task status and writes an audit log within a single SQL transaction.
func (r *TaskRepository) UpdateStatusTx(ctx context.Context, id string, newStatus TaskStatus) (Task, error) {
	// TODO: Begin SQL transaction (db.BeginTx)
	// TODO: Query current task status within transaction
	// TODO: Update task status and updated_at
	// TODO: Insert row into task_audit_logs recording from_status and to_status
	// TODO: Roll back on any error, Commit if all operations succeed
	return Task{}, errors.New("not implemented")
}

func (r *TaskRepository) GetAuditLogs(ctx context.Context, taskID string) ([]TaskAuditLog, error) {
	// TODO: Query audit logs for taskID ordered by changed_at ASC
	return nil, errors.New("not implemented")
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	// TODO: Delete task by id, verify RowsAffected > 0 or return ErrNotFound
	return errors.New("not implemented")
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
