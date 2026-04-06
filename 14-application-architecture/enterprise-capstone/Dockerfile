# STAGE 1: Builder
# Because `go.mod` is in the root directory, we construct the build context properly 
# from the main docker-compose, but we execute build in the Capstone context!
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Ensure we have SSL certs and gcc (if needed)
RUN apk --no-cache add ca-certificates gcc musl-dev

# Copy Go modules first for caching advantages
COPY go.mod go.sum ./
RUN go mod download

# Let's cleanly copy the entire Go codebase
COPY . .

# Build the binary explicitly from the capstone api entrypoint
RUN CGO_ENABLED=0 GOOS=linux go build -o api-binary ./14-application-architecture/enterprise-capstone/cmd/api


# STAGE 2: Runner
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Pull the binary from Phase 1
COPY --from=builder /app/api-binary .

EXPOSE 8080

CMD ["./api-binary"]
