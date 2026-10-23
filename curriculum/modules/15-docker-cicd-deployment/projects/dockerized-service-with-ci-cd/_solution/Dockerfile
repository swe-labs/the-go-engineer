# Stage 1: Build static Go binary with cached layers
FROM golang:1.24-alpine AS builder

# Install CA certificates and create unprivileged non-root user
RUN apk update && apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' -u 10001 appuser

WORKDIR /app

# Cache dependency download layer
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and compile optimized static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X main.Version=v1.0.0 -X main.CommitSHA=release" -o /bin/server .

# Stage 2: Minimal scratch image containing zero attack surface
FROM scratch

# Import CA certs and unprivileged user from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy compiled static binary
COPY --from=builder /bin/server /bin/server

# Run as non-root user
USER appuser:appuser

# Expose HTTP port
EXPOSE 8080

# Run internal healthcheck probe
ENTRYPOINT ["/bin/server"]
