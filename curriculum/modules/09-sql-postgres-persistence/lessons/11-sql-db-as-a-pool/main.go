package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

func configurePool(db *sql.DB) {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(30 * time.Second)
	db.SetConnMaxIdleTime(10 * time.Second)
}

func MonitorPool(db *sql.DB) *sql.DBStats {
	stats := db.Stats()
	if stats.WaitCount > 100 {
		log.Printf("WARNING: pool contention — %d waits", stats.WaitCount)
	}
	return &stats
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	configurePool(db)

	db.Exec(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, ts TEXT, msg TEXT)`)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				_, err := db.Exec(`INSERT INTO events (ts, msg) VALUES (datetime('now'), ?)`,
					fmt.Sprintf("event %d-%d", n, j))
				if err != nil {
					log.Printf("insert error: %v", err)
				}
			}
		}(i)
	}
	wg.Wait()

	stats := MonitorPool(db)
	fmt.Printf("Pool stats:\n")
	fmt.Printf("  MaxOpen:     %d\n", stats.MaxOpenConnections)
	fmt.Printf("  Open:        %d\n", stats.OpenConnections)
	fmt.Printf("  InUse:       %d\n", stats.InUse)
	fmt.Printf("  Idle:        %d\n", stats.Idle)
	fmt.Printf("  WaitCount:   %d\n", stats.WaitCount)
	fmt.Printf("  WaitDuration: %v\n", stats.WaitDuration)

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM events`).Scan(&count)
	fmt.Printf("Total events: %d\n", count)
}
