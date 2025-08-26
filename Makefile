# Gopher Patterns - Central Makefile
.PHONY: check test fmt example clean help
.PHONY: test-all test-db-transaction test-repository-pattern test-db-testing test-grpc-interceptors test-db-codegen test-migration-management

# Main targets (Nova-style)
check: test-all
	@echo "🎉 All pattern checks completed!"

test: test-all

fmt: fmt-all

# Test all implemented patterns
test-all: test-db-transaction test-repository-pattern test-db-testing test-grpc-interceptors test-db-codegen test-migration-management

# Format all implemented patterns
fmt-all: fmt-db-transaction

fmt-db-transaction:
	@echo "🔧 Formatting DB Transaction pattern..."
	cd db-transaction && make fmt

# Individual pattern tests
test-db-transaction:
	@echo "🔄 Testing DB Transaction pattern..."
	cd db-transaction && make check

test-repository-pattern:
	@echo "🏪 Testing Repository pattern..."
	@if [ -d "repository-pattern" ]; then cd repository-pattern && go test -v; else echo "Repository pattern not implemented yet"; fi

test-db-testing:
	@echo "🧪 Testing DB Testing pattern..."
	@if [ -d "db-testing" ]; then cd db-testing && go test -v; else echo "DB Testing pattern not implemented yet"; fi

test-grpc-interceptors:
	@echo "🛡️ Testing gRPC Interceptors pattern..."
	@if [ -d "grpc-interceptors" ]; then cd grpc-interceptors && go test -v; else echo "gRPC Interceptors pattern not implemented yet"; fi

test-db-codegen:
	@echo "🏗️ Testing DB Codegen pattern..."
	@if [ -d "db-codegen" ]; then cd db-codegen && go test -v; else echo "DB Codegen pattern not implemented yet"; fi

test-migration-management:
	@echo "📦 Testing Migration Management pattern..."
	@if [ -d "migration-management" ]; then cd migration-management && go test -v; else echo "Migration Management pattern not implemented yet"; fi

# Run all examples
example: example-db-transaction

example-db-transaction:
	@echo "🏦 Running DB Transaction example..."
	cd db-transaction && make example

# Clean all patterns
clean: clean-db-transaction
	@echo "🧹 Global cleanup complete!"

clean-db-transaction:
	@echo "🧹 Cleaning DB Transaction pattern..."
	cd db-transaction && make clean

# Show help
help:
	@echo "Gopher Patterns - Central Commands:"
	@echo ""
	@echo "🎯 Main Commands (like Nova):"
	@echo "  make check         - Run format + test on all patterns"
	@echo "  make test          - Test all patterns"
	@echo "  make fmt           - Format all patterns"
	@echo "  make example       - Run all examples"
	@echo "  make clean         - Clean all patterns"
	@echo ""
	@echo "📋 Individual Pattern Commands:"
	@echo "  make test-db-transaction       - Test specific pattern"
	@echo "  make example-db-transaction    - Run specific example"
	@echo ""
	@echo "📖 Available Patterns:"
	@echo "  🔄 db-transaction     - Context-based transaction management"
	@echo "  🏪 repository-pattern - Clean data access layer (TODO)"
	@echo "  🧪 db-testing         - Isolated database testing (TODO)"
	@echo "  🛡️ grpc-interceptors  - gRPC middleware (TODO)"
	@echo "  🏗️ db-codegen         - GORM model generation (TODO)"
	@echo "  📦 migration-management - Database migrations (TODO)"
	@echo ""
	@echo "💡 Quick start: make check && make example"