# Gopher Patterns

Go service patterns extracted from production microservice systems. Each pattern includes code, documentation, and working examples.

## Purpose

Common patterns for Go microservices:

- Database management (transactions, testing, repositories)
- gRPC service architecture (interceptors, middleware, authentication)
- Code generation (ORM models, API scaffolding)
- Testing infrastructure (isolated tests, mocking patterns)

## Pattern Catalog

| Pattern | Purpose | Complexity | Dependencies |
|---------|---------|------------|--------------|
| [DB Transaction](./db-transaction/) | Context-based transaction management | Simple | `gorm` |
| [DB Setup](./db-setup/) | Local PostgreSQL with Docker | Simple | `docker` |
| [Repository Pattern](./repository-pattern/) | Clean data access layer with testing | Medium | `gorm`, `testify` |
| [DB Testing](./db-testing/) | Isolated database testing utilities | Simple | `gorm`, `testify` |
| [gRPC Interceptors](./grpc-interceptors/) | Authentication, logging, metrics middleware | Complex | `grpc`, `zap` |
| [DB Codegen](./db-codegen/) | Automated GORM model generation | Medium | `gorm/gen` |
| [Migration Management](./migration-management/) | Database migration patterns | Medium | `goose` |

## Pattern Structure

Each pattern follows this structure:
```
pattern-name/
├── README.md           # Problem, solution, when to use
├── *.go               # Core pattern implementation
├── *_test.go          # Usage examples and unit tests
└── example_test.go     # Complete working demo as test
```

## Usage

### Testing
```bash
# Test and format all patterns
make check

# Test specific pattern
make test-db-transaction

# Run examples
make example
```

### Individual Pattern
```bash
# From pattern directory
cd db-transaction
make check        # Format + test
make example      # Run example
```

### For Developers
1. Browse the Pattern Catalog above
2. Read the pattern README to understand when to use it
3. Copy the pattern code and adapt to your needs
4. Use the example as a reference implementation

## Design Principles

- Minimal dependencies: Each pattern uses only essential external packages
- Self-contained: Every pattern can be understood and used independently  
- Copy and adapt: Patterns are meant to be copied and customized, not imported as dependencies

## Contributing

1. Pattern improvements: Open an issue with specific feedback
2. New patterns: Follow the existing pattern structure (README + code + tests + example)

## Pattern Development Checklist

- [ ] Problem definition: Clear explanation of what problem this solves
- [ ] Solution overview: How the pattern addresses the problem
- [ ] Implementation: Clean, minimal code with good abstractions
- [ ] Usage tests: Tests that demonstrate how to use the pattern
- [ ] Working example: Complete demo that can be run
- [ ] Dependencies: Minimal external dependencies, clearly listed
- [ ] Tradeoffs: When to use this pattern vs alternatives
