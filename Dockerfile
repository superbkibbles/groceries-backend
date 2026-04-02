# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
RUN swag init -g main.go -o ./docs

# Build API, seed, and repair-categories binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./cmd/seed/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o repair-categories ./cmd/repair-categories/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata wget

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /home/appuser

# Copy binaries from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/seed .
COPY --from=builder /app/repair-categories .

# Copy any additional files
COPY --from=builder /app/docs ./docs

# Change ownership to appuser
RUN chown -R appuser:appuser /home/appuser
USER appuser

# Expose port
EXPOSE 80

# Health check
# HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
#     CMD wget --no-verbose --tries=1 --spider http://localhost:80/api/v1/health || exit 1

# Run the application
CMD ["./main"]
