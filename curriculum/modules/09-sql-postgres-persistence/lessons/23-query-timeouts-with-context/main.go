package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE numbers (id INTEGER PRIMARY KEY, value TEXT NOT NULL)`)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 2000; i++ {
		db.Exec("INSERT INTO numbers (value) VALUES (?)", fmt.Sprintf("number-%d", i))
	}

	slowQuery := `SELECT COUNT(*) FROM numbers a
		CROSS JOIN numbers b
		CROSS JOIN numbers c
		CROSS JOIN numbers d`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()

	start := time.Now()
	_, err = db.QueryContext(ctx, slowQuery)
	duration := time.Since(start)

	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Printf("Query timed out after %v: %v\n", duration, err)
	} else if err != nil {
		fmt.Printf("Query error: %v\n", err)
	} else {
		fmt.Println("Query succeeded (unexpected with 1us timeout)")
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	start = time.Now()
	var count int
	err = db.QueryRowContext(ctx2, "SELECT COUNT(*) FROM numbers").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Count query completed in %v: %d rows\n", time.Since(start), count)
}
