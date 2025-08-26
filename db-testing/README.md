# DB Testing Pattern

## Problem

Database tests need isolation and fast setup without manual database management.

## Solution

Simple test database utilities that create isolated databases for each test:
- `CreateTestDB()` - Clean database per test
- `CreateTestDBWithTx()` - Transaction-wrapped for rollback
- Automatic cleanup with `t.Cleanup()`
- Works with db-setup pattern (PostgreSQL on localhost:5432)

## Quick Start

```go
func TestMyRepository(t *testing.T) {
    db := CreateTestDB(t)
    
    // Auto-migrate your models
    db.AutoMigrate(&User{})
    
    // Test your repository
    repo := NewUserRepository(db)
    user := &User{Name: "Alice"}
    err := repo.Create(user)
    assert.NoError(t, err)
}
```

## Transaction Isolation

```go
func TestWithTransaction(t *testing.T) {
    tx := CreateTestDBWithTx(t) // Rolls back automatically
    
    // All changes are isolated and cleaned up
    tx.Create(&User{Name: "Test"})
}
```

## When to Use

Essential for repository testing with real PostgreSQL. Each test gets clean database state.