# Multi-stage production Dockerfile for Clinic Backend
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o /bin/api ./cmd/api

# Build Worker binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o /bin/worker ./cmd/worker

# Build Migrate CLI binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o /bin/migrate ./cmd/migrate

# Production runner stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Jakarta

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /home/appuser

COPY --from=builder /bin/api /usr/local/bin/api
COPY --from=builder /bin/worker /usr/local/bin/worker
COPY --from=builder /bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/api"]

