# DB Testing Pattern Makefile
.PHONY: fmt test check example deps clean help

all: check

fmt:
	@echo "🔧 Formatting code..."
	go fmt ./...

test:
	@echo "🧪 Running tests..."
	go test -timeout 30s ./...

check: fmt test
	@echo "✅ All checks passed!"

example:
	@echo "📚 Running DB testing example..."
	go test -run TestRepositoryExample

deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

clean:
	@echo "🧹 Cleaning..."
	go clean -testcache

help:
	@echo "DB Testing Pattern:"
	@echo ""
	@echo "  make check         - Format + test"
	@echo "  make example       - Run repository example"