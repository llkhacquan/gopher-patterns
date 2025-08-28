# DB Transaction Pattern Makefile
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
	go test -timeout 30s ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -timeout 30s -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Check: format + test (simplified version of Nova's check)
check: fmt test
	@echo "✅ All checks passed!"

# Run the example (now part of tests)
example:
	@echo "🏦 Running banking transaction example..."
	go test -run TestBankingTransactionExample

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Clean generated files
clean:
	@echo "🧹 Cleaning..."
	rm -f coverage.out coverage.html
	rm -f *.db *.sqlite3

# Show help
help:
	@echo "DB Transaction Pattern - Available commands:"
	@echo ""
	@echo "  make fmt           - Format code with goimports/go fmt"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make check         - Run fmt + test (recommended)"
	@echo "  make example       - Run the banking transaction example"
	@echo "  make deps          - Install/update dependencies"
	@echo "  make clean         - Clean generated files"
	@echo "  make help          - Show this help"
	@echo ""
	@echo "Quick start:"
	@echo "  make check && make example"