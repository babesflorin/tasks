FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies for building
RUN apk add --no-cache git

# Copy go mod file and source code
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

# Initialize go module and download dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /server /app/server

# Expose port
EXPOSE 8080

# Default command runs the server
CMD ["/app/server"]