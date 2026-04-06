FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build test binary - only compile the test package, not the cmd package
RUN CGO_ENABLED=0 GOOS=linux go test -c -o /e2e-test .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /e2e-test /app/e2e-test

# Default command runs tests (use -test.v for verbose output)
CMD ["/app/e2e-test"]