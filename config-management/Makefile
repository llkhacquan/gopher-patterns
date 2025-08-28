SHELL := /bin/bash

.PHONY: fmt test check clean

fmt:
	@echo "Formatting code..."
	go fmt ./...

test:
	@echo "Running config management tests..."
	go test ./...

check: fmt test
	@echo "Config pattern check passed!"

clean:
	@echo "Cleaning build artifacts..."
	go clean ./...