package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(driverName, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(5)
	return db, nil
}

func mustRowsAffected(res sql.Result) int64 {
	n, _ := res.RowsAffected()
	return n
}

func main() {
	db, err := InitDB("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Connected (pool ready)")

	res, err := db.Exec(`CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, val TEXT)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Table created: %d rows affected\n", mustRowsAffected(res))

	res, err = db.Exec(`INSERT INTO test (val) VALUES ('hello')`)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	fmt.Printf("Inserted row id=%d\n", id)

	var val string
	err = db.QueryRow(`SELECT val FROM test WHERE id = ?`, id).Scan(&val)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Queried: %s\n", val)

	rows, err := db.Query(`SELECT id, val FROM test ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var val string
		if err := rows.Scan(&id, &val); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Row: id=%d, val=%s\n", id, val)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
