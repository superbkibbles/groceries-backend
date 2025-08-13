# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=groceries-api
BINARY_UNIX=$(BINARY_NAME)_unix
SEED_BINARY=seed

# Swagger parameters
SWAGGER_CMD=swag

# Build directories
BUILD_DIR=build
DIST_DIR=dist

# Default target
.PHONY: all
all: test build

# Help target - shows available commands
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          Build the main application"
	@echo "  build-seed     Build the seed command"
	@echo "  build-all      Build both main app and seed command"
	@echo "  run            Run the application"
	@echo "  run-seed       Run the seed command"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage"
	@echo "  clean          Clean build files"
	@echo "  deps           Download dependencies"
	@echo "  deps-update    Update dependencies"
	@echo "  swagger        Generate/update Swagger documentation"
	@echo "  swagger-serve  Serve Swagger UI locally"
	@echo "  lint           Run linter"
	@echo "  format         Format code"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Run Docker container"
	@echo "  dev            Run in development mode with live reload"
	@echo "  setup          Setup development environment"

# Build targets
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v ./main.go

.PHONY: build-seed
build-seed:
	@echo "Building $(SEED_BINARY)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(SEED_BINARY) -v ./cmd/seed/main.go

.PHONY: build-all
build-all: build build-seed

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_UNIX) -v ./main.go

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME).exe -v ./main.go

.PHONY: build-mac
build-mac:
	@echo "Building for macOS..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)_mac -v ./main.go

.PHONY: build-all-platforms
build-all-platforms: build-linux build-windows build-mac

# Run targets
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: run-seed
run-seed: build-seed
	@echo "Running seed command..."
	./$(BUILD_DIR)/$(SEED_BINARY)

.PHONY: dev
dev:
	@echo "Running in development mode..."
	@which air > /dev/null || (echo "Installing air for live reload..." && go install github.com/cosmtrek/air@latest)
	air

# Test targets
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -v -race ./...

# Swagger documentation targets
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g main.go -o ./docs
	@echo "Swagger documentation updated!"

.PHONY: swagger-serve
swagger-serve:
	@echo "Serving Swagger UI at http://localhost:8080/swagger/index.html"
	@echo "Make sure the application is running..."

# Dependency management
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	$(GOMOD) tidy
	$(GOGET) -u ./...

.PHONY: deps-clean
deps-clean:
	@echo "Cleaning module cache..."
	$(GOCMD) clean -modcache

# Code quality targets
.PHONY: lint
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

.PHONY: format
format:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	goimports -w .

.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Clean targets
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html

.PHONY: clean-all
clean-all: clean deps-clean

# Docker targets
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name $(BINARY_NAME)-container $(BINARY_NAME):latest

.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(BINARY_NAME)-container || true
	docker rm $(BINARY_NAME)-container || true

# Development setup
.PHONY: setup
setup:
	@echo "Setting up development environment..."
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/cosmtrek/air@latest
	@echo "Downloading dependencies..."
	$(GOMOD) download
	@echo "Generating Swagger documentation..."
	swag init -g main.go -o ./docs
	@echo "Setup complete!"

# Database targets
.PHONY: seed
seed: run-seed

.PHONY: db-reset
db-reset:
	@echo "This will reset the database. Are you sure? [y/N]"
	@read ans && [ $${ans:-N} = y ]
	@echo "Resetting database and seeding..."
	$(MAKE) run-seed

# Git hooks
.PHONY: install-hooks
install-hooks:
	@echo "Installing git hooks..."
	@mkdir -p .git/hooks
	@echo '#!/bin/bash\nmake format lint test' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed!"

# Quick development workflow
.PHONY: check
check: format lint test

.PHONY: quick-start
quick-start: deps swagger build run

# Environment specific targets
.PHONY: dev-deps
dev-deps:
	@echo "Starting development dependencies (MongoDB, Redis)..."
	@command -v docker-compose >/dev/null 2>&1 || { echo "docker-compose is required but not installed. Aborting." >&2; exit 1; }
	docker-compose up -d mongodb redis

.PHONY: dev-stop
dev-stop:
	@echo "Stopping development dependencies..."
	docker-compose down

# Performance and profiling
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: profile-cpu
profile-cpu:
	@echo "CPU profiling (run this while app is running)..."
	go tool pprof http://localhost:8080/debug/pprof/profile

.PHONY: profile-mem
profile-mem:
	@echo "Memory profiling (run this while app is running)..."
	go tool pprof http://localhost:8080/debug/pprof/heap

# Security
.PHONY: security-scan
security-scan:
	@echo "Running security scan..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	gosec ./...

# Release
.PHONY: release
release: clean test lint build-all-platforms
	@echo "Release build complete!"
	@echo "Binaries available in $(DIST_DIR)/"

# Show current version
.PHONY: version
version:
	@echo "Go version: $(shell go version)"
	@echo "Module: $(shell grep '^module' go.mod | cut -d' ' -f2)"

# Verify installation
.PHONY: verify
verify:
	@echo "Verifying installation..."
	@$(MAKE) build
	@echo "✓ Build successful"
	@$(MAKE) test
	@echo "✓ Tests passed"
	@$(MAKE) swagger
	@echo "✓ Swagger generation successful"
	@echo "All checks passed!"
