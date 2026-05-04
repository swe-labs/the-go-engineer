# STAGE 1: Builder
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk --no-cache add ca-certificates gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o opslane cmd/server


# STAGE 2: Runner
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /app/opslane .

# Health check for container orchestration
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

EXPOSE 8080
CMD ["./opslane"]
