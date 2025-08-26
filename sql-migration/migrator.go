package migration

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Config holds database connection configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// ConnString returns PostgreSQL connection string
func (c Config) ConnString() string {
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, sslMode)
}

// Migrator handles database migrations using embedded SQL files
type Migrator struct {
	db *sql.DB
}

// NewMigrator creates a new migrator with database connection
func NewMigrator(config Config) (*Migrator, error) {
	db, err := sql.Open("postgres", config.ConnString())
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return &Migrator{db: db}, nil
}

// NewMigratorFromDB creates a migrator from existing database connection
func NewMigratorFromDB(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// Up runs all pending migrations
func (m *Migrator) Up(ctx context.Context) error {
	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "failed to set dialect")
	}

	if err := goose.UpContext(ctx, m.db, "migrations"); err != nil {
		return errors.Wrap(err, "failed to run migrations")
	}

	return nil
}

// Down rolls back one migration
func (m *Migrator) Down(ctx context.Context) error {
	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "failed to set dialect")
	}

	if err := goose.DownContext(ctx, m.db, "migrations"); err != nil {
		return errors.Wrap(err, "failed to rollback migration")
	}

	return nil
}

// Status returns migration status
func (m *Migrator) Status(ctx context.Context) error {
	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "failed to set dialect")
	}

	if err := goose.StatusContext(ctx, m.db, "migrations"); err != nil {
		return errors.Wrap(err, "failed to get migration status")
	}

	return nil
}

// Version returns current migration version
func (m *Migrator) Version(ctx context.Context) (int64, error) {
	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return 0, errors.Wrap(err, "failed to set dialect")
	}

	version, err := goose.GetDBVersionContext(ctx, m.db)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get database version")
	}

	return version, nil
}

// Close closes the database connection
func (m *Migrator) Close() error {
	return m.db.Close()
}

// GetEmbeddedMigrations returns list of embedded migration files for inspection
func GetEmbeddedMigrations() ([]string, error) {
	var files []string

	err := fs.WalkDir(migrationFS, "migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && path != "migrations" {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
