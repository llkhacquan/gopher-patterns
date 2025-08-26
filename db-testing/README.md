# DB Testing Pattern

## Problem

Database tests need isolation, flexibility across environments, and configurable behavior for different testing scenarios.

## Solution

Multi-environment database testing utilities with configurable options:
- **Environments**: `EnvTest` (isolated databases) and `EnvDev` (shared development DB)
- **Options**: `DBDebugOff` (clean output) and `DBNoWrapInTransaction` (skip auto-rollback)
- **Connection caching** for improved performance
- **Automatic cleanup** with `t.Cleanup()`
- **Backwards compatibility** with legacy APIs

## Quick Start

### Basic Usage
```go
func TestMyRepository(t *testing.T) {
    // Isolated database with transaction wrapping (default)
    db := CreateTestDB(t, EnvTest)
    
    db.AutoMigrate(&User{})
    
    repo := NewUserRepository(db)
    user := &User{Name: "Alice"}
    err := repo.Create(user)
    assert.NoError(t, err)
}
```

### With Options
```go
func TestWithOptions(t *testing.T) {
    // Clean output + no transaction wrapping
    db := CreateTestDB(t, EnvTest, DBDebugOff, DBNoWrapInTransaction)
    
    db.AutoMigrate(&User{})
    // Test logic here
}
```

### Development Database
```go
func TestAgainstDev(t *testing.T) {
    // Uses shared development database (may skip if unavailable)
    db := CreateTestDB(t, EnvDev, DBDebugOff)
    if db == nil {
        t.Skip("Development database not available")
        return
    }
    // Integration tests here
}
```

## Environments

### EnvTest (Recommended)
- Creates unique database per test
- Complete isolation between tests
- Automatic cleanup
- Slower startup but guaranteed clean state

### EnvDev
- Uses shared development database
- Faster startup for integration tests
- Requires external database setup
- May skip tests if database unavailable

## Options

### DBDebugOff
Disables SQL query logging for cleaner test output.

### DBNoWrapInTransaction
Skips automatic transaction wrapping when you need to test transaction logic directly.

## Migration Integration

The pattern works seamlessly with sql-migration pattern:

```go
func TestWithMigrations(t *testing.T) {
    db := CreateTestDB(t, EnvTest)
    
    // Run migrations
    migrator := NewMigrator(db)
    err := migrator.Up(context.Background())
    require.NoError(t, err)
    
    // Test against migrated schema
}
```

## When to Use Each Environment

**EnvTest**: Unit tests, repository tests, isolated testing scenarios
**EnvDev**: Integration tests, testing against realistic data, performance testing

## Connection Caching

Connections are cached for performance. Multiple `CreateTestDB` calls reuse base connections while maintaining test isolation through unique databases or transactions.

## Backwards Compatibility

Legacy functions still work:
- `CreateTestDBLegacy(t)` equivalent to `CreateTestDB(t, EnvTest)`
- `CreateTestDBWithTx(t)` equivalent to `CreateTestDB(t, EnvTest)` (default includes transaction)