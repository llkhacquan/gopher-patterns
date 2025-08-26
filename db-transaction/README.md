# DB Transaction Pattern

## 🎯 Problem

Managing database transactions across multiple repository operations without tight coupling between service layers and database implementation details.

**Common Issues:**
- Services need to coordinate multiple repository calls atomically  
- Repositories should work with or without transactions seamlessly
- Testing requires transaction rollback for isolation
- Hard to mock transaction behavior
- Transaction management scattered across service code

## 💡 Solution

**Context-based transaction injection** that allows repositories to automatically detect and use transactions when available, while working normally without them.

### Key Benefits
- ✅ **Transparent**: Repositories don't need to know about transactions
- ✅ **Flexible**: Works with or without transactions
- ✅ **Testable**: Easy transaction rollback in tests
- ✅ **Clean**: No transaction logic in business code
- ✅ **GORM Compatible**: Works seamlessly with GORM's transaction system

## 🔧 Implementation

The pattern uses Go's `context.Context` to inject database transactions:

```go
// Service layer - manages transactions
func (s *UserService) TransferFunds(ctx context.Context, fromID, toID string, amount decimal.Decimal) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Inject transaction into context
        ctx = transaction.SetTx(ctx, tx)
        
        // Repository methods automatically use the transaction
        if err := s.userRepo.DebitBalance(ctx, fromID, amount); err != nil {
            return err
        }
        return s.userRepo.CreditBalance(ctx, toID, amount)
    })
}

// Repository layer - automatically uses transactions when available
func (r *UserRepository) DebitBalance(ctx context.Context, userID string, amount decimal.Decimal) error {
    db := transaction.GetTx(ctx) // Gets transaction if available, otherwise regular DB
    return db.Model(&User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance - ?", amount)).Error
}
```

## 🚀 When to Use

**Perfect for:**
- Multi-repository operations that need atomicity
- Services that sometimes need transactions, sometimes don't  
- Testing repositories with automatic rollback
- Clean separation between business logic and transaction management

**Avoid when:**
- Single database operations (overkill)
- Using multiple databases (transactions don't span databases)
- Need distributed transactions (this is local only)

## ⚡ Quick Start

1. **Copy the transaction utilities**:
   ```bash
   cp transaction.go your-project/pkg/transaction/
   ```

2. **Setup in your repository**:
   ```go
   type UserRepo struct {
       db func(ctx context.Context) *gorm.DB
   }

   func NewUserRepo(db *gorm.DB) *UserRepo {
       return &UserRepo{
           db: transaction.GetTxOrDefault(db), // Use transaction if available, otherwise default DB
       }
   }
   ```

3. **Use in your service**:
   ```go
   func (s *Service) AtomicOperation(ctx context.Context) error {
       return s.db.Transaction(func(tx *gorm.DB) error {
           ctx = transaction.SetTx(ctx, tx)
           // All repository calls in this context will use the transaction
           return s.repo.SomeOperation(ctx)
       })
   }
   ```

4. **See the complete example**:
   ```bash
   go test -v -run TestBankingTransactionExample
   ```

## 📊 Tradeoffs

| Pros | Cons |
|------|------|
| Transparent to repositories | Context magic can be unclear |
| Easy testing with rollback | Only works with single database |
| Clean service code | Requires context threading |
| Flexible (works with/without tx) | GORM specific implementation |

## 🔗 Related Patterns

- **[Repository Pattern](../repository-pattern/)** - How to structure repositories that use this pattern
- **[DB Testing](../db-testing/)** - Testing strategies that leverage transaction rollback
- **[DB Codegen](../db-codegen/)** - Generated models work seamlessly with this pattern

## 📝 Notes for AI Assistants

This pattern is ideal when users need:
- "Atomic database operations across multiple tables"
- "Transaction management without coupling"  
- "Database testing with rollback"
- "Clean service layer architecture"

The key insight is using `context.Context` as a dependency injection mechanism for database transactions, allowing repositories to be transaction-aware without explicit transaction parameters.