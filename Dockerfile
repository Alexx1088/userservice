# --- builder stage ---
FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migrate ./cmd/migrate/main.go

# --- runtime stage ---
FROM gcr.io/distroless/base-debian11 AS runtime
WORKDIR /app
COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/bin/migrate /app/migrate
COPY migrations /app/migrations
COPY .env /app/.env
ENTRYPOINT ["/app/app"]
