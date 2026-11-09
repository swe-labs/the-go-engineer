package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

type Migration struct {
	Version int
	Up      string
	Down    string
}

var migrations = []Migration{
	{
		Version: 1,
		Up: `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		)`,
		Down: "DROP TABLE IF EXISTS users",
	},
	{
		Version: 2,
		Up:      "ALTER TABLE users ADD COLUMN age INTEGER DEFAULT 0",
		Down: `CREATE TABLE IF NOT EXISTS users_backup AS SELECT id, name, email FROM users;
DROP TABLE users;
CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL);
INSERT INTO users SELECT id, name, email FROM users_backup;
DROP TABLE IF EXISTS users_backup;`,
	},
}

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) ensureTable() error {
	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func (m *Migrator) CurrentVersion() (int, error) {
	if err := m.ensureTable(); err != nil {
		return 0, err
	}
	var version sql.NullInt64
	err := m.db.QueryRow("SELECT MAX(version) FROM schema_migrations").Scan(&version)
	if err != nil {
		return 0, err
	}
	if version.Valid {
		return int(version.Int64), nil
	}
	return 0, nil
}

func (m *Migrator) Up() error {
	current, err := m.CurrentVersion()
	if err != nil {
		return err
	}
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	for _, mig := range migrations {
		if mig.Version <= current {
			continue
		}
		log.Printf("Applying migration %d...\n", mig.Version)
		if err := m.apply(mig); err != nil {
			return fmt.Errorf("migration %d: %w", mig.Version, err)
		}
	}
	return nil
}

func (m *Migrator) apply(mig Migration) (err error) {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	statements := strings.Split(mig.Up, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err = tx.Exec(stmt)
		if err != nil {
			return fmt.Errorf("statement %q: %w", stmt[:clip(len(stmt), 60)], err)
		}
	}

	_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", mig.Version)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func clip(n, max int) int {
	if n > max {
		return max
	}
	return n
}

func (m *Migrator) Down() error {
	current, err := m.CurrentVersion()
	if err != nil {
		return err
	}
	if current == 0 {
		return nil
	}
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version > migrations[j].Version
	})
	for _, mig := range migrations {
		if mig.Version > current {
			continue
		}
		log.Printf("Reverting migration %d...\n", mig.Version)
		if err := m.revert(mig); err != nil {
			return fmt.Errorf("revert migration %d: %w", mig.Version, err)
		}
		break
	}
	return nil
}

func (m *Migrator) revert(mig Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := strings.Split(mig.Down, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := tx.Exec(stmt); err != nil {
			return err
		}
	}

	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = ?", mig.Version); err != nil {
		return err
	}
	return tx.Commit()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	m := NewMigrator(db)
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	var tableCount int
	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableCount)
	fmt.Printf("Users table exists: %v\n", tableCount == 1)

	var hasAge int
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull bool
		var dflt sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			log.Fatal(err)
		}
		if name == "age" {
			hasAge = 1
		}
	}
	fmt.Printf("Age column exists: %v\n", hasAge == 1)

	version, _ := m.CurrentVersion()
	fmt.Printf("Current migration version: %d\n", version)
}
