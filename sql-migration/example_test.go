package migration

import (
	"context"
	"fmt"
	"log"
	"testing"
)

// TestMigrationExample demonstrates complete migration usage
func TestMigrationExample(t *testing.T) {
	// Database configuration (assumes db-setup pattern is running)
	config := Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Database: "postgres",
		SSLMode:  "disable",
	}

	// Create migrator
	migrator, err := NewMigrator(config)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer migrator.Close()

	ctx := context.Background()

	// Show embedded migrations
	fmt.Println("📁 Embedded migrations:")
	files, err := GetEmbeddedMigrations()
	if err != nil {
		log.Fatalf("Failed to get migrations: %v", err)
	}
	for _, file := range files {
		fmt.Printf("  - %s\n", file)
	}

	// Check current version
	version, err := migrator.Version(ctx)
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}
	fmt.Printf("\n📊 Current migration version: %d\n", version)

	// Run migrations
	fmt.Println("\n🚀 Running migrations...")
	if err := migrator.Up(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	fmt.Println("✅ Migrations completed successfully")

	// Check new version
	version, err = migrator.Version(ctx)
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}
	fmt.Printf("📊 New migration version: %d\n", version)

	// Show migration status
	fmt.Println("\n📋 Migration status:")
	if err := migrator.Status(ctx); err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}

	// Clean up for demo (rollback migrations)
	fmt.Println("\n🔄 Rolling back for cleanup...")
	if err := migrator.Down(ctx); err != nil {
		log.Printf("Rollback 1 failed (may be expected): %v", err)
	}
	if err := migrator.Down(ctx); err != nil {
		log.Printf("Rollback 2 failed (may be expected): %v", err)
	}
	fmt.Println("✅ Example completed")
}
