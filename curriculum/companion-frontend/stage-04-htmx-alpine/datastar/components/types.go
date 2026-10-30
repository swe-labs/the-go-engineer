package components

import "time"

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
}

