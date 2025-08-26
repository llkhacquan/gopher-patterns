# SQL Migration Pattern Makefile
.PHONY: fmt test check example deps clean help

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

# Check: format + test
check: fmt test
	@echo "✅ All checks passed!"

# Run the example (as test)
example:
	@echo "📚 Running SQL migration example..."
	go test -v -run TestMigrationExample

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Clean generated files
clean:
	@echo "🧹 Cleaning..."
	go clean -testcache
	rm -f *.log

# Show help
help:
	@echo "SQL Migration Pattern - Available commands:"
	@echo ""
	@echo "  make fmt           - Format code with go fmt"
	@echo "  make test          - Run tests"
	@echo "  make check         - Run fmt + test (recommended)"
	@echo "  make example       - Run the migration example"
	@echo "  make deps          - Install/update dependencies"
	@echo "  make clean         - Clean generated files"
	@echo "  make help          - Show this help"
	@echo ""
	@echo "Quick start:"
	@echo "  make check && make example"