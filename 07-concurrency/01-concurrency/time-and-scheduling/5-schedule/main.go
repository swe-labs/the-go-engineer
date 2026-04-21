// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Stage 07: Time & Scheduling — Cron-like Scheduling
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Building a custom task scheduler from scratch
//   - Managing state spanning multiple goroutines with `sync.Mutex`
//   - Waiting for massive concurrency trees to resolve via `sync.WaitGroup`
//   - Safe cancellation using channel blocks
//
// ENGINEERING DEPTH:
//   This custom Scheduler implements an Actor Model pattern. Instead of spawning
//   cron-jobs globally and losing control of them, the `Scheduler` struct holds
//   a map of Task pointers and a global `WaitGroup`. If the application needs to
//   shut down (e.g. graceful shutdown on SIGTERM), it loops through the map, sends
//   a signal to every child's `stopChan`, and safely blocks until the `WaitGroup`
//   counter hits zero. This guarantees zero zombie processes during deployments!
//
// RUN: go run ./07-concurrency/01-concurrency/time-and-scheduling/5-schedule
// ============================================================================

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID          string
	Name        string
	IsRecurring bool
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

type Scheduler struct {
	tasks    map[string]*Task
	mu       sync.Mutex
	globalWg sync.WaitGroup
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*Task),
	}
}

func generateTaskID() string {
	return fmt.Sprintf("task-%d", time.Now().UnixNano())
}

func (s *Scheduler) ScheduleOnce(name string, delay time.Duration, action func()) string {
	taskID := generateTaskID()

	fmt.Printf("[%s] SCHEDULER: Scheduling one-off task '%s' (ID: %s) to run after %s\n",
		time.Now().Format("15:04:05.000"), name, taskID, delay)

	s.globalWg.Add(1)
	time.AfterFunc(delay, func() {
		defer s.globalWg.Done()
		fmt.Printf("[%s] TASK '%s' (ID: %s): Executing one-off action.\n",
			time.Now().Format("15:04:05.000"), name, taskID)
		action()
		fmt.Printf("[%s] TASK '%s' (ID: %s): Finished one-off action.\n",
			time.Now().Format("15:04:05.000"), name, taskID)
	})

	return taskID
}

func (s *Scheduler) ScheduleInterval(name string, initialDelay time.Duration, interval time.Duration, action func()) (string, error) {
	if interval <= 0 {
		return "", fmt.Errorf("interval must be positive, got %v", interval)
	}

	taskID := generateTaskID()

	task := &Task{
		ID:          taskID,
		Name:        name,
		IsRecurring: true,
		stopChan:    make(chan struct{}),
	}

	// 1. Thread-Safe Map Mutations
	// Maps in Go are NOT thread-safe. If two Goroutines try to schedule a task
	// at exactly the same microsecond, the application will trigger a fatal panic.
	// We MUST acquire a Lock, write to the map memory, and immediately Unlock it.
	s.mu.Lock()
	s.tasks[taskID] = task
	s.mu.Unlock()

	fmt.Printf("[%s] SCHEDULER: Scheduling interval task '%s' (ID: %s) with initial delay %s, interval %s\n",
		time.Now().Format("15:04:05.000"), name, taskID, initialDelay, interval)

	s.globalWg.Add(1)
	task.wg.Add(1)
	go func() {
		defer s.globalWg.Done()
		defer task.wg.Done()

		fmt.Printf("[%s] TASK '%s' (ID: %s): Goroutine started. Initial delay: %s.\n",
			time.Now().Format("15:04:05.000"), task.Name, task.ID, initialDelay)

		initTimer := time.NewTimer(initialDelay)
		select {
		case <-initTimer.C:
		// allow it to flow down
		case <-task.stopChan:
			initTimer.Stop()
			fmt.Printf("[%s] TASK '%s' (ID: %s): Stopped before initial run.\n",
				time.Now().Format("15:04:05.000"), task.Name, task.ID)
			return
		}

		fmt.Printf("[%s] TASK '%s' (ID: %s): Executing first action.\n",
			time.Now().Format("15:04:05.000"), task.Name, task.ID)
		action()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		fmt.Printf("[%s] TASK '%s' (ID: %s): Ticker started, interval %s.\n",
			time.Now().Format("15:04:05.000"), task.Name, task.ID, interval)
		for {
			select {
			case t := <-ticker.C:
				fmt.Printf("[%s] TASK '%s' (ID: %s): Executing recurring action at %v.\n",
					time.Now().Format("15:04:05.000"), task.Name, task.ID, t.Format("15:04:05.000"))
				action()
			case <-task.stopChan:
				// 2. Asynchronous Cancellation
				// If `StopAll()` closes `task.stopChan`, this case instantly
				// becomes unblocked, terminating the infinite `for` loop and
				// allowing the `defer s.globalWg.Done()` to execute.
				fmt.Printf("[%s] TASK '%s' (ID: %s): Received stop signal. Exiting goroutine.\n",
					time.Now().Format("15:04:05.000"), task.Name, task.ID)
				return
			}
		}

	}()

	return taskID, nil
}

func (s *Scheduler) StopTask(taskID string) bool {
	s.mu.Lock()
	task, exists := s.tasks[taskID]
	if !exists || !task.IsRecurring {
		s.mu.Unlock()
		fmt.Printf("[%s] SCHEDULER: Task ID '%s' not found or not a stoppable recurring task.\n",
			time.Now().Format("15:04:05.000"), taskID)
		return false
	}
	delete(s.tasks, taskID)
	s.mu.Unlock()
	fmt.Printf("[%s] SCHEDULER: Sending stop signal to task '%s' (ID: %s).\n",
		time.Now().Format("15:04:05.000"), task.Name, taskID)
	close(task.stopChan) // am assumption that this would be closed once

	task.wg.Wait()
	fmt.Printf("[%s] SCHEDULER: Task '%s' (ID: %s) confirmed stopped.\n",
		time.Now().Format("15:04:05.000"), task.Name, taskID)
	return true
}

func (s *Scheduler) StopAll() {
	fmt.Printf("[%s] SCHEDULER: Initiating shutdown of all tasks...\n", time.Now().Format("15:04:05.000"))
	s.mu.Lock()

	var taskIDsToStop []string
	for id := range s.tasks {
		taskIDsToStop = append(taskIDsToStop, id)
	}
	s.mu.Unlock()

	for _, id := range taskIDsToStop {
		s.StopTask(id)
	}
	fmt.Printf("[%s] SCHEDULER: Waiting for all task goroutines to complete...\n", time.Now().Format("15:04:05.000"))
	s.globalWg.Wait()

	fmt.Printf("[%s] SCHEDULER: All tasks shut down. Scheduler stopped.\n", time.Now().Format("15:04:05.000"))
}

func (s *Scheduler) AutoStopAll(after time.Duration) {
	go func() {
		timer := time.NewTimer(after)
		defer timer.Stop()
		<-timer.C
		fmt.Printf("[%s] SCHEDULER: AutoStopAll triggered after %s.\n",
			time.Now().Format("15:04:05.000"), after)
		s.StopAll()
	}()
}

func main() {

	scheduler := NewScheduler()
	scheduler.ScheduleOnce("test", 2*time.Second, func() {
		fmt.Println("Running test...")
	})

	_, _ = scheduler.ScheduleInterval("test-interval", 1*time.Second, 3*time.Second, func() {
		fmt.Println("Running database backup in realtime")
	})

	//scheduler.AutoStopAll(5 * time.Second)
	//time.AfterFunc(5*time.Second, scheduler.StopAll)

	scheduler.globalWg.Wait()
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: TM.6 timezones")
	fmt.Println("   Current: TM.5 (scheduler)")
	fmt.Println("---------------------------------------------------")
}
