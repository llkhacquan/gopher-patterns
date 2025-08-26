# Gopher Patterns - Central Makefile
.PHONY: check test example clean help
.PHONY: test-all test-db-transaction test-db-setup

# Main targets (Nova-style)
check: test-all
	@echo "🎉 All pattern checks completed!"

test: test-all

# Test all implemented patterns (db-setup must run first)
test-all: test-db-setup test-db-transaction

# Individual pattern tests
test-db-transaction:
	@echo "🔄 Testing DB Transaction pattern..."
	cd db-transaction && make check

test-db-setup:
	@echo "🐘 Starting DB Setup..."
	cd db-setup && make db
	@echo "✅ DB Setup complete"


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
	@echo "  make check         - Test all patterns"
	@echo "  make test          - Test all patterns"
	@echo "  make example       - Run all examples"
	@echo "  make clean         - Clean all patterns"
	@echo ""
	@echo "📋 Individual Pattern Commands:"
	@echo "  make test-db-transaction       - Test specific pattern"
	@echo "  make example-db-transaction    - Run specific example"
	@echo ""
	@echo "📖 Available Patterns:"
	@echo "  🔄 db-transaction     - Context-based transaction management"
	@echo "  🐘 db-setup          - Docker PostgreSQL setup"
	@echo ""
	@echo "💡 Quick start: make check && make example"