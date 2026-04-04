# Section 22: Enterprise Capstone Project

Welcome to the absolute pinnacle of **The Go Engineer**.

If you've made it this far, you are ready to construct a production-grade Backend service. We have stripped away the magic, and now we are putting the pieces together properly.

This is a multi-package, Dockerized REST API connected to a PostgreSQL database with automated startup schema migrations.

## System Architecture

We are strictly following the Standard Go Package Layout:

- **`cmd/api/main.go`**: The entry point. Handles DB connections, migrations, and shutting down securely.
- **`internal/models/`**: Domain structs (e.g., `User`, `Post`).
- **`internal/repository/`**: Abstractions over SQL queries.
- **`internal/middleware/`**: HTTP request wrappers (Logging, Secure Headers, Panic Recovery).
- **`internal/handlers/`**: The HTTP logic itself.

## How to run the entire backend stack

Ensure you have Docker Desktop installed.

1. Open your terminal in this directory (`22-enterprise-capstone`)
2. Run the command:

   ```bash
   docker-compose up -d --build
   ```

3. Watch as a PostgreSQL database spins up instantly, the Go application compiles via a Multi-Stage Dockerfile, runs its embedded `.sql` migrations automatically, and attaches to port `:8080`.

## Testing the API

```bash
# Register a user
curl -X POST -H "Content-Type: application/json" -d '{"email":"test@go.dev", "password":"password123"}' http://localhost:8080/register

# Login
curl -X POST -H "Content-Type: application/json" -d '{"email":"test@go.dev", "password":"password123"}' http://localhost:8080/login

# Create a Post
curl -X POST -H "Authorization: Bearer <ID_FROM_LOGIN>" -H "Content-Type: application/json" -d '{"title":"Capstone", "content":"Docker is amazing"}' http://localhost:8080/posts

# View Posts
curl http://localhost:8080/posts
```


## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
