package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`
		CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE orders (id INTEGER PRIMARY KEY, user_id INTEGER NOT NULL, total REAL NOT NULL, FOREIGN KEY (user_id) REFERENCES users(id));
		INSERT INTO users VALUES (1, 'Alice'), (2, 'Bob'), (3, 'Carol');
		INSERT INTO orders VALUES (101, 1, 50.00), (102, 1, 30.00), (103, 3, 20.00);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestInnerJoinReturnsMatchingRows(t *testing.T) {
	db := setupTestDB(t)

	rows, err := db.Query(`SELECT orders.id, orders.user_id, orders.total, users.name
		FROM orders INNER JOIN users ON orders.user_id = users.id ORDER BY orders.id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var o OrderWithUser
		if err := rows.Scan(&o.OrderID, &o.UserID, &o.Total, &o.UserName); err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 joined rows, got %d", count)
	}
}

func TestLeftJoinShowsNullForNoMatch(t *testing.T) {
	db := setupTestDB(t)

	rows, err := db.Query(`SELECT users.id, users.name, orders.id, orders.total
		FROM users LEFT JOIN orders ON users.id = orders.user_id ORDER BY users.id, orders.id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var rowCount int
	var bobHasNullOrder bool
	for rows.Next() {
		var uwo UserWithOrder
		if err := rows.Scan(&uwo.UserID, &uwo.UserName, &uwo.OrderID, &uwo.Total); err != nil {
			t.Fatal(err)
		}
		rowCount++
		if uwo.UserName == "Bob" && !uwo.OrderID.Valid {
			bobHasNullOrder = true
		}
	}
	// Alice (2 orders) + Bob (1 null row) + Carol (1 order) = 4 rows
	if rowCount != 4 {
		t.Errorf("expected 4 rows in left join, got %d", rowCount)
	}
	if !bobHasNullOrder {
		t.Error("expected Bob to appear with NULL order")
	}
}

func TestThreeTableJoin(t *testing.T) {
	db := setupTestDB(t)

	// Add order_items and products for a 3-table join
	_, err := db.Exec(`
		CREATE TABLE products (id INTEGER PRIMARY KEY, name TEXT NOT NULL, price REAL NOT NULL);
		CREATE TABLE order_items (id INTEGER PRIMARY KEY, order_id INTEGER NOT NULL, product_id INTEGER NOT NULL, quantity INTEGER NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id), FOREIGN KEY (product_id) REFERENCES products(id));
		INSERT INTO products VALUES (1, 'Widget', 10.00), (2, 'Gadget', 25.00);
		INSERT INTO order_items VALUES (1, 101, 1, 2), (2, 101, 2, 1), (3, 102, 1, 3);
	`)
	if err != nil {
		t.Fatal(err)
	}

	type ItemWithCustomer struct {
		OrderID  int
		Customer string
		Product  string
		Quantity int
	}

	rows, err := db.Query(`
		SELECT o.id, u.name, p.name, oi.quantity
		FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN order_items oi ON oi.order_id = o.id
		JOIN products p ON p.id = oi.product_id
		ORDER BY o.id, p.id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var iwc ItemWithCustomer
		if err := rows.Scan(&iwc.OrderID, &iwc.Customer, &iwc.Product, &iwc.Quantity); err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 joined rows from 3-table join, got %d", count)
	}
}

func TestMissingJoinConditionCausesCartesianProduct(t *testing.T) {
	db := setupTestDB(t)

	// Without ON clause: Cartesian product
	rows, err := db.Query(`SELECT COUNT(*) FROM users, orders`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var count int
	rows.Next()
	rows.Scan(&count)
	// 3 users * 3 orders = 9
	if count != 9 {
		t.Errorf("expected 9 from cartesian product, got %d", count)
	}
}
