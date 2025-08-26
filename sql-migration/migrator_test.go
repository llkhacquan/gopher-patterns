package migration

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrator(t *testing.T) {
	// Use db-setup pattern - assumes PostgreSQL is running on localhost:5432
	config := Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Database: "postgres",
		SSLMode:  "disable",
	}

	t.Run("Migration Up and Down", func(t *testing.T) {
		migrator, err := NewMigrator(config)
		require.NoError(t, err)
		defer migrator.Close()

		ctx := context.Background()

		// Get initial version
		initialVersion, err := migrator.Version(ctx)
		require.NoError(t, err)

		// Run migrations up
		err = migrator.Up(ctx)
		require.NoError(t, err)

		// Check version increased
		newVersion, err := migrator.Version(ctx)
		require.NoError(t, err)
		assert.Greater(t, newVersion, initialVersion)

		// Verify tables were created
		db := migrator.db

		// Check users table exists
		var exists bool
		err = db.QueryRow(`SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'users'
		)`).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "users table should exist")

		// Check orders table exists
		err = db.QueryRow(`SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'orders'
		)`).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "orders table should exist")

		// Test data insertion (verify schema works)
		_, err = db.Exec(`INSERT INTO users (name, email) VALUES ($1, $2)`, "Test User", "test@example.com")
		require.NoError(t, err)

		var userID int
		err = db.QueryRow(`SELECT id FROM users WHERE email = $1`, "test@example.com").Scan(&userID)
		require.NoError(t, err)

		_, err = db.Exec(`INSERT INTO orders (user_id, product, quantity, price) VALUES ($1, $2, $3, $4)`,
			userID, "Test Product", 2, 19.99)
		require.NoError(t, err)

		// Cleanup: Roll back migrations
		err = migrator.Down(ctx)
		require.NoError(t, err)

		err = migrator.Down(ctx)
		require.NoError(t, err)
	})

	t.Run("Migrator from existing DB connection", func(t *testing.T) {
		// Connect directly
		db, err := sql.Open("postgres", config.ConnString())
		require.NoError(t, err)
		defer db.Close()

		// Create migrator from existing connection
		migrator := NewMigratorFromDB(db)
		ctx := context.Background()

		// Test status (should work even with no migrations)
		err = migrator.Status(ctx)
		require.NoError(t, err)
	})

	t.Run("Get embedded migrations", func(t *testing.T) {
		files, err := GetEmbeddedMigrations()
		require.NoError(t, err)

		// Should have our test migrations
		assert.Len(t, files, 2)
		assert.Contains(t, files, "migrations/001_create_users.sql")
		assert.Contains(t, files, "migrations/002_create_orders.sql")
	})
}
