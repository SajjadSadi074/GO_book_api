# Stage 1: Build the Go binary
FROM golang:1.23.0 AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# ✅ Build the binary with CGO disabled (static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o book-api-server .

# Stage 2: Minimal runtime container
FROM debian:bullseye-slim

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/book-api-server .

EXPOSE 8080

CMD ["./book-api-server", "serve"]
