# Makefile for InnovoSens API

.PHONY: help build run test test-unit test-integration test-coverage clean deps docker-build docker-run

# Default target
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  clean            - Clean build artifacts"
	@echo "  deps             - Install dependencies"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Run application in Docker"
	@echo "  lint             - Run linters"
	@echo "  fmt              - Format code"

# Build the application
build:
	@echo "Building InnovoSens API..."
	go build -o bin/innovasense_api ./main.go

# Run the application
run:
	@echo "Running InnovoSens API..."
	go run main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v ./tests/*_test.go -run "TestUnit|Test.*Service.*" -short

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration_test.go -run "TestIntegration"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	go tool cover -func=coverage.out | grep total:

# Run tests with detailed coverage by package
test-coverage-detailed:
	@echo "Running detailed coverage analysis..."
	go test -v -coverprofile=coverage.out -coverpkg=./... ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out
	@echo "Detailed coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/innovasense_api
	rm -f coverage.out
	rm -f coverage.html
	rm -f test_results.xml

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Run security checks
security:
	@echo "Running security checks..."
	gosec ./...

# Generate swagger docs
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g main.go -o ./docs

# Run API tests with newman (if Postman collection exists)
test-api:
	@echo "Running API tests..."
	@if [ -f "postman/innovasense_api.postman_collection.json" ]; then \
		newman run postman/innovasense_api.postman_collection.json; \
	else \
		echo "Postman collection not found. Running curl tests instead..."; \
		bash test_api.sh; \
	fi

# Performance test
test-performance:
	@echo "Running performance tests..."
	go test -bench=. -benchmem ./tests/...

# Load test with wrk (if available)
test-load:
	@echo "Running load tests..."
	@if command -v wrk >/dev/null 2>&1; then \
		wrk -t12 -c400 -d30s http://localhost:8500/health; \
	else \
		echo "wrk not installed. Install with: brew install wrk (macOS) or apt-get install wrk (Ubuntu)"; \
	fi

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t innovasense_api:latest .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8500:8500 --env-file .env innovasense_api:latest

# Database targets
db-setup:
	@echo "Setting up test database..."
	mysql -u root -p < lakshmana.sql

db-migrate:
	@echo "Running database migrations..."
	# Add migration commands here

# Development targets
dev-setup: deps db-setup
	@echo "Setting up development environment..."
	cp env.example .env
	@echo "Please edit .env file with your database configuration"

# Watch and restart on changes (requires air: go install github.com/cosmtrek/air@latest)
dev-watch:
	@echo "Starting development server with auto-reload..."
	air

# Generate test report
test-report:
	@echo "Generating test report..."
	go test -v -json ./tests/... > test_results.json
	@echo "Test results saved to test_results.json"

# All quality checks
quality: fmt vet lint security test-coverage
	@echo "All quality checks completed"

# CI/CD pipeline simulation
ci: deps quality test-coverage
	@echo "CI pipeline completed successfully"

# Production build
prod-build:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/innovasense_api ./main.go

# Health check
health:
	@echo "Checking API health..."
	curl -f http://localhost:8500/health || exit 1

# Benchmark specific functions
bench-bmi:
	@echo "Benchmarking BMI calculations..."
	go test -bench=BenchmarkBMI -benchmem ./tests/...

bench-hydration:
	@echo "Benchmarking hydration calculations..."
	go test -bench=BenchmarkHydration -benchmem ./tests/...
