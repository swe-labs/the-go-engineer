# UNOPTIMIZED STARTER DOCKERFILE
# Flaws in this starter:
# 1. Single stage includes entire Go toolchain in final artifact (~800MB vs ~15MB)
# 2. Runs as root user (security hazard)
# 3. Source code copied before go mod download (breaks layer caching)
#
# TODO: Task 1 - Refactor into multi-stage build:
#   Stage 1: FROM golang:1.24-alpine AS builder (create non-root appuser, copy go.mod first, build static CGO_ENABLED=0 binary)
#   Stage 2: FROM scratch (copy ca-certificates, /etc/passwd, binary, USER appuser, ENTRYPOINT)

FROM golang:1.24-alpine

WORKDIR /app
COPY . .
RUN go build -o /bin/server .

EXPOSE 8080
CMD ["/bin/server"]
