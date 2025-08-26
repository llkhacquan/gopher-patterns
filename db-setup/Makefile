# DB Setup Pattern Makefile
.PHONY: db

# Start PostgreSQL if not running, ensure postgres database exists
db:
	@if nc -z localhost 5432 2>/dev/null; then \
		echo "✅ PostgreSQL is already running"; \
		pgcli postgres://postgres:password@localhost:5432/postgres -c "SELECT 1" >/dev/null 2>&1 || \
		pgcli postgres://postgres:password@localhost:5432/template1 -c "CREATE DATABASE postgres" >/dev/null 2>&1 || true; \
	else \
		echo "🚀 Starting PostgreSQL..."; \
		docker-compose up -d; \
		echo "⏳ Waiting for PostgreSQL..."; \
		until nc -z localhost 5432; do sleep 1; done; \
		echo "✅ PostgreSQL ready!"; \
	fi

db.connect:
	@echo "Connecting to PostgreSQL..."
	@pgcli postgres://postgres:password@localhost:5432/postgres
