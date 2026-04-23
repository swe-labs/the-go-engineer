# STAGE 1: Builder
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk --no-cache add ca-certificates gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o opslane ./11-flagship/01-opslane/cmd/server


# STAGE 2: Runner
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/opslane .

EXPOSE 8080
CMD ["./opslane"]
