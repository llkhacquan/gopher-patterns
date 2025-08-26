package transaction

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ctxKey is used to store the transaction in the context
var ctxKey = new(int)

// Fix creates a database function that always uses the provided database instance
// Useful when you want to force using a specific DB connection (e.g., in tests)
func Fix(db *gorm.DB) func(ctx context.Context) *gorm.DB {
	return func(ctx context.Context) *gorm.DB {
		return db.WithContext(ctx)
	}
}

// GetTxOrDefault creates a database function that uses a transaction if available in context,
// otherwise falls back to the provided default database
// This is the most common pattern for repositories
func GetTxOrDefault(defaultDB *gorm.DB) func(ctx context.Context) *gorm.DB {
	return func(ctx context.Context) *gorm.DB {
		if tx := GetTx(ctx); tx != nil {
			return tx.WithContext(ctx)
		}
		return defaultDB.WithContext(ctx)
	}
}

// selectForUpdateKey is used to store SELECT FOR UPDATE preference in context
var selectForUpdateKey = new(int)

// IsSelectForUpdate checks if the context has SELECT FOR UPDATE enabled
func IsSelectForUpdate(ctx context.Context) bool {
	if v := ctx.Value(selectForUpdateKey); v != nil {
		return v.(bool)
	}
	return false
}

// SelectForUpdate creates a context with SELECT FOR UPDATE enabled
// This will cause queries to lock rows for update
func SelectForUpdate(ctx context.Context) context.Context {
	return context.WithValue(ctx, selectForUpdateKey, true)
}

// GetTx retrieves the transaction from the context
// Returns nil if no transaction is set
func GetTx(ctx context.Context) *gorm.DB {
	if tx := ctx.Value(ctxKey); tx != nil {
		if db := tx.(*gorm.DB); db != nil {
			// Apply SELECT FOR UPDATE if context requests it
			if IsSelectForUpdate(ctx) {
				return db.Clauses(clause.Locking{Strength: "UPDATE"})
			}
			return db
		}
	}
	return nil
}

// SetTx stores a transaction in the context
// This is typically called by the service layer when starting a transaction
func SetTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey, tx)
}

// SetTxFunc stores a transaction function in the context
// Alternative approach that stores a function instead of the transaction directly
func SetTxFunc(ctx context.Context, txFunc func(ctx context.Context) *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey, fromTxFunc(txFunc)(ctx))
}

// fromTxFunc converts a transaction function to a transaction
// Internal helper function
func fromTxFunc(txFunc func(ctx context.Context) *gorm.DB) func(ctx context.Context) *gorm.DB {
	return func(ctx context.Context) *gorm.DB {
		if tx := GetTx(ctx); tx != nil {
			return tx
		}
		return txFunc(ctx)
	}
}

// MustGetTx retrieves the transaction from context or panics
// Use this when you're certain a transaction should be present
func MustGetTx(ctx context.Context) *gorm.DB {
	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	// Log error before panicking for debugging
	if logger := ctx.Value("logger"); logger != nil {
		if zapLogger, ok := logger.(*zap.Logger); ok {
			zapLogger.Panic("transaction not found in context")
		}
	}
	panic("transaction not found in context - ensure SetTx was called")
}
