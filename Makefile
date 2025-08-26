# Gopher Patterns - Central Makefile
.PHONY: check test example clean help
.PHONY: test-all test-db-transaction test-db-setup test-sql-migration

# Main targets (Nova-style)
check: test-all
	@echo "🎉 All pattern checks completed!"

test: test-all

# Test all implemented patterns (db-setup must run first)
test-all: test-db-setup test-db-transaction test-sql-migration

# Individual pattern tests
test-db-transaction:
	@echo "🔄 Testing DB Transaction pattern..."
	cd db-transaction && make check

test-db-setup:
	@echo "🐘 Starting DB Setup..."
	cd db-setup && make db
	@echo "✅ DB Setup complete"

test-sql-migration:
	@echo "🚀 Testing SQL Migration pattern..."
	cd sql-migration && make check


# Show help
help:
	@echo "Gopher Patterns:"
	@echo ""
	@echo "  make check         - Test all patterns"
	@echo "  make test          - Test all patterns"
	@echo ""
	@echo "Available Patterns:"
	@echo "  🔄 db-transaction  - Context-based transaction management"
	@echo "  🐘 db-setup        - Docker PostgreSQL setup" 
	@echo "  🚀 sql-migration   - Embedded SQL migrations with Goose"