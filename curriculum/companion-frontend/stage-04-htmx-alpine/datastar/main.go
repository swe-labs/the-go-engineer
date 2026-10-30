package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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

	"github.com/a-h/templ"
	"github.com/rasel9t6/the-go-engineer/curriculum/companion-frontend/stage-04-htmx-alpine/datastar/components"
)

//go:embed templates/*.html
var defaultTemplatesFS embed.FS

//go:embed static/*
var defaultStaticFS embed.FS

// Task represents a unit of work processed by the background worker pool.
type Task = components.Task

// MetricsSnapshot represents system and worker runtime stats.
type MetricsSnapshot = components.MetricsSnapshot

// DatastarSSE provides zero-dependency helpers to stream Datastar v1.0 Server-Sent Events.
type DatastarSSE struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

// NewDatastarSSE initializes SSE headers and checks for streaming support.
func NewDatastarSSE(w http.ResponseWriter) (*DatastarSSE, bool) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, false
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	return &DatastarSSE{w: w, flusher: flusher}, true
}

// PatchElements formats and streams a `datastar-patch-elements` SSE event.
// Mode can be: "outer", "inner", "prepend", "append", "before", "after", "replace".
func (d *DatastarSSE) PatchElements(elements string, selector string, mode string) error {
	var buf bytes.Buffer
	buf.WriteString("event: datastar-patch-elements\n")
	if selector != "" {
		buf.WriteString("data: selector ")
		buf.WriteString(selector)
		buf.WriteByte('\n')
	}
	if mode != "" {
		buf.WriteString("data: mode ")
		buf.WriteString(mode)
		buf.WriteByte('\n')
	}
	for _, line := range strings.Split(elements, "\n") {
		buf.WriteString("data: elements ")
		buf.WriteString(line)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')

	_, err := d.w.Write(buf.Bytes())
	if err == nil {
		d.flusher.Flush()
	}
	return err
}

// PatchComponent renders a type-safe templ.Component and streams a `datastar-patch-elements` event.
func (d *DatastarSSE) PatchComponent(ctx context.Context, c templ.Component, selector string, mode string) error {
	var buf bytes.Buffer
	if err := c.Render(ctx, &buf); err != nil {
		return err
	}
	return d.PatchElements(buf.String(), selector, mode)
}

// PatchSignals formats and streams a `datastar-patch-signals` SSE event.
func (d *DatastarSSE) PatchSignals(signalsJSON string) error {
	var buf bytes.Buffer
	buf.WriteString("event: datastar-patch-signals\n")
	buf.WriteString("data: signals ")
	buf.WriteString(signalsJSON)
	buf.WriteString("\n\n")

	_, err := d.w.Write(buf.Bytes())
	if err == nil {
		d.flusher.Flush()
	}
	return err
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

func (b *SSEBroker) BroadcastRaw(rawEvent string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.clients {
		select {
		case ch <- rawEvent:
		default:
			// Non-blocking drop if client buffer is saturated
		}
	}
}

func (b *SSEBroker) BroadcastPatchElements(elements string, selector string, mode string) {
	var buf bytes.Buffer
	buf.WriteString("event: datastar-patch-elements\n")
	if selector != "" {
		buf.WriteString("data: selector ")
		buf.WriteString(selector)
		buf.WriteByte('\n')
	}
	if mode != "" {
		buf.WriteString("data: mode ")
		buf.WriteString(mode)
		buf.WriteByte('\n')
	}
	for line := range strings.SplitSeq(elements, "\n") {
		buf.WriteString("data: elements ")
		buf.WriteString(line)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')

	b.BroadcastRaw(buf.String())
}

// BroadcastComponent renders a type-safe templ.Component and broadcasts a `datastar-patch-elements` event.
func (b *SSEBroker) BroadcastComponent(ctx context.Context, c templ.Component, selector string, mode string) error {
	var buf bytes.Buffer
	if err := c.Render(ctx, &buf); err != nil {
		return err
	}
	b.BroadcastPatchElements(buf.String(), selector, mode)
	return nil
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

	wp.seedDemoTasks()

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
		taskCopy := *task
		wp.tasksMu.Unlock()

		atomic.AddInt32(&wp.activeWorkers, 1)

		// Broadcast log update to Datastar SSE subscribers
		wp.broadcastLog(fmt.Sprintf("[Worker-%d] Started executing Task #%d (%s)...", workerID, task.ID, task.Type))

		// Broadcast live task row update (status -> RUNNING)
		wp.broadcastTaskRow(&taskCopy)
		wp.broadcastMetrics()

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
			completedCopy := *task
			wp.tasksMu.Unlock()

			wp.broadcastTaskRow(&completedCopy)

		case <-ctx.Done():
			wp.tasksMu.Lock()
			task.Status = "CANCELLED"
			task.StatusLower = "cancelled"
			delete(wp.cancelFuncs, task.ID)
			cancelledCopy := *task
			wp.tasksMu.Unlock()

			wp.broadcastLog(fmt.Sprintf("[Worker-%d] Task #%d cancelled.", workerID, task.ID))
			wp.broadcastTaskRow(&cancelledCopy)
		}

		atomic.AddInt32(&wp.activeWorkers, -1)
		cancel()

		wp.broadcastMetrics()
	}
}

func (wp *WorkerPool) broadcastLog(msg string) {
	now := time.Now().Format("15:04:05")
	if wp.templates == nil {
		var buf bytes.Buffer
		if err := components.LogEntry(now, msg, false).Render(context.Background(), &buf); err == nil {
			wp.broker.BroadcastPatchElements(buf.String(), "#log-stream", "append")
		}
		return
	}
	entryHTML := fmt.Sprintf(`<div class="log-entry"><span class="log-time font-mono">%s</span><span class="log-msg">%s</span></div>`, now, msg)
	wp.broker.BroadcastPatchElements(entryHTML, "#log-stream", "append")
}

func (wp *WorkerPool) broadcastTaskRow(task *Task) {
	selector := fmt.Sprintf("#task-%d", task.ID)
	if wp.templates == nil {
		var buf bytes.Buffer
		if err := components.TaskRow(task).Render(context.Background(), &buf); err == nil {
			wp.broker.BroadcastPatchElements(buf.String(), selector, "outer")
		}
		return
	}
	var buf bytes.Buffer
	_ = wp.templates.ExecuteTemplate(&buf, "task_row.html", task)
	wp.broker.BroadcastPatchElements(buf.String(), selector, "outer")
}

func (wp *WorkerPool) broadcastMetrics() {
	snap := wp.Snapshot()
	if wp.templates == nil {
		var buf bytes.Buffer
		if err := components.MetricsCard(snap).Render(context.Background(), &buf); err == nil {
			wp.broker.BroadcastPatchElements(buf.String(), "#metrics-panel", "outer")
		}
		return
	}
	var buf bytes.Buffer
	_ = wp.templates.ExecuteTemplate(&buf, "metrics_card.html", snap)
	wp.broker.BroadcastPatchElements(buf.String(), "#metrics-panel", "outer")
}

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

	taskCopy := *task
	return &taskCopy
}

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

	taskCopy := *task
	return &taskCopy, true
}

func (wp *WorkerPool) ListTasks() []*Task {
	wp.tasksMu.RLock()
	defer wp.tasksMu.RUnlock()

	var result []*Task
	for i := wp.nextID - 1; i >= 1; i-- {
		if t, ok := wp.tasks[i]; ok {
			tCopy := *t
			result = append(result, &tCopy)
		}
	}
	return result
}

func (wp *WorkerPool) Snapshot() MetricsSnapshot {
	wp.tasksMu.RLock()
	defer wp.tasksMu.RUnlock()

	return MetricsSnapshot{
		ActiveWorkers: int(atomic.LoadInt32(&wp.activeWorkers)),
		TotalWorkers:  wp.workerCount,
		QueueDepth:    len(wp.taskQueue),
		CompletedJobs: atomic.LoadInt64(&wp.completedJobs),
		Goroutines:    runtime.NumGoroutine(),
	}
}

// DashboardServer coordinates Datastar HTTP requests and SSE streams.
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

	metrics := s.pool.Snapshot()
	tasks := s.pool.ListTasks()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if s.templates == nil {
		_ = components.Dashboard(metrics, tasks).Render(r.Context(), w)
		return
	}

	data := struct {
		Metrics MetricsSnapshot
		Tasks   []*Task
	}{
		Metrics: metrics,
		Tasks:   tasks,
	}
	_ = s.templates.ExecuteTemplate(w, "dashboard.html", data)
}

func (s *DashboardServer) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		taskType := r.URL.Query().Get("type")
		durationSec := 3

		// Check for signals in body if not provided in URL
		if taskType == "" {
			body, _ := io.ReadAll(r.Body)
			if len(body) > 0 {
				var sigs struct {
					CustomType     string `json:"customType"`
					CustomDuration any    `json:"customDuration"`
				}
				if err := json.Unmarshal(body, &sigs); err == nil {
					if sigs.CustomType != "" {
						taskType = sigs.CustomType
					}
					switch v := sigs.CustomDuration.(type) {
					case float64:
						durationSec = int(v)
					case string:
						if d, err := strconv.Atoi(v); err == nil {
							durationSec = d
						}
					}
				}
			}
		}

		if taskType == "" {
			taskType = "Async Background Job"
		}

		if durStr := r.URL.Query().Get("duration"); durStr != "" {
			if d, err := strconv.Atoi(durStr); err == nil && d > 0 {
				durationSec = d
			}
		}

		task := s.pool.Enqueue(taskType, time.Duration(durationSec)*time.Second)

		sse, ok := NewDatastarSSE(w)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		if s.templates == nil {
			_ = sse.PatchComponent(r.Context(), components.TaskRow(task), "#tasks-table-body", "prepend")
			_ = sse.PatchComponent(r.Context(), components.MetricsCard(s.pool.Snapshot()), "#metrics-panel", "outer")
			_ = sse.PatchSignals(`{"modalOpen": false}`)
			return
		}

		// 1. Prepend new task row to table
		var taskBuf bytes.Buffer
		_ = s.templates.ExecuteTemplate(&taskBuf, "task_row.html", task)
		_ = sse.PatchElements(taskBuf.String(), "#tasks-table-body", "prepend")

		// 2. Patch metrics panel
		var metricsBuf bytes.Buffer
		snap := s.pool.Snapshot()
		_ = s.templates.ExecuteTemplate(&metricsBuf, "metrics_card.html", snap)
		_ = sse.PatchElements(metricsBuf.String(), "#metrics-panel", "outer")

		// 3. Patch client signals: automatically close the dispatch modal
		_ = sse.PatchSignals(`{"modalOpen": false}`)

	case http.MethodGet:
		tasks := s.pool.ListTasks()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if s.templates == nil {
			for _, t := range tasks {
				_ = components.TaskRow(t).Render(r.Context(), w)
			}
			return
		}
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

		sse, ok := NewDatastarSSE(w)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		selector := fmt.Sprintf("#task-%d", task.ID)
		if s.templates == nil {
			_ = sse.PatchComponent(r.Context(), components.TaskRow(task), selector, "outer")
			return
		}

		var rowBuf bytes.Buffer
		_ = s.templates.ExecuteTemplate(&rowBuf, "task_row.html", task)
		_ = sse.PatchElements(rowBuf.String(), selector, "outer")
		return
	}

	http.NotFound(w, r)
}

func (s *DashboardServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.pool.Snapshot()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if s.templates == nil {
		_ = components.MetricsCard(snap).Render(r.Context(), w)
		return
	}
	_ = s.templates.ExecuteTemplate(w, "metrics_card.html", snap)
}

// RenderComponent renders a type-safe templ.Component directly to the response writer.
func (s *DashboardServer) RenderComponent(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return c.Render(r.Context(), w)
}

func (s *DashboardServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	sse, ok := NewDatastarSSE(w)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	clientCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(clientCh)

	// Send initial connection log entry via Datastar patch elements
	now := time.Now().Format("15:04:05")
	if s.templates == nil {
		_ = sse.PatchComponent(r.Context(), components.LogEntry(now, "[Datastar] Live SSE stream established with browser", true), "#log-stream", "append")
	} else {
		initLog := fmt.Sprintf(`<div class="log-entry log-system"><span class="log-time font-mono">%s</span><span class="log-msg">[Datastar] Live SSE stream established with browser</span></div>`, now)
		_ = sse.PatchElements(initLog, "#log-stream", "append")
	}

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
			sse.flusher.Flush()
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default to 8081 so it can run alongside Stage 4's 8080
	}

	tmpl, err := template.ParseFS(defaultTemplatesFS, "templates/*.html")
	if err != nil {
		log.Fatalf("Fatal: failed to parse templates: %v", err)
	}

	broker := NewSSEBroker()
	pool := NewWorkerPool(4, broker, tmpl)
	server := NewDashboardServer(pool, broker, tmpl, defaultStaticFS)

	addr := ":" + port
	log.Printf("[Datastar Edition] Opslane Worker & Metrics Dashboard listening at http://localhost%s", addr)
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
