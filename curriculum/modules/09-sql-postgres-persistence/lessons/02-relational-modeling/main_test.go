package main

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func openMemory(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestCreateTables(t *testing.T) {
	db := openMemory(t)
	if err := createTables(db); err != nil {
		t.Fatal(err)
	}
}

func TestInsertAndJoin(t *testing.T) {
	db := openMemory(t)
	if err := createTables(db); err != nil {
		t.Fatal(err)
	}

	uid, err := insertUser(db, &User{Name: "Bob", Email: "bob@test.com", CreatedAt: time.Now()})
	if err != nil {
		t.Fatal("insert user:", err)
	}
	if uid == 0 {
		t.Fatal("expected non-zero user id")
	}

	oid, err := insertOrder(db, &Order{UserID: uid, Total: 19.99, Status: "shipped"})
	if err != nil {
		t.Fatal("insert order:", err)
	}
	if oid == 0 {
		t.Fatal("expected non-zero order id")
	}

	var name string
	var total float64
	err = db.QueryRow(`SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id WHERE o.id = ?`, oid).Scan(&name, &total)
	if err != nil {
		t.Fatal("join query:", err)
	}
	if name != "Bob" {
		t.Fatalf("got name %q, want %q", name, "Bob")
	}
	if total != 19.99 {
		t.Fatalf("got total %f, want %f", total, 19.99)
	}
}

func TestMultipleOrdersPerUser(t *testing.T) {
	db := openMemory(t)
	if err := createTables(db); err != nil {
		t.Fatal(err)
	}

	uid, err := insertUser(db, &User{Name: "Carol", Email: "carol@test.com", CreatedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		_, err := insertOrder(db, &Order{UserID: uid, Total: float64(i+1) * 10, Status: "pending"})
		if err != nil {
			t.Fatalf("insert order %d: %v", i, err)
		}
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM orders WHERE user_id = ?`, uid).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("expected 3 orders, got %d", count)
	}
}
