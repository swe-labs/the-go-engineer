package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Product struct {
	ID    int64
	Name  string
	Price float64
}

type Queries struct {
	db DBTX
}

func NewQueries(db DBTX) *Queries {
	return &Queries{db: db}
}

func (q *Queries) CreateProduct(ctx context.Context, name string, price float64) (sql.Result, error) {
	return q.db.ExecContext(ctx, "INSERT INTO products (name, price) VALUES (?, ?)", name, price)
}

func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error) {
	row := q.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = ?", id)
	var p Product
	err := row.Scan(&p.ID, &p.Name, &p.Price)
	return p, err
}

func (q *Queries) ListProducts(ctx context.Context) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT id, name, price FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func (q *Queries) UpdateProductPrice(ctx context.Context, id int64, price float64) error {
	_, err := q.db.ExecContext(ctx, "UPDATE products SET price = ? WHERE id = ?", price, id)
	return err
}

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	return err
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL
	)`)

	q := NewQueries(db)

	q.CreateProduct(context.Background(), "Widget", 9.99)
	q.CreateProduct(context.Background(), "Gadget", 24.99)
	q.CreateProduct(context.Background(), "Doohickey", 49.99)

	p, _ := q.GetProduct(context.Background(), 1)
	fmt.Printf("Product 1: %s ($%.2f)\n", p.Name, p.Price)

	q.UpdateProductPrice(context.Background(), 1, 12.99)
	p, _ = q.GetProduct(context.Background(), 1)
	fmt.Printf("After update: %s ($%.2f)\n", p.Name, p.Price)

	q.DeleteProduct(context.Background(), 2)
	products, _ := q.ListProducts(context.Background())
	fmt.Printf("Products after delete: %d\n", len(products))
	for _, p := range products {
		fmt.Printf("  %d: %s ($%.2f)\n", p.ID, p.Name, p.Price)
	}
}
