package dbtesting

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database connection configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// ConnString returns PostgreSQL connection string
func (c Config) ConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database)
}

// DefaultConfig returns config for db-setup pattern
func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Database: "postgres",
	}
}

// CreateTestDB creates isolated test database with automatic cleanup
func CreateTestDB(t *testing.T) *gorm.DB {
	config := DefaultConfig()

	// Connect to default database
	baseDB, err := gorm.Open(postgres.Open(config.ConnString()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	require.NoError(t, err)

	// Create unique test database
	testDBName := fmt.Sprintf("test_db_%d", rand.Intn(1000000))
	err = baseDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName)).Error
	require.NoError(t, err)

	// Connect to test database
	config.Database = testDBName
	testDB, err := gorm.Open(postgres.Open(config.ConnString()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	require.NoError(t, err)

	// Cleanup on test completion
	t.Cleanup(func() {
		sqlDB, _ := testDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		baseDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	})

	return testDB
}

// CreateTestDBWithTx creates test database wrapped in transaction for isolation
func CreateTestDBWithTx(t *testing.T) *gorm.DB {
	db := CreateTestDB(t)

	// Start transaction
	tx := db.Begin()
	require.NoError(t, tx.Error)

	// Rollback transaction on cleanup
	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}
