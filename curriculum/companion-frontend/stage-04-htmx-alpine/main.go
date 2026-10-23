package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//go:embed templates/*.html
var defaultTemplatesFS embed.FS

//go:embed static/*
var defaultStaticFS embed.FS

// Task represents a unit of work processed by the background worker pool.
type Task struct {
	ID          int
	Type        string
	Status      string // "QUEUED", "RUNNING", "COMPLETED", "CANCELLED"
	StatusLower string
	WorkerID    int
	Duration    time.Duration
	DurationStr string
	CreatedAt   time.Time
}

// MetricsSnapshot represents system and worker runtime stats.
type MetricsSnapshot struct {
	ActiveWorkers int
	TotalWorkers  int
	QueueDepth    int
	CompletedJobs int64
	Goroutines    int
	OOB           bool // Out-Of-Band swap flag for HTMX
}

// SSEBroker coordinates Server-Sent Events subscribers.
type SSEBroker struct {
	mu      sync.RWMutex
	clients map[chan string]bool
}

func NewSSEBroker() *SSEBroker {
	return &SSEBroker{
		clients: make(map[chan string]bool),
	}
}

func (b *SSEBroker) Subscribe() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 32)
	b.clients[ch] = true
	return ch
}

func (b *SSEBroker) Unsubscribe(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.clients[ch]; exists {
		delete(b.clients, ch)
		close(ch)
	}
}

func (b *SSEBroker) Broadcast(event, data string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	msg := fmt.Sprintf("event: %s\ndata: %s\n\n", event, strings.ReplaceAll(data, "\n", " "))
	for ch := range b.clients {
		select {
		case ch <- msg:
		default:
			// Non-blocking drop if client is saturated
		}
	}
}

// WorkerPool manages concurrent goroutines executing background tasks.
type WorkerPool struct {
	workerCount   int
	taskQueue     chan *Task
	tasksMu       sync.RWMutex
	tasks         map[int]*Task
	cancelFuncs   map[int]context.CancelFunc
	nextID        int
	activeWorkers int32
	completedJobs int64
	broker        *SSEBroker
	templates     *template.Template
}

func NewWorkerPool(workerCount int, broker *SSEBroker, tmpl *template.Template) *WorkerPool {
	wp := &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan *Task, 100),
		tasks:       make(map[int]*Task),
		cancelFuncs: make(map[int]context.CancelFunc),
		nextID:      1,
		broker:      broker,
		templates:   tmpl,
	}

	// Seed with initial demo history
	wp.seedDemoTasks()

	// Launch worker goroutines
	for i := 1; i <= workerCount; i++ {
		go wp.runWorker(i)
	}

	return wp
}

func (wp *WorkerPool) seedDemoTasks() {
	types := []string{"Data Ingestion", "Payment Reconciliation", "System Health Audit"}
	for _, t := range types {
		task := &Task{
			ID:          wp.nextID,
			Type:        t,
			Status:      "COMPLETED",
			StatusLower: "completed",
			WorkerID:    ((wp.nextID - 1) % wp.workerCount) + 1,
			Duration:    2 * time.Second,
			DurationStr: "2.0s",
			CreatedAt:   time.Now().Add(-time.Duration(wp.nextID) * time.Minute),
		}
		wp.tasks[wp.nextID] = task
		wp.nextID++
		wp.completedJobs++
	}
}

func (wp *WorkerPool) runWorker(workerID int) {
	for task := range wp.taskQueue {
		ctx, cancel := context.WithCancel(context.Background())

		wp.tasksMu.Lock()
		if task.Status == "CANCELLED" {
			wp.tasksMu.Unlock()
			cancel()
			continue
		}
		task.Status = "RUNNING"
		task.StatusLower = "running"
		task.WorkerID = workerID
		wp.cancelFuncs[task.ID] = cancel
		wp.tasksMu.Unlock()

		atomic.AddInt32(&wp.activeWorkers, 1)

		// Broadcast log update to SSE clients
		wp.broadcastLog(fmt.Sprintf("[Worker-%d] Started executing Task #%d (%s)...", workerID, task.ID, task.Type))

		// Simulate processing work
		workDone := make(chan struct{})
		go func() {
			time.Sleep(task.Duration)
			close(workDone)
		}()

		select {
		case <-workDone:
			wp.tasksMu.Lock()
			if task.Status != "CANCELLED" {
				task.Status = "COMPLETED"
				task.StatusLower = "completed"
				atomic.AddInt64(&wp.completedJobs, 1)
				wp.broadcastLog(fmt.Sprintf("[Worker-%d] Finished Task #%d successfully in %s.", workerID, task.ID, task.DurationStr))
			}
			delete(wp.cancelFuncs, task.ID)
			wp.tasksMu.Unlock()

		case <-ctx.Done():
			wp.tasksMu.Lock()
			task.Status = "CANCELLED"
			task.StatusLower = "cancelled"
			delete(wp.cancelFuncs, task.ID)
			wp.tasksMu.Unlock()

			wp.broadcastLog(fmt.Sprintf("[Worker-%d] Task #%d cancelled.", workerID, task.ID))
		}

		atomic.AddInt32(&wp.activeWorkers, -1)
		cancel()

		// Stream updated metrics to all active browsers via SSE
		wp.broadcastMetrics()
	}
}

func (wp *WorkerPool) broadcastLog(msg string) {
	now := time.Now().Format("15:04:05")
	entryHTML := fmt.Sprintf(`<div class="log-entry"><span class="log-time font-mono">%s</span><span class="log-msg">%s</span></div>`, now, msg)
	wp.broker.Broadcast("log", entryHTML)
}

func (wp *WorkerPool) broadcastMetrics() {
	var buf bytes.Buffer
	snap := wp.Snapshot(false)
	_ = wp.templates.ExecuteTemplate(&buf, "metrics_card.html", snap)
	wp.broker.Broadcast("metric", buf.String())
}

// Enqueue adds a new task to the queue and returns it.
func (wp *WorkerPool) Enqueue(taskType string, duration time.Duration) *Task {
	wp.tasksMu.Lock()
	defer wp.tasksMu.Unlock()

	task := &Task{
		ID:          wp.nextID,
		Type:        taskType,
		Status:      "QUEUED",
		StatusLower: "queued",
		Duration:    duration,
		DurationStr: fmt.Sprintf("%.1fs", duration.Seconds()),
		CreatedAt:   time.Now(),
	}
	wp.tasks[wp.nextID] = task
	wp.nextID++

	wp.taskQueue <- task
	wp.broadcastLog(fmt.Sprintf("[Queue] Enqueued new job #%d (%s, target: %s)", task.ID, task.Type, task.DurationStr))

	return task
}

// Cancel attempts to abort an active or queued task.
func (wp *WorkerPool) Cancel(id int) (*Task, bool) {
	wp.tasksMu.Lock()
	defer wp.tasksMu.Unlock()

	task, exists := wp.tasks[id]
	if !exists {
		return nil, false
	}

	if cancel, running := wp.cancelFuncs[id]; running {
		cancel()
		delete(wp.cancelFuncs, id)
	}
	task.Status = "CANCELLED"
	task.StatusLower = "cancelled"

	return task, true
}

func (wp *WorkerPool) ListTasks() []*Task {
	wp.tasksMu.RLock()
	defer wp.tasksMu.RUnlock()

	var result []*Task
	// Return in reverse chronological order
	for i := wp.nextID - 1; i >= 1; i-- {
		if t, ok := wp.tasks[i]; ok {
			result = append(result, t)
		}
	}
	return result
}

func (wp *WorkerPool) Snapshot(oob bool) MetricsSnapshot {
	wp.tasksMu.RLock()
	defer wp.tasksMu.RUnlock()

	return MetricsSnapshot{
		ActiveWorkers: int(atomic.LoadInt32(&wp.activeWorkers)),
		TotalWorkers:  wp.workerCount,
		QueueDepth:    len(wp.taskQueue),
		CompletedJobs: atomic.LoadInt64(&wp.completedJobs),
		Goroutines:    runtime.NumGoroutine(),
		OOB:           oob,
	}
}

// DashboardServer serves the HTMX + Alpine.js web interface.
type DashboardServer struct {
	pool      *WorkerPool
	broker    *SSEBroker
	templates *template.Template
	staticFS  fs.FS
}

func NewDashboardServer(pool *WorkerPool, broker *SSEBroker, tmpl *template.Template, staticFS fs.FS) *DashboardServer {
	return &DashboardServer{
		pool:      pool,
		broker:    broker,
		templates: tmpl,
		staticFS:  staticFS,
	}
}

func (s *DashboardServer) Handler() http.Handler {
	mux := http.NewServeMux()

	// Static assets
	staticSub, err := fs.Sub(s.staticFS, "static")
	if err == nil {
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))
	}

	mux.HandleFunc("/", s.handleDashboard)
	mux.HandleFunc("/tasks", s.handleTasks)
	mux.HandleFunc("/tasks/", s.handleTaskActions)
	mux.HandleFunc("/metrics", s.handleMetrics)
	mux.HandleFunc("/events", s.handleSSE)

	return mux
}

func (s *DashboardServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Metrics MetricsSnapshot
		Tasks   []*Task
	}{
		Metrics: s.pool.Snapshot(false),
		Tasks:   s.pool.ListTasks(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = s.templates.ExecuteTemplate(w, "dashboard.html", data)
}

func (s *DashboardServer) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		taskType := r.FormValue("type")
		if taskType == "" {
			taskType = "Async Background Job"
		}

		durationSec := 3
		if durStr := r.FormValue("duration"); durStr != "" {
			if d, err := strconv.Atoi(durStr); err == nil && d > 0 {
				durationSec = d
			}
		}

		task := s.pool.Enqueue(taskType, time.Duration(durationSec)*time.Second)

		// HTMX Response: Primary Partial Row + Out-of-Band (OOB) Metrics Update
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = s.templates.ExecuteTemplate(w, "task_row.html", task)

		// Append OOB swapped metrics card
		metricsSnap := s.pool.Snapshot(true)
		_ = s.templates.ExecuteTemplate(w, "metrics_card.html", metricsSnap)

	case http.MethodGet:
		tasks := s.pool.ListTasks()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		for _, t := range tasks {
			_ = s.templates.ExecuteTemplate(w, "task_row.html", t)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *DashboardServer) handleTaskActions(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")

	if len(parts) >= 2 && parts[1] == "cancel" && r.Method == http.MethodPost {
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		task, found := s.pool.Cancel(id)
		if !found {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = s.templates.ExecuteTemplate(w, "task_row.html", task)
		return
	}

	http.NotFound(w, r)
}

func (s *DashboardServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.pool.Snapshot(false)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = s.templates.ExecuteTemplate(w, "metrics_card.html", snap)
}

func (s *DashboardServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(clientCh)

	// Send initial connection event
	_, _ = fmt.Fprintf(w, "event: log\ndata: <div class=\"log-entry log-system\"><span class=\"log-time font-mono\">%s</span><span class=\"log-msg\">[SSE] Live stream established with browser</span></div>\n\n", time.Now().Format("15:04:05"))
	flusher.Flush()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg, open := <-clientCh:
			if !open {
				return
			}
			_, _ = fmt.Fprint(w, msg)
			flusher.Flush()
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tmpl, err := template.ParseFS(defaultTemplatesFS, "templates/*.html")
	if err != nil {
		log.Fatalf("Fatal: failed to parse templates: %v", err)
	}

	broker := NewSSEBroker()
	pool := NewWorkerPool(4, broker, tmpl)
	server := NewDashboardServer(pool, broker, tmpl, defaultStaticFS)

	addr := ":" + port
	log.Printf("[Stage 4 Starter] Opslane Worker & Metrics Dashboard listening at http://localhost%s", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      server.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server stopped: %v", err)
	}
}
