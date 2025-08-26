# PATTERN_NAME Pattern Makefile
# Template for creating standardized Makefiles for each pattern
# Replace PATTERN_NAME and EXAMPLE_DESCRIPTION with actual values

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: fmt test check example clean help

# Default target
all: check

# Format code
fmt:
	@echo "🔧 Formatting code..."
	go fmt ./...

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v -timeout 30s ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -v -timeout 30s -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Check: format + test (simplified version of Nova's check)
check: fmt test
	@echo "✅ All checks passed!"

# Run the example (as test)
example:
	@echo "🎯 Running PATTERN_NAME example..."
	go test -v -run TestExamplePattern

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Clean generated files
clean:
	@echo "🧹 Cleaning..."
	rm -f coverage.out coverage.html
	rm -f *.db *.sqlite3 *.log

# Show help
help:
	@echo "PATTERN_NAME Pattern - Available commands:"
	@echo ""
	@echo "  make fmt           - Format code with go fmt"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make check         - Run fmt + test (recommended)"
	@echo "  make example       - Run the EXAMPLE_DESCRIPTION"
	@echo "  make deps          - Install/update dependencies"
	@echo "  make clean         - Clean generated files"
	@echo "  make help          - Show this help"
	@echo ""
	@echo "Quick start:"
	@echo "  make check && make example"