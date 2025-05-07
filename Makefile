.PHONY: all build test bench clean

# Default target
all: build

# Build the project
build:
	@echo "Building..."
	go build -o gendo ./cmd/gendo

# Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f gendo
	rm -f coverage.out

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Help target
help:
	@echo "Available targets:"
	@echo "  all        - Build the project (default)"
	@echo "  build      - Build the project"
	@echo "  test       - Run all tests"
	@echo "  bench      - Run benchmarks"
	@echo "  coverage   - Run tests with coverage report"
	@echo "  clean      - Remove build artifacts"
	@echo "  deps       - Install dependencies"
	@echo "  lint       - Run linter"
	@echo "  help       - Show this help message" 