package dbtesting

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Env represents different test environments
type Env int

const (
	// EnvTest creates unique database per test (isolated, slower startup)
	EnvTest Env = iota
	// EnvDev uses shared development database (faster, requires external setup)
	EnvDev Env = iota
)

func (e Env) String() string {
	switch e {
	case EnvTest:
		return "test"
	case EnvDev:
		return "dev"
	default:
		return "unknown"
	}
}

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

// GetConfig returns database config for environment
func GetConfig(env Env) Config {
	switch env {
	case EnvTest:
		return Config{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "postgres",
		}
	case EnvDev:
		return Config{
			Host:     "localhost",
			Port:     5433, // Different port for dev
			User:     "postgres",
			Password: "devpassword",
			Database: "nova_dev",
		}
	default:
		return Config{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "postgres",
		}
	}
}

// Database options for flexible test configuration
type dbOptions struct {
	DebugOff            bool                   // Turn off SQL query logging
	NoWrapInTransaction bool                   // Skip transaction wrapping
	PostInitHooks       []func(*gorm.DB) error // Hooks to run after DB initialization (in committed transaction)
}

// DBOption configures database behavior
type DBOption func(*dbOptions)

// DBDebugOff disables SQL query logging for cleaner test output
var DBDebugOff DBOption = func(o *dbOptions) {
	o.DebugOff = true
}

// DBNoWrapInTransaction skips automatic transaction wrapping
var DBNoWrapInTransaction DBOption = func(o *dbOptions) {
	o.NoWrapInTransaction = true
}

// DBWithHook adds a post-initialization hook that runs in a committed transaction
func DBWithHook(hook func(*gorm.DB) error) DBOption {
	return func(o *dbOptions) {
		o.PostInitHooks = append(o.PostInitHooks, hook)
	}
}

// Connection cache for performance
var connections = map[string]*gorm.DB{}
var connectionsMutex = &sync.Mutex{}

func getCachedDB(connString string) (*gorm.DB, error) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	if db, exists := connections[connString]; exists {
		return db, nil
	}

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	connections[connString] = db
	return db, nil
}

// DefaultConfig returns config for db-setup pattern (backwards compatibility)
func DefaultConfig() Config {
	return GetConfig(EnvTest)
}

// CreateTestDB creates test database with environment and options support
func CreateTestDB(t *testing.T, env Env, options ...DBOption) *gorm.DB {
	var opts dbOptions
	for _, option := range options {
		option(&opts)
	}

	config := GetConfig(env)
	var db *gorm.DB

	switch env {
	case EnvTest:
		// Connect to base database using cache
		baseDB, err := getCachedDB(config.ConnString())
		require.NoError(t, err, "failed to connect to base database")

		// Test database connectivity
		var version string
		err = baseDB.Raw("SELECT version()").Row().Scan(&version)
		require.NoError(t, err)
		require.NotEmpty(t, version)
		t.Logf("Database version: %s", version)

		// Create unique test database
		testDBName := fmt.Sprintf("test_db_%d", rand.Intn(10000000))
		err = baseDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName)).Error
		require.NoError(t, err)

		// Connect to test database
		config.Database = testDBName
		logLevel := logger.Info
		if opts.DebugOff {
			logLevel = logger.Error
		}

		testDB, err := gorm.Open(postgres.Open(config.ConnString()), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
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

		db = testDB

	case EnvDev:
		// Connect to shared development database
		logLevel := logger.Info
		if opts.DebugOff {
			logLevel = logger.Error
		}

		devDB, err := gorm.Open(postgres.Open(config.ConnString()), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})

		if err != nil {
			t.Skipf("Dev database not available: %v", err)
			return nil
		}

		// Test connectivity
		var version string
		err = devDB.Raw("SELECT version()").Row().Scan(&version)
		if err != nil {
			t.Skipf("Dev database not accessible: %v", err)
			return nil
		}
		t.Logf("Dev database version: %s", version)

		db = devDB

	default:
		t.Fatalf("Unknown environment: %v", env)
		return nil
	}

	// Run post-initialization hooks in committed transactions
	for i, hook := range opts.PostInitHooks {
		t.Logf("Running post-init hook %d", i+1)
		err := hook(db)
		require.NoError(t, err, "Post-init hook %d failed", i+1)
	}

	// Wrap in transaction unless disabled
	if !opts.NoWrapInTransaction {
		tx := db.Begin()
		require.NoError(t, tx.Error)

		t.Cleanup(func() {
			tx.Rollback()
		})

		db = tx
	}

	return db
}

// CreateTestDB creates isolated test database (backwards compatibility)
func CreateTestDBLegacy(t *testing.T) *gorm.DB {
	return CreateTestDB(t, EnvTest)
}

// CreateTestDBWithTx creates test database wrapped in transaction (backwards compatibility)
func CreateTestDBWithTx(t *testing.T) *gorm.DB {
	return CreateTestDB(t, EnvTest) // Default behavior includes transaction wrapping
}
