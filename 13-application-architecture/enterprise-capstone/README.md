# Flagship Seed: Enterprise Capstone

## Mission

This surface is the current flagship project seed for the beta curriculum.

It is a multi-package Go service that brings together backend flow, persistence, architecture,
runtime behavior, and deployment-oriented thinking inside one longer-running system.

## Beta Stage Ownership

This project belongs to [11 Flagship](../../docs/stages/11-flagship.md).

Within the beta public shell, it is the main source surface for the flagship stage.

## Why This Project Matters

This project exists so the curriculum has one integrated system where earlier stage skills meet:

- backend request and data flow
- repository-style persistence boundaries
- package and handler structure
- runtime and deployment-oriented behavior
- longer feedback loops than a single exercise can provide

## Current Project Shape

| Area | Role |
| --- | --- |
| `Dockerfile` | packages the application for containerized execution |
| `docker-compose.yml` | runs the application with supporting services |
| `internal/` | holds the service internals and application boundaries |

## How To Use It In Beta

1. run the project and explain what each major part is responsible for
2. treat it as a staged flagship path, not as one giant final assignment
3. use checkpoint thinking before making large changes

## Run The Project

Make sure Docker is available, then run from this directory:

```bash
docker-compose up -d --build
```

That starts the current flagship seed and its supporting services.

## Example API Flow

```bash
curl -X POST -H "Content-Type: application/json" -d '{"email":"test@go.dev", "password":"password123"}' http://localhost:8080/register
curl -X POST -H "Content-Type: application/json" -d '{"email":"test@go.dev", "password":"password123"}' http://localhost:8080/login
curl -X POST -H "Authorization: Bearer <ID_FROM_LOGIN>" -H "Content-Type: application/json" -d '{"title":"Capstone", "content":"Docker is amazing"}' http://localhost:8080/posts
curl http://localhost:8080/posts
```

## Next Step

After you understand this seed, return to the
[Flagship Project stage](../../docs/stages/11-flagship.md)
and use the
[checkpoint set](../../docs/stages/flagship-project/checkpoints/README.md)
to decide what kind of work should come next.


