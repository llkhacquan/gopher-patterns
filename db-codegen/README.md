# Database Code Generation

Automated GORM code generation that creates type-safe models and query builders from database schemas.

## Overview

This pattern creates a temporary database, applies schema definitions, and uses GORM Gen to generate Go models and query builders. The generated code provides type-safe database operations without manual model maintenance.

## Architecture

```go
type codeGenerator struct {
    connString string
    tempDB     string
}

func (c *codeGenerator) Run() error {
    // Create temporary database
    // Apply schema definitions
    // Generate GORM models and queries
    // Cleanup resources
}
```

## Dependencies

- PostgreSQL database (configured via db-setup)
- GORM Gen for code generation
- Direct SQL for schema creation

## Usage

```bash
make gen
```

This creates:
- `model/*.gen.go` - GORM model structs
- `query/*.gen.go` - Type-safe query builders

## Generated Code Usage

```go
import "db-codegen/query"

q := query.Use(db)
users, err := q.User.Where(q.User.Email.Like("%@example.com")).Find()
```

## Configuration

The generator creates tables for users and orders with standard fields:
- Primary keys (BIGSERIAL)
- Timestamps (created_at, updated_at)
- Business fields (name, email, product, etc.)

Database connection uses localhost PostgreSQL with credentials from db-setup configuration.