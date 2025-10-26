# Groceries Backend API

![CI/CD Pipeline](https://github.com/superbkibbles/groceries-backend/workflows/CI%2FCD%20Pipeline/badge.svg)

A modern e-commerce backend API built with Go, following hexagonal architecture principles.

## Features

- 🏗️ **Hexagonal Architecture** - Clean separation of concerns
- 📦 **Product Management** - Single variation products with attributes
- 🛒 **Shopping Cart** - Persistent cart management
- 📋 **Order Management** - Complete order lifecycle
- 👥 **User Management** - Authentication and user profiles
- 📚 **Categories** - Hierarchical product categorization
- ⭐ **Reviews** - Product review system
- 🔧 **Settings** - Configurable application settings
- 📖 **API Documentation** - Auto-generated Swagger docs
- 🐳 **Docker Support** - Containerized deployment
- 🔧 **Development Tools** - Live reload, linting, testing

## Quick Start

### Prerequisites

- Go 1.23 or later
- MongoDB 7.0+
- Redis 7.2+
- Docker & Docker Compose (optional)

### Setup Development Environment

```bash
# Clone the repository
git clone <repository-url>
cd groceries-backend

# Setup development environment (installs tools and dependencies)
make setup

# Start development dependencies (MongoDB & Redis)
make dev-deps

# Run database seed
make seed

# Start the application in development mode (with live reload)
make dev
```

The API will be available at `http://localhost:8080` and Swagger documentation at `http://localhost:8080/swagger/index.html`.

## Makefile Commands

### 🏗️ Build Commands

```bash
make build           # Build the main application
make build-seed      # Build the seed command
make build-all       # Build both main app and seed command
make build-linux     # Build for Linux
make build-windows   # Build for Windows
make build-mac       # Build for macOS
make build-all-platforms  # Build for all platforms
```

### 🚀 Run Commands

```bash
make run            # Run the application
make run-seed       # Run the seed command
make dev            # Run in development mode with live reload
```

### 🧪 Test Commands

```bash
make test           # Run tests
make test-coverage  # Run tests with coverage report
make test-race      # Run tests with race detection
make benchmark      # Run benchmarks
```

### 📖 Documentation Commands

```bash
make swagger        # Generate/update Swagger documentation
make swagger-serve  # Serve Swagger UI locally (app must be running)
```

### 🔧 Development Commands

```bash
make setup          # Setup development environment
make deps           # Download dependencies
make deps-update    # Update dependencies
make clean          # Clean build files
make format         # Format code
make lint           # Run linter
make vet            # Run go vet
make check          # Run format, lint, and test
```

### 🐳 Docker Commands

```bash
make docker-build   # Build Docker image
make docker-run     # Run Docker container
make docker-stop    # Stop Docker container
make dev-deps       # Start development dependencies (MongoDB, Redis)
make dev-stop       # Stop development dependencies
```

### 💾 Database Commands

```bash
make seed           # Run database seed
make run-seed       # Same as seed
make db-reset       # Reset database and seed (with confirmation)
```

### 🛠️ Utility Commands

```bash
make help           # Show all available commands
make version        # Show Go and module version
make verify         # Verify installation (build, test, swagger)
make quick-start    # Quick development workflow
make install-hooks  # Install git pre-commit hooks
```

### 🔒 Security Commands

```bash
make security-scan  # Run security vulnerability scan
```

### 📊 Profiling Commands

```bash
make profile-cpu    # CPU profiling (app must be running)
make profile-mem    # Memory profiling (app must be running)
```

### 🚢 Release Commands

```bash
make release        # Create release build for all platforms
```

## Development Workflow

### First Time Setup
```bash
make setup          # Install tools and dependencies
make dev-deps       # Start MongoDB and Redis
make seed           # Populate database with sample data
make dev            # Start development server with live reload
```

### Daily Development
```bash
make dev            # Start development server
# Make your changes - the server will automatically reload
make check          # Run format, lint, and tests before committing
```

### Before Committing
```bash
make check          # Format, lint, and test
make swagger        # Update API documentation if endpoints changed
```

### Testing
```bash
make test           # Run all tests
make test-coverage  # Generate coverage report
make test-race      # Check for race conditions
```

## Docker Development

### Using Docker Compose for Dependencies
```bash
# Start only the dependencies (recommended for development)
make dev-deps

# With management tools (Mongo Express & Redis Commander)
docker-compose --profile tools up -d

# Access management interfaces
# MongoDB: http://localhost:8081
# Redis: http://localhost:8082
```

### Full Docker Setup
```bash
# Build and run the entire application stack
make docker-build
make docker-run

# Or use docker-compose (uncomment app service in docker-compose.yml)
docker-compose up -d
```

## Configuration

The application uses environment variables for configuration. Create a `.env` file or set these variables:

```bash
# Database
MONGODB_URI=mongodb://localhost:27017/groceries_db
REDIS_URL=localhost:6379

# Server
SERVER_PORT=8080
GIN_MODE=debug

# JWT
JWT_SECRET=your-secret-key

# SMS (optional)
SMS_API_KEY=your-sms-api-key
```

## API Documentation

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI Spec**: `http://localhost:8080/swagger/doc.json`

To update the documentation after adding new endpoints:
```bash
make swagger
```

## Project Structure

```
.
├── cmd/                    # Application commands
│   └── seed/              # Database seeding command
├── docs/                  # Auto-generated Swagger documentation
├── internal/              # Private application code
│   ├── adapters/          # External interfaces (HTTP, DB, etc.)
│   ├── application/       # Application services (business logic)
│   ├── config/           # Configuration management
│   ├── domain/           # Domain entities and interfaces
│   └── utils/            # Utility functions
├── scripts/              # Database and deployment scripts
├── Dockerfile            # Container definition
├── docker-compose.yml    # Development dependencies
├── Makefile             # Build and development commands
└── .air.toml           # Live reload configuration
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make check` to ensure code quality
5. Update documentation if needed (`make swagger`)
6. Commit your changes
7. Push to your fork
8. Create a Pull Request

## Troubleshooting

### Common Issues

**Build fails with module errors:**
```bash
make deps           # Download dependencies
make deps-update    # Update to latest versions
```

**Swagger generation fails:**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
make swagger
```

**Database connection issues:**
```bash
make dev-deps       # Ensure MongoDB and Redis are running
```

**Live reload not working:**
```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Reset Everything
```bash
make clean-all      # Clean all build files and module cache
make dev-stop       # Stop development dependencies
make dev-deps       # Restart development dependencies
make setup          # Reinstall tools
make db-reset       # Reset and seed database
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
