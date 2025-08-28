SHELL := /bin/bash

.PHONY: fmt test gen check clean

fmt:
	@echo "Formatting code..."
	go fmt ./...

test:
	@echo "Running tests..."
	go test -timeout 30s ./...

gen:
	@echo "Running GORM code generation..."
	go mod tidy
	go run .

check: fmt test
	@echo "All checks passed!"

clean:
	@echo "Cleaning generated files..."
	rm -rf query model