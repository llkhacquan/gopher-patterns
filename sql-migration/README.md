# SQL Migration Pattern

## Problem

Database schema changes need to be versioned and reproducible across environments. Teams need to reliably recreate the same database state.

## Solution

Embedded SQL migrations using Goose:
- SQL files bundled into Go binary with `go:embed`
- Automatic versioning and rollback support
- Reproduces exact database state across dev/test/prod
- No external migration files to manage

## How It Works

1. **Write migrations** in `migrations/001_name.sql` format
2. **Embed files** automatically with `go:embed`
3. **Run migrations** with `migrator.Up()` to reach target state
4. **Track progress** - Goose maintains version table

Each environment runs the same embedded migrations → same database state.

## Quick Start

```go
migrator, err := NewMigrator(Config{
    Host: "localhost", Port: 5432, 
    User: "postgres", Password: "password", Database: "myapp"
})

// Reproduce database state
err = migrator.Up(context.Background())
```

## Migration Format

```sql
-- migrations/001_create_users.sql
-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- +goose Down  
DROP TABLE IF EXISTS users;
```

## Real-World Benefits

**Production scenarios where this pattern helps:**

- **New developer onboarding**: `make db && migrator.Up()` recreates entire production schema locally
- **CI/CD deployments**: Same migrations run in test/staging/prod → guaranteed schema consistency  
- **Feature branches**: Each branch carries its schema changes, no external file conflicts
- **Database rollbacks**: `migrator.Down()` safely reverts problematic schema changes
- **Docker deployments**: Single binary contains both app and all schema versions
- **Multi-environment**: Dev/staging/prod all reach identical database state from same codebase

**Nova example**: Trading engine deploys with embedded migrations. Every environment (local dev, CI, staging, prod) runs identical `UpAllPendingScripts()` → same schema version → reliable trading operations.

## When to Use

Essential for production systems requiring reproducible database state across environments.