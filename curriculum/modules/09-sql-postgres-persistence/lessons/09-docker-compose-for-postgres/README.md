# Docker Compose for Postgres

## Learning objective

Write a Docker Compose file to run PostgreSQL locally, connect to it from Go, and understand how volumes, ports, and environment variables configure the database service.

## Why this matters

Setting up PostgreSQL manually (install, service start, role creation, database creation) is error-prone and time-consuming. Docker Compose provides a repeatable, shareable, one-command setup for the exact PostgreSQL version and configuration your project needs. Every Go team uses Docker Compose for local development — it ensures every developer and CI pipeline runs against the same database environment.

## Mental model

Docker Compose is a recipe for a multi-container application. The recipe for PostgreSQL says: "Use this image, expose port 5432, mount this directory for persistent data, and set these environment variables." Running `docker compose up` executes the recipe, starting a PostgreSQL server in a container that behaves exactly like a production server but runs on your local machine.

## Core idea

A **Docker Compose** file (`compose.yaml` or `docker-compose.yml`) defines services, networks, and volumes. For PostgreSQL, the key configuration:

```yaml
services:
  postgres:
    image: postgres:17-alpine
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: myapp
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

- **image**: Which PostgreSQL version (`17-alpine` is lightweight).
- **environment**: Sets the default user, password, and database.
- **ports**: Maps host port 5432 to container port 5432.
- **volumes**: Persists database data outside the container. Without this, all data is lost when the container stops.

## Under the hood

When Docker Compose starts a PostgreSQL container:

1. Docker pulls the image if not present locally.
2. A container is created with the specified environment variables.
3. The PostgreSQL entrypoint script runs: it initializes the data directory (`/var/lib/postgresql/data`), creates the `POSTGRES_USER` role, and creates `POSTGRES_DB`.
4. The postmaster process starts listening on port 5432.
5. When the container stops, the volume `pgdata` keeps the database files. Starting a new container with the same volume restores the data.

## How Go uses it

Go connects to the Docker-hosted PostgreSQL using the host port:

```go
dsn := "postgres://app:secret@localhost:5432/myapp?sslmode=disable"
db, _ := sql.Open("pgx", dsn)
```

When running Go tests, you can start PostgreSQL via Docker Compose in a `TestMain` function or use test containers (testcontainers-go) for per-test databases.

## Go example

This example cannot start Docker, so it shows the pattern with SQLite. The `compose.yaml` and the Go connection code are provided for reference.

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// docker-compose.yaml content is shown in the README.
// In production, connect with pgx:
// db, err := sql.Open("pgx", "postgres://app:secret@localhost:5432/myapp?sslmode=disable")

func main() {
	fmt.Println("--- Docker Compose for PostgreSQL ---")
	fmt.Println()
	fmt.Println("To start PostgreSQL:")
	fmt.Println("  docker compose up -d")
	fmt.Println()
	fmt.Println("Connection string:")
	fmt.Println("  postgres://app:secret@localhost:5432/myapp?sslmode=disable")
	fmt.Println()

	// Simulate with SQLite
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)`)

	users := []string{"Alice", "Bob", "Carol"}
	for _, name := range users {
		db.Exec(`INSERT INTO users (name) VALUES (?)`, name)
	}

	rows, err := db.Query(`SELECT id, name FROM users ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Users in database:")
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %d: %s\n", id, name)
	}
}
```

## Step-by-step execution

1. The Docker Compose file defines a `postgres` service with environment variables, port mapping, and a named volume.
2. Running `docker compose up -d` starts the PostgreSQL container in the background.
3. The Go program connects to the database via `localhost:5432` (the mapped port).
4. The program creates a table, inserts sample data, and queries it — using the same `database/sql` API as SQLite.
5. When the container is stopped (`docker compose down`), the `pgdata` volume retains the data.

## Common mistakes

- **Forgetting the volume**: Without a volume, `docker compose down` destroys all data. Add a named volume to persist data.
- **Using `sslmode=disable` locally but not realizing it must change in production**: Production should use `sslmode=verify-full`. Configure this via environment variables.
- **Port conflicts**: If another PostgreSQL is running on the host, change the host port: `"5433:5432"`.
- **Not setting `POSTGRES_DB`**: Without it, the default database name matches `POSTGRES_USER`. Always set it explicitly.

## Debugging walkthrough

**Symptom**: Go application cannot connect with `dial tcp 127.0.0.1:5432: connect: connection refused`.

**Root cause**: Docker container is not running or the port mapping is wrong.

**Fix**:
```bash
docker ps                    # Verify container is running
docker compose logs postgres # Check PostgreSQL startup logs
docker compose ps            # Verify port mapping
```

If the port is mapped to `5433`:
```go
dsn := "postgres://app:secret@localhost:5433/myapp?sslmode=disable"
```

## Production notes

- For local development, use the `compose.yaml` from the project repository. Check it into version control.
- Never commit the production password to `compose.yaml`. Use `.env` files (referenced with `${POSTGRES_PASSWORD}`) or Docker secrets.
- Set resource limits in production: `deploy.resources.limits.memory: 512M`.
- Use health checks to ensure PostgreSQL is ready before connecting:

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U app"]
  interval: 5s
  timeout: 5s
  retries: 5
```

- For CI, use Docker Compose or testcontainers-go to spin up PostgreSQL per test run.

## Performance implications

- Docker containers add negligible overhead for PostgreSQL — within 1-2% of native performance.
- Named volumes on Docker Desktop (macOS/Windows) are slower than native Linux volumes due to filesystem translation. Use `gRPC FUSE` or a remote Docker host for serious performance testing.
- The Docker networking bridge adds sub-millisecond latency to localhost connections, which is insignificant for most applications.

## Practice task

1. Write a `compose.yaml` that:
   - Uses `postgres:17-alpine`.
   - Sets the user to `goapp`, password to `devpassword`, and creates a database called `batch09`.
   - Maps port `5432`.
   - Uses a named volume `pgdata`.
2. Write a Go function `WaitForDB(dsn string, maxRetries int) (*sql.DB, error)` that:
   - Attempts to connect and ping the database.
   - Retries up to `maxRetries` times with a 1-second delay between attempts.
   - Returns the connection on success.
3. Test the function with SQLite (simulating retries).

## Tests / verification

```bash
go run ./curriculum/modules/09-sql-postgres-persistence/lessons/09-docker-compose-for-postgres
go test ./curriculum/modules/09-sql-postgres-persistence/lessons/09-docker-compose-for-postgres
```

## Review questions

1. What is the purpose of the `volumes` section in a Docker Compose PostgreSQL service?
2. What does `"5432:5432"` in the ports section mean?
3. Why is `sslmode=disable` acceptable for local development but not for production?
4. How would you debug a "connection refused" error when connecting to Docker PostgreSQL?
5. What does the `POSTGRES_DB` environment variable do?

## NEXT UP

database/sql — the standard library package that abstracts SQL database access in Go.
