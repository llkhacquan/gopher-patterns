package dbtesting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Simple migrator for testing - mimics sql-migration pattern
type testMigrator struct {
	db *gorm.DB
}

func newTestMigrator(db *gorm.DB) *testMigrator {
	return &testMigrator{db: db}
}

func (m *testMigrator) up(ctx context.Context) error {
	// Simple migration: create orders table
	return m.db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			product VARCHAR(100) NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`).Error
}

func (m *testMigrator) down(ctx context.Context) error {
	return m.db.Exec("DROP TABLE IF EXISTS orders").Error
}

// Order model for testing migrated schema
type Order struct {
	ID      uint    `gorm:"primaryKey"`
	UserID  uint    `gorm:"not null"`
	Product string  `gorm:"size:100;not null"`
	Amount  float64 `gorm:"type:decimal(10,2);not null"`
}

func TestSQLMigrationIntegration(t *testing.T) {
	t.Run("Hooks run migrations before tests", func(t *testing.T) {
		// Create migration hook
		migrationHook := func(db *gorm.DB) error {
			migrator := newTestMigrator(db)
			return migrator.up(context.Background())
		}

		// Create database with migration hook
		db := CreateTestDB(t, EnvTest,
			DBDebugOff,                // Clean output
			DBWithHook(migrationHook), // Run migration
		)

		// Test that migrated schema works
		order := Order{
			UserID:  1,
			Product: "Test Product",
			Amount:  99.99,
		}

		err := db.Create(&order).Error
		require.NoError(t, err)
		assert.NotZero(t, order.ID)

		// Verify order was created
		var found Order
		err = db.First(&found, order.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Test Product", found.Product)
		assert.Equal(t, 99.99, found.Amount)
	})

	t.Run("Multiple hooks run in sequence", func(t *testing.T) {
		// Hook 1: Create orders table
		createOrdersHook := func(db *gorm.DB) error {
			return db.Exec(`
				CREATE TABLE IF NOT EXISTS orders (
					id BIGSERIAL PRIMARY KEY,
					user_id BIGINT NOT NULL,
					product VARCHAR(100) NOT NULL,
					amount DECIMAL(10,2) NOT NULL,
					created_at TIMESTAMP DEFAULT NOW()
				)
			`).Error
		}

		// Hook 2: Insert test data
		seedDataHook := func(db *gorm.DB) error {
			return db.Exec(`
				INSERT INTO orders (user_id, product, amount) VALUES
				(1, 'Seeded Product 1', 19.99),
				(2, 'Seeded Product 2', 29.99)
			`).Error
		}

		// Create database with multiple hooks
		db := CreateTestDB(t, EnvTest,
			DBDebugOff,
			DBWithHook(createOrdersHook),
			DBWithHook(seedDataHook),
		)

		// Verify seeded data exists
		var count int64
		err := db.Model(&Order{}).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// Verify specific seeded data
		var orders []Order
		err = db.Find(&orders).Error
		require.NoError(t, err)
		assert.Len(t, orders, 2)
		assert.Equal(t, "Seeded Product 1", orders[0].Product)
		assert.Equal(t, "Seeded Product 2", orders[1].Product)
	})

	t.Run("Hooks work with transaction isolation", func(t *testing.T) {
		migrationHook := func(db *gorm.DB) error {
			migrator := newTestMigrator(db)
			return migrator.up(context.Background())
		}

		// Create two isolated test databases
		db1 := CreateTestDB(t, EnvTest, DBWithHook(migrationHook))
		db2 := CreateTestDB(t, EnvTest, DBWithHook(migrationHook))

		// Create different orders in each database
		order1 := Order{UserID: 1, Product: "Product A", Amount: 10.00}
		order2 := Order{UserID: 2, Product: "Product B", Amount: 20.00}

		err := db1.Create(&order1).Error
		require.NoError(t, err)
		err = db2.Create(&order2).Error
		require.NoError(t, err)

		// Verify isolation - each database only sees its own data
		var count1, count2 int64
		db1.Model(&Order{}).Count(&count1)
		db2.Model(&Order{}).Count(&count2)

		assert.Equal(t, int64(1), count1)
		assert.Equal(t, int64(1), count2)

		// Verify different products in each database
		var found1, found2 Order
		db1.First(&found1)
		db2.First(&found2)

		assert.Equal(t, "Product A", found1.Product)
		assert.Equal(t, "Product B", found2.Product)
	})

	t.Run("Hooks work without transaction wrapping", func(t *testing.T) {
		migrationHook := func(db *gorm.DB) error {
			migrator := newTestMigrator(db)
			return migrator.up(context.Background())
		}

		// Create database with migration but no transaction wrapping
		db := CreateTestDB(t, EnvTest,
			DBDebugOff,
			DBNoWrapInTransaction, // Data persists
			DBWithHook(migrationHook),
		)

		// Create order
		order := Order{UserID: 1, Product: "Persistent Product", Amount: 50.00}
		err := db.Create(&order).Error
		require.NoError(t, err)

		// Verify order persists (no transaction rollback)
		var found Order
		err = db.First(&found, order.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Persistent Product", found.Product)
	})
}

func TestHooksWithMigrationCleanup(t *testing.T) {
	t.Run("Migration up and down hooks", func(t *testing.T) {
		var migrator *testMigrator

		// Hook to create migrator and run up migration
		upHook := func(db *gorm.DB) error {
			migrator = newTestMigrator(db)
			return migrator.up(context.Background())
		}

		db := CreateTestDB(t, EnvTest,
			DBDebugOff,
			DBNoWrapInTransaction, // Need persistence for cleanup test
			DBWithHook(upHook),
		)

		// Test that migration worked
		order := Order{UserID: 1, Product: "Migration Test", Amount: 75.00}
		err := db.Create(&order).Error
		require.NoError(t, err)

		// Clean up with down migration
		t.Cleanup(func() {
			if migrator != nil {
				err := migrator.down(context.Background())
				assert.NoError(t, err)
			}
		})
	})
}
