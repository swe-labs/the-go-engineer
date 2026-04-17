# ============================================================================
# Section 20: Docker & Deployment — Multi-Stage Builds
# Level: Advanced
# ============================================================================
#
# WHAT YOU'LL LEARN:
#   - Writing an optimized Dockerfile for Go
#   - Multi-stage builds (Builder pattern)
#   - Caching go.mod dependencies for faster builds
#
# ENGINEERING DEPTH:
#   A "scratch" or "alpine" image contains virtually nothing. No bash, no curl,
#   no OS packages. Since Go binaries are statically linked (they contain everything
#   they need to run natively), we can drop them into these tiny containers.
#   This drastically reduces deployment size AND security attack surface area.
# ============================================================================

# ----------------------------------------------------------------------------
# STAGE 1: Builder
# ----------------------------------------------------------------------------
# Use the official Golang image to build the app.
# We use Alpine here as well to keep the build process lightweight.
FROM golang:1.26-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Download dependencies first. 
# We copy ONLY the go.mod and go.sum files first. 
# Why? Docker caches layers. If we copy all source code first, Docker will 
# re-download dependencies every time we change a single line of go code.
# By grabbing the mod files first, `go mod download` is cached unless dependencies change!
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the binary!
# CGO_ENABLED=0 ensures the binary is statically linked (no C libraries required)
# GOOS=linux ensures we are compiling for a Linux environment
# -o main outputs the binary to a file named 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# ----------------------------------------------------------------------------
# STAGE 2: Runner
# ----------------------------------------------------------------------------
# Start completely fresh with a tiny 5MB Alpine Linux image.
FROM alpine:latest

# Install CA certificates so our app can make secure HTTPS requests to the outside world
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# EXTREMELY IMPORTANT: We copy the binary from the 'builder' stage above!
# This leaves behind the 800MB Golang compiler and source code.
COPY --from=builder /app/main .

# Expose the port our app runs on
EXPOSE 8080

# Run the binary
CMD ["./main"]
