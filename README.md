# Geolocation Service

A RESTful API service built in Go that enables users to register geolocated stations with latitude and longitude coordinates, and query for the nearest station to any given point.

## Features

- Register new locations with name, latitude, and longitude
- Find the nearest location to a given point
- List all registered locations
- Delete locations by name
- Auto-generated OpenAPI documentation at `/docs`
- Support for both in-memory and PostgreSQL storage

## Setup Instructions

### Using Docker Compose (Recommended)

```bash
# Start the service with PostgreSQL database
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

### Local Development

```bash
# Install dependencies
go mod download

# Run with in-memory storage
go run ./cmd/api

# Run with PostgreSQL (ensure database is running)
export STORAGE_TYPE=postgres
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=geolocation
go run ./cmd/api
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|----------|
| `SERVER_PORT` | HTTP server port | `8080` |
| `STORAGE_TYPE` | Storage type: "memory" or "postgres" | `memory` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_USER` | PostgreSQL username | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `postgres` |
| `DB_NAME` | PostgreSQL database name | `geolocation` |

## API Usage Examples

### API Documentation
Interactive API documentation is available at `http://localhost:8080/docs` when the service is running.

### Using curl

```bash
# Register a new location
curl -X POST http://localhost:8080/locations \
  -H "Content-Type: application/json" \
  -d '{"name":"Central Park","latitude":40.7829,"longitude":-73.9654}'

# List all locations
curl http://localhost:8080/locations

# Find nearest location
curl "http://localhost:8080/nearest?lat=40.7589&lng=-73.9851"

# Find nearest with specific unit
curl "http://localhost:8080/nearest?lat=40.7589&lng=-73.9851&unit=miles"

# Delete a location
curl -X DELETE "http://localhost:8080/locations/Central%20Park"
```

## How to Run Tests

### Prerequisites
- Go 1.21 or later
- PostgreSQL running (for integration tests)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run only unit tests (exclude integration tests)
go test -short ./...

# Run specific test
go test -run TestLocationService ./internal/service

# Run tests with verbose output
go test -v ./...
```

### Test Categories
- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test database interactions and API endpoints
- **Performance Tests**: Benchmark spatial queries and API performance

### Test Database Setup
Integration tests use a separate test database. Ensure PostgreSQL is running and accessible with the environment variables set in your `.env` file.

## Running the Service

### Using Docker Compose (Recommended)

```bash
# Start the service with PostgreSQL database
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

### Using Docker

```bash
# Build the image
docker build -t geolocation-service .

# Run with in-memory storage
docker run -p 8080:8080 geolocation-service

# Run with environment variables
docker run -p 8080:8080 \
  -e STORAGE_TYPE=memory \
  -e SERVER_PORT=8080 \
  geolocation-service
```

### Local Development

```bash
# Install dependencies
go mod download

# Run with in-memory storage
go run ./cmd/api

# Run with PostgreSQL (ensure database is running)
export STORAGE_TYPE=postgres
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=geolocation
go run ./cmd/api
```

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_PORT` | HTTP server port | `8080` | No |
| `STORAGE_TYPE` | Storage type: "memory" or "postgres" | `memory` | No |
| `DB_HOST` | PostgreSQL host | `localhost` | If using postgres |
| `DB_PORT` | PostgreSQL port | `5432` | No |
| `DB_USER` | PostgreSQL username | `postgres` | If using postgres |
| `DB_PASSWORD` | PostgreSQL password | `postgres` | If using postgres |
| `DB_NAME` | PostgreSQL database name | `geolocation` | If using postgres |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` | No |

## Development

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (optional, for database storage)
- Docker and Docker Compose (optional, for containerized deployment)

### Building

```bash
# Build the binary
go build -o geolocation-service ./cmd/api

# Run the binary
./geolocation-service
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test package
go test ./internal/service
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Vet code
go vet ./...
```

## Architecture

The service follows a clean architecture pattern with the following layers:

- **Handlers**: HTTP request/response handling using Huma.rocks
- **Service**: Business logic and validation
- **Repository**: Data persistence abstraction
- **Domain**: Core business entities and interfaces

### Project Structure

```
.
├── cmd/api/                 # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── domain/             # Domain entities and interfaces
│   ├── handlers/           # HTTP handlers
│   ├── repository/         # Data persistence layer
│   │   ├── memory/         # In-memory implementation
│   │   └── postgres/       # PostgreSQL implementation
│   └── service/            # Business logic layer
├── pkg/geospatial/         # Geospatial utilities
├── tests/                  # Integration tests
├── docs/                   # Additional documentation
├── docker-compose.yml      # Docker Compose configuration
├── Dockerfile             # Docker image definition
└── README.md              # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.